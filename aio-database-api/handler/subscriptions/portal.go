package subscriptions

import (
	"database-api/user"
	"net/http"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/billingportal/session"
)

func Portal(w http.ResponseWriter, r *http.Request) {

	user, err := user.FromRequest(r)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusUnauthorized)
		return
	}

	portal, err := session.New(&stripe.BillingPortalSessionParams{Customer: stripe.String(string(user.User.Subscription.Customer_ID))})
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	responses.Redirect(w, r, portal.URL)

}
