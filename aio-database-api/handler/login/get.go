package login

import (
	"database-api/mail"
	"database-api/user"
	"net/http"
	"net/url"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/cookies"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/responses"

	"golang.org/x/exp/slices"

	"github.com/Copped-Inc/aio-types/discord"
)

func get(w http.ResponseWriter, r *http.Request) {
	if _, err := helper.GetClaim("id", r); err == nil {
		redirect(w, r)
		return
	}

	keys, ok := r.URL.Query()["code"]
	if !ok || len(keys[0]) < 1 {
		cookies.Add(w, "redirect", r.URL.Path)
		responses.Redirect(w, r, "https://discord.com/api/oauth2/authorize?client_id="+discord.Application_ID+"&redirect_uri="+url.QueryEscape(helper.ActiveData+"/")+"&response_type=code&scope=identify%20email%20guilds.join")
		return
	}

	mResp, _, err := helper.GetDiscordResp(keys[0], helper.ActiveData)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	if mResp.Id == "" {
		console.Log("Invalid code")
		cookies.Add(w, "redirect", r.URL.Path)
		responses.Redirect(w, r, "https://discord.com/api/oauth2/authorize?client_id="+discord.Application_ID+"&redirect_uri="+url.QueryEscape(helper.ActiveData+"/")+"&response_type=code&scope=identify%20email%20guilds.join")
		return
	}

	data, err := user.FromId(mResp.Id)
	if err != nil {
		responses.Redirect(w, r, helper.Active+"/utility/403")
		return
	}

	jwt, err := data.User.Jwt()
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	ip, err := helper.GetIP(w, r)
	if err != nil {
		if helper.System != "windows" {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}

		ip = "localhost"
		err = nil
	}

	if !slices.Contains(data.IPs, ip) {
		go func() {
			mail.New().
				SetTitle("New Login detected").
				SetSubtitle("A new device has been logged in. If you have not logged in, please open a ticket in the Copped AIO Discord server.").
				SetBelowButton("Browser used: " + r.UserAgent() + "\nIP Address: " + ip).
				Send(data.User.Email)

			data.IPs = append(data.IPs, ip)
			err = data.Update()
			if err != nil {
				console.ErrorLog(err)
			}
		}()
	}

	cookies.Add(w, "authorization", jwt)
	cookies.Remove(w, "redirect")
	cookies.Remove(w, "code")
	redirect(w, r)

}

func redirect(w http.ResponseWriter, r *http.Request) {
	cookie := cookies.Get(r, "after_login")
	if cookie != nil {
		responses.Redirect(w, r, helper.Active+cookie.Value)
		return
	}

	responses.Redirect(w, r, helper.Active)
}
