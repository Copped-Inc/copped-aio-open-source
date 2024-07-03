package linked_roles

import (
	"bytes"
	"database-api/user"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/responses"

	"github.com/Copped-Inc/aio-types/discord"
	secret "github.com/Copped-Inc/aio-types/secrets"
	"github.com/Copped-Inc/aio-types/subscriptions"
	"github.com/infinitare/disgo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	var (
		code, state []string
		ok          bool
		resData     helper.DiscordTokenResp
		userData    disgo.User
		secrets     user.Secrets
		client      = &http.Client{}
		rawSecrets  []byte
	)

	//
	// validate redirect / request
	//

	query := r.URL.Query()

	if code, ok = query["code"]; !ok {
		console.ErrorRequest(w, r, errors.New("request doesn't include OAuth2 code query string parameter"), http.StatusBadRequest)
		return
	}

	if state, ok = query["state"]; !ok {
		console.ErrorRequest(w, r, errors.New("request missing state query string parameter"), http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("client_state")
	if err != nil {
		console.ErrorRequest(w, r, errors.New("validation failed: couldn't find state cookie"), http.StatusBadRequest)
		return
	}

	if cookie.Value != state[0] {
		console.ErrorRequest(w, r, errors.New("authorization failed: state doesn't match"), http.StatusUnauthorized)
		return
	}

	//
	// get tokens from code
	//

	form := url.Values{}
	form.Add("code", code[0])
	form.Add("grant_type", "authorization_code")
	form.Add("redirect_uri", "https://database.copped-inc.com/discord/linked-roles")
	form.Add("scope", "identify%20email%20guilds.join%20role_connections.write")

	req, err := http.NewRequest(http.MethodPost, "https://discord.com/api/v"+discord.API_Version+"/oauth2/token", strings.NewReader(form.Encode()))
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(discord.Application_ID, discord.Oauth2_Secret)

	res, err := client.Do(req)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	//
	// get user
	//

	req, err = http.NewRequest(http.MethodGet, "https://discord.com/api/v"+discord.API_Version+"/users/@me", nil)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	req.Header.Add("Authorization", "Bearer "+resData.AccessToken)

	res, err = client.Do(req)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	if err = json.NewDecoder(res.Body).Decode(&userData); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	//
	// try getting db userdata
	//

	db, err := user.FromId(string(userData.ID))
	if err != nil {
		if status.Code(err) == codes.NotFound {
			responses.Redirect(w, r, helper.Active+"/utility/403?message="+url.QueryEscape("You must have an active subsciption in order to connect Copped AIO to your Discord server profile."))
		} else {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	if db.User.Subscription.Plan == subscriptions.None {
		responses.Redirect(w, r, helper.Active+"/utility/403?message="+url.QueryEscape("You must have an active subsciption in order to connect Copped AIO to your Discord server profile."))
		return
	}

	//
	// update Oauth2 data / database
	//

	if len(db.Secrets) > 0 {
		rawSecrets, err = helper.Decrypt(db.Secrets, secret.JWT_Secret)
		if err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}

		if err = json.Unmarshal(rawSecrets, &secrets); err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}
	}

	secrets.Oauth2.Expiry = time.Now().Add(time.Second * time.Duration(resData.ExpiresIn))
	secrets.Oauth2.AccessToken = resData.AccessToken
	secrets.Oauth2.RefreshToken = resData.RefreshToken

	if rawSecrets, err = json.Marshal(secrets); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	db.Secrets, err = helper.Encrypt(rawSecrets, secret.JWT_Secret)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	if err = db.Update(); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	//
	// update user role metadata
	//

	payload, err := json.Marshal(struct {
		Platform_name string            `json:"platform_name"`
		Metadata      map[string]string `json:"metadata"`
	}{
		Platform_name: "Copped AIO",
		Metadata:      map[string]string{strconv.Itoa(int(db.User.Subscription.Plan)): "1"},
	})
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	req, err = http.NewRequest(http.MethodPut, "https://discord.com/api/v"+discord.API_Version+"/users/@me/applications/"+discord.Application_ID+"/role-connection", bytes.NewBuffer(payload))
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+resData.AccessToken)

	if _, err = client.Do(req); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	responses.Redirect(w, r, helper.Active+"/utility/discord")
}
