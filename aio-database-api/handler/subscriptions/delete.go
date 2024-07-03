package subscriptions

import (
	"database-api/user"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	u "github.com/Copped-Inc/aio-types/user"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/discord"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/responses"
	"github.com/Copped-Inc/aio-types/secrets"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/webhook"
)

func delete(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	event, err := webhook.ConstructEvent(data, r.Header.Get("Stripe-Signature"), "whsec_faIXgMnp0tfIBIGcpFrcHtl1hF0iusLW")
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusForbidden)
		return
	}

	var payload *stripe.Subscription

	if err = json.Unmarshal(event.Data.Raw, &payload); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	customer, err := customer.Get(payload.Customer.ID, nil)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	db, err := user.FromId(customer.Metadata["discord_id"])
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	var dbSecrets user.Secrets

	if len(db.Secrets) > 0 {

		rawSecrets, err := helper.Decrypt(db.Secrets, secrets.JWT_Secret)
		if err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}

		if err = json.Unmarshal(rawSecrets, &dbSecrets); err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}
	} else {
		console.ErrorRequest(w, r, errors.New("no user secrets found in db"), http.StatusUnauthorized)
		return
	}

	form := url.Values{}
	form.Add("token", dbSecrets.Oauth2.RefreshToken)

	req, err := http.NewRequest(http.MethodPost, "https://discord.com/api/v"+discord.API_Version+"/oauth2/token/revoke", strings.NewReader(form.Encode()))
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(discord.Application_ID, discord.Oauth2_Secret)

	if res, err := (&http.Client{}).Do(req); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	} else if res.StatusCode != http.StatusOK {
		console.ErrorRequest(w, r, errors.New("oauth2 token invalidation failed with "+strconv.Itoa(res.StatusCode)+" "+res.Status), http.StatusInternalServerError)
		return
	}

	rawSecrets, err := json.Marshal(dbSecrets)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	if db.Secrets, err = helper.Encrypt(rawSecrets, secrets.JWT_Secret); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	db.User.Subscription.State = u.Expired

	if err = db.Update(); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	responses.SendOk(w, r)
}
