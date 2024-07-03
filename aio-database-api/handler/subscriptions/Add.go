package subscriptions

import (
	"database-api/handler/subscriptions/payment"
	"net/http"

	"github.com/Copped-Inc/aio-types/console"

	"github.com/gorilla/mux"
)

func Add(r *mux.Router) {
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := r.Header["Stripe-Signature"]; !ok {
				console.ErrorRequest(w, r, nil, http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		})
	})

	r.HandleFunc("/checkout-session-completed", create)
	r.HandleFunc("/customer-subscription-delete", delete)
	r.HandleFunc("/invoice-payment-failed", payment.Failed)
	r.HandleFunc("/invoice-paid", payment.Done)
}
