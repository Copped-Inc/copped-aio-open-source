package purchase

import (
	"database-api/link"
	"database-api/user"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/cookies"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/responses"
	"github.com/Copped-Inc/aio-types/subscriptions"
	"github.com/gorilla/mux"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"

	"github.com/Copped-Inc/aio-types/discord"
)

var running = false

func getHandle(w http.ResponseWriter, r *http.Request) {

	for running {
		time.Sleep(time.Millisecond * 100)
	}
	running = true

	get(w, r)

	running = false

}

func get(w http.ResponseWriter, r *http.Request) {

	keys, ok := r.URL.Query()["code"]

	if !ok || len(keys[0]) < 1 {
		cookies.Add(w, "redirect", r.URL.Path)
		responses.Redirect(w, r, "https://discord.com/api/oauth2/authorize?client_id="+discord.Application_ID+"&redirect_uri="+url.QueryEscape(helper.ActiveData+"/")+"&response_type=code&scope=identify%20email%20guilds.join")
		return
	}

	l, err := link.Get(mux.Vars(r)["code"])
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	mResp, token, err := helper.GetDiscordResp(keys[0], helper.ActiveData)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	if l.Stock < 1 {
		d, err := user.FromId(mResp.Id)
		if err == nil {
			jwt, err := d.User.Jwt()
			if err != nil {
				console.ErrorRequest(w, r, err, http.StatusInternalServerError)
				return
			}

			cookies.Add(w, "authorization", jwt)
			cookies.Remove(w, "redirect")
			cookies.Remove(w, "code")
			responses.Redirect(w, r, helper.Active)
			return
		}
		console.ErrorRequest(w, r, errors.New("link is out of stock"), 601)
		return
	}

	if !mResp.Verified {
		console.ErrorRequest(w, r, errors.New("not verified"), http.StatusUnauthorized)
		return
	}

	d, err := user.New(mResp, l.Plan, l.InstanceLimit)
	if err != nil {
		responses.Redirect(w, r, helper.Active)
		return
	}

	if err = l.Use().Update(); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	jwt, err := d.User.Jwt()
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	if err = helper.JoinServer(mResp.Id, token); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	cookies.Add(w, "authorization", jwt)
	cookies.Remove(w, "redirect")
	cookies.Remove(w, "code")

	if l.Plan != subscriptions.Developer {
		checkout, err := session.New(&stripe.CheckoutSessionParams{
			SuccessURL: stripe.String(helper.Active),
			Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
			LineItems: []*stripe.CheckoutSessionLineItemParams{
				{Price: stripe.String(l.Plan.GetData().Price)},
			},
			ClientReferenceID:        stripe.String(mResp.Id),
			BillingAddressCollection: stripe.String(string(stripe.CheckoutSessionBillingAddressCollectionRequired)),
			PaymentMethodCollection:  stripe.String(string(stripe.CheckoutSessionPaymentMethodCollectionAlways)),
			ExpiresAt:                stripe.Int64(time.Now().Add(time.Hour).Unix()),
			AfterExpiration: &stripe.CheckoutSessionAfterExpirationParams{
				Recovery: &stripe.CheckoutSessionAfterExpirationRecoveryParams{
					Enabled: stripe.Bool(true),
				},
			},
		})
		if err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}

		d.User.Subscription.Customer_ID = checkout.ID
		if err = d.Update(); err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}

		responses.Redirect(w, r, checkout.URL)
	} else {
		responses.Redirect(w, r, helper.Active)
	}
}
