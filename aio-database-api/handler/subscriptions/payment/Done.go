package payment

import (
	"database-api/user"
	"encoding/json"
	"io"
	"net/http"

	consts "github.com/Copped-Inc/aio-types/user"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/subscription"
	"github.com/stripe/stripe-go/v74/webhook"
)

func Done(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	event, err := webhook.ConstructEvent(data, r.Header.Get("Stripe-Signature"), "whsec_FRhxTs25W30gDUzjlE9xWpO6sUFrPfDj")
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusForbidden)
		return
	}

	var payload *stripe.Invoice

	if err = json.Unmarshal(event.Data.Raw, &payload); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	subscription, err := subscription.Get(payload.Subscription.ID, nil)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	if subscription.Status == stripe.SubscriptionStatusActive {
		customer, err := customer.Get(payload.Customer.ID, nil)
		if err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}

		u, err := user.FromId(customer.Metadata["discord_id"])
		if err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}

		if u.User.Subscription.State != consts.Active {
			u.User.Subscription.State = consts.Active

			if err = u.Update(); err != nil {
				console.ErrorRequest(w, r, err, http.StatusInternalServerError)
				return
			}
		}
	}

	responses.SendOk(w, r)
}
