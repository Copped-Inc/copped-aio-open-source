package subscriptions

import (
	"database-api/user"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"
	u "github.com/Copped-Inc/aio-types/user"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/webhook"
)

func create(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	event, err := webhook.ConstructEvent(data, r.Header.Get("Stripe-Signature"), "") // INSERT Stripe Webhook Secret
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusForbidden)
		return
	}

	var payload stripe.CheckoutSession

	if err = json.Unmarshal(event.Data.Raw, &payload); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	db, err := user.FromId(payload.ClientReferenceID)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	db.User.Subscription.Customer_ID = payload.Customer.ID
	db.User.Subscription.State = u.Active

	for _, item := range payload.Subscription.Items.Data {
		if item.Price.ID == db.User.Subscription.Plan.GetData().Price {
			db.User.Subscription.Subscription_ID = item.ID
			break
		}
	}

	if err = db.Update(); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	params := &stripe.CustomerParams{}
	params.AddMetadata("discord_id", db.User.ID)

	if _, err = customer.Update(payload.Customer.ID, params); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	responses.SendOk(w, r)
}
