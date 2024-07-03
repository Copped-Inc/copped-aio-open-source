package payment

import (
	"database-api/mail"
	"database-api/user"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	consts "github.com/Copped-Inc/aio-types/user"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/invoice"
	"github.com/stripe/stripe-go/v74/webhook"
)

func Failed(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	event, err := webhook.ConstructEvent(data, r.Header.Get("Stripe-Signature"), "whsec_dyGHRvuPAQQ3bKUR5tNXUCYRSzTQcOOU")
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusForbidden)
		return
	}

	var payload *stripe.Invoice

	if err = json.Unmarshal(event.Data.Raw, &payload); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	var attempted int

	if attempts, ok := payload.Metadata["attempts"]; ok {
		if attemptes, err := strconv.ParseInt(attempts, 10, 64); err != nil {
			console.ErrorRequest(w, r, err, http.StatusBadRequest)
			return
		} else {
			attempted = int(attemptes)
		}
	}

	mail := mail.New()
	mail.Title = "Subscription fee payment failed!"
	mail.Button = true
	mail.ButtonUrl = payload.HostedInvoiceURL
	mail.ButtonText = "pay now"

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

	if attempted++; attempted < 2 {
		mail.Title = "Subscription fee payment failed!"
		mail.Text = "Automatic payment for your most recent subscription fee failed. Note that your subscription will be deactived after 5 days and finally canceled after 10 days without payment."
	} else {

		if attempted == 2 {
			u.User.Subscription.State = consts.Pending
			mail.Text = "Automatic payment for your most recent subscription fee failed for a second time. Note that your subscription will be canceled if no payment is provided within the next 5 days. Also, we might take further actions to ensure payment is collected including legal ones."

		} else {
			u.User.Subscription.State = consts.Expired
			mail.Title = "Subscription canceled!"
			mail.Text = "Automatic payment for your most recent subscription fee failed three times in a row. Since you didn't take appropiate action within 10 days your subscription was canceled. Please make sure to pay the remaining charges. Else we may have to take further actions to ensure payment is collected including legal ones."
			// TODO delete user subscription
		}

		if err = u.Update(); err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}
	}

	mail.Send(u.User.Email)

	params := &stripe.InvoiceParams{}
	params.Metadata = map[string]string{"attempts": strconv.Itoa(attempted)}
	payload, err = invoice.Update(payload.ID, params)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	responses.SendOk(w, r)
}
