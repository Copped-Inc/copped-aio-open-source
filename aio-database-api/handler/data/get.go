package data

import (
	"bytes"
	"database-api/product"
	"database-api/user"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/responses"
	u "github.com/Copped-Inc/aio-types/user"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"golang.org/x/exp/slices"
)

func get(w http.ResponseWriter, r *http.Request, database *user.Database) {
	if state := database.User.Subscription.State; state != u.Active {
		if state == u.Disabled {
			responses.Redirect(w, r, helper.Active+"utility/403?error=forbidden&message="+url.QueryEscape("You have pending payments left open for more than 5 days already. Please pay them first to regain access to the dashboard and Copped AIO.\nFor more details check your inbox at "+database.User.Email+".")+"&button=false")
		} else if state == u.Expired {
			responses.Redirect(w, r, helper.Active+"utility/403?error=forbidden&message="+url.QueryEscape("Your subscription expired or was cancelled, thus you can no longer access Copped AIO.")+"&button=false")
		} else if state == u.Pending {
			session, err := session.Get(database.User.Subscription.Customer_ID, nil)
			if err != nil {
				console.ErrorRequest(w, r, err, http.StatusInternalServerError)
				return
			}

			if payment_status := session.Status; payment_status == stripe.CheckoutSessionStatusOpen {
				responses.Redirect(w, r, helper.Active+"utility/409?error="+url.QueryEscape("payment pending")+"&message="+url.QueryEscape("You have yet to complete the stripe checkout in order to gain access to Copped AIO.&location")+"&location="+url.QueryEscape(session.URL)+"&title="+url.QueryEscape("checkout"))
			} else if payment_status == stripe.CheckoutSessionStatusExpired {
				if time.Since(time.Unix(session.Created, 0)) < time.Hour*24*7 /* a week */ {
					responses.Redirect(w, r, helper.Active+"utility/409?error="+url.QueryEscape("payment pending")+"&message="+url.QueryEscape("You have yet to complete the stripe checkout in order to gain access to Copped AIO.\nSince your previous stripe checkout session expired, click below to get to a new checkout session and give it another try.")+"&location="+url.QueryEscape(session.AfterExpiration.Recovery.URL)+"&title="+url.QueryEscape("checkout"))
				} else {
					database.User.Subscription.State = u.Expired

					if err = database.Update(); err != nil {
						console.ErrorRequest(w, r, err, http.StatusInternalServerError)
						return
					}

					responses.Redirect(w, r, helper.Active+"utility/410?error=forbidden&message="+url.QueryEscape("You failed to complete checkout within 7 days. Your checkout link expired and you will have to participate in another drop to get access to Copped AIO.")+"&button=false")
				}
			} else {
				console.ErrorRequest(w, r, errors.New("pending user with unexpected checkout session state: "+string(payment_status)), http.StatusInternalServerError)
			}
		} else {
			console.ErrorRequest(w, r, errors.New("unexpected userstate:"+strconv.Itoa(int(state))), http.StatusBadRequest)
		}

		return
	}

	password := r.Header.Get("password")
	confirm := r.Header.Get("confirm")

	if password == "" && confirm == "" || len(database.Password) == 0 && confirm == "" {
		responses.SendJson(struct {
			Status string `json:"status"`
		}{
			Status: http.StatusText(http.StatusNonAuthoritativeInfo),
		}, http.StatusNonAuthoritativeInfo, w, r)
		return
	}

	if confirm != "" {
		database.Password = helper.CreateHash(password)
		d := user.Data{}
		j, err := json.Marshal(d)
		if err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}

		database.Data, err = helper.Encrypt(j, password)
		if err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}

		if err = database.Update(); err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}
	}

	if !bytes.Equal(database.Password, helper.CreateHash(password)) {
		console.ErrorRequest(w, r, errors.New("password is incorrect"), http.StatusForbidden)
		return
	}

	j, err := helper.Decrypt(database.Data, password)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	var d user.Data
	err = json.Unmarshal(j, &d)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	whitelist, err := product.GlobalWhitelisted()
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	for _, s := range database.Products {
		p, err := product.Get(s)
		if err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}

		for _, state := range p.UserState {
			if state.ID != database.User.ID {
				continue
			}

			if state.State == product.Whitelisted && !slices.Contains(whitelist, p.SKU) {
				whitelist = append(whitelist, p.SKU)
			} else if state.State == product.Blacklisted && slices.Contains(whitelist, p.SKU) {
				index := slices.Index(whitelist, p.SKU)
				whitelist = append(whitelist[:index], whitelist[index+1:]...)
			}
			break
		}
	}

	database.Products = whitelist
	res, err := database.ToDataResp(&d)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	responses.SendJson(res, http.StatusOK, w, r)
}
