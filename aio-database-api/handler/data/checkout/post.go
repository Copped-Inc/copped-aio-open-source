package checkout

import (
	"database-api/handler/websocket"
	"database-api/product"
	"database-api/user"
	"encoding/json"
	"github.com/Copped-Inc/aio-types/modules"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"
	"github.com/Copped-Inc/aio-types/subscriptions"
	"github.com/Copped-Inc/aio-types/webhook"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/usagerecord"
)

func post(w http.ResponseWriter, r *http.Request, database *user.Database) {
	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusBadRequest)
		return
	}

	handle := strings.Split(req.Link, "/")[len(strings.Split(req.Link, "/"))-1]
	productData, err := product.Get(handle)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusBadRequest)
		return
	}

	p := user.Product{
		User:  database.User.ID,
		Date:  time.Now(),
		Name:  productData.Name,
		Link:  req.Link,
		Image: productData.Image,
		Store: req.Site,
		Size:  req.Size,
		Price: productData.Price,
	}

	if err = database.AddCheckout(p); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	if database.User.Subscription.Plan != subscriptions.Developer {
		params := &stripe.UsageRecordParams{
			Quantity:         stripe.Int64(int64(p.Price * 100)),
			SubscriptionItem: stripe.String(database.User.Subscription.Subscription_ID),
		}
		params.IdempotencyKey = stripe.String(uuid.NewString())

		_, err = usagerecord.New(params)

		for i := 0; err != nil; i++ {
			if i > 4 {
				console.ErrorRequest(w, r, err, http.StatusServiceUnavailable)
				return
			}

			_, err = usagerecord.New(params)

			time.Sleep(time.Second * time.Duration(2^i))
		}
	}

	websocket.Websocket{
		Action: websocket.UpdateSession,
		Body: struct {
			Session user.Session `json:"session"`
		}{
			Session: *database.Session,
		},
	}.Send(database.User.ID)

	websocket.Websocket{
		Action: websocket.AddCheckout,
		Body: struct {
			Checkout user.Product `json:"checkout"`
		}{
			Checkout: p,
		},
	}.Send(database.User.ID)

	go func() {
		wh := webhook.New()
		wh.AddEmbed(
			webhook.NewCheckout,
			p.Name,
			req.Link,
			p.Image,
			string(p.Store),
			req.Size,
			strconv.FormatFloat(p.Price, 'f', 2, 64),
		)

		if !strings.Contains(strings.ToLower(p.Name), "test") {
			_ = wh.Send("") // INSERT Webhook URL here
		}

		if req.Checkout != "" {
			wh.AddEmbed(webhook.NewCheckoutLink, req.Checkout)
		}
		wh.SendMultiple(database.Settings.Webhooks)
	}()

	responses.SendOk(w, r)
}

type request struct {
	Link     string       `json:"link"`
	Site     modules.Site `json:"site"`
	Size     string       `json:"size"`
	Checkout string       `json:"checkout"`
}
