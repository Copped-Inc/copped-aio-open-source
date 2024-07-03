package instock

import (
	"database-api/global"
	"database-api/handler/websocket"
	"database-api/product"
	"encoding/json"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"
	"github.com/Copped-Inc/aio-types/webhook"
	"golang.org/x/exp/slices"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func post(w http.ResponseWriter, r *http.Request) {

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		console.ErrorRequest(w, r, err, http.StatusBadRequest)
		return
	}

	go global.AddPing()

	split := strings.Split(req.Link, "/")
	p, err := product.Get(req.SKU)
	if err != nil {
		p = product.New(req.SKU, product.Blacklisted)
		go global.AddProduct()

		err = p.AddHandle(split[len(split)-1]).
			UpdateName(strings.Split(req.Name, " [")[0]).
			UpdateImage(req.Image).
			UpdateState(product.Blacklisted).
			UpdatePrice(req.Price).
			Save()

		if err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}

		responses.SendJson(p, http.StatusOK, w, r)
		return
	}

	update := false
	if req.Name != "" && p.Name == "" {
		p.UpdateName(strings.Split(req.Name, " [")[0])
		update = true
	}

	if req.Image != "" && p.Image == "" {
		p.UpdateImage(req.Image)
		update = true
	}

	if !slices.Contains(p.Handles, split[len(split)-1]) {
		p.AddHandle(split[len(split)-1])
		update = true
	}

	if req.Price != 0 && p.Price == 0.0 {
		p.UpdatePrice(req.Price)
		update = true
	}

	if p.State == product.Blacklisted {
		for _, v := range p.UserState {
			if v.State == product.Whitelisted {
				go websocket.UserMonitor(req, v.ID)
			}
		}
		responses.SendJson(p, http.StatusOK, w, r)
		go save(p, update)
		return
	}

	var l []string
	for _, v := range p.UserState {
		if v.State == product.Blacklisted {
			l = append(l, v.ID)
		}
	}

	websocket.Monitor(req, l)
	responses.SendJson(p, http.StatusOK, w, r)

	go save(p, update)
	go func() {
		if strings.Contains(strings.ToLower(req.Name), "test") {
			return
		}

		webhook := webhook.New().AddEmbed(
			webhook.NewProduct,
			req.Name,
			"[Link]("+req.Link+")",
			func() string {
				if p.StockX != "" {
					return "[StockX](" + p.StockX + ")"
				} else {
					return "[StockX](https://stockx.com/search?s=" + p.SKU + ")"
				}
			}(),
			strconv.FormatFloat(req.Price, 'f', 2, 64),
			req.Image,
		)

		if strings.Contains(strings.ToLower(req.Link), "vfb") {
			webhook.SendMultiple([]string{
				"", // INSERT Webhook URL here
			})
		} else {
			webhook.SendMultiple([]string{
				"", // INSERT Webhook URL here
			})
		}
	}()

}

func save(p product.Product, update bool) {
	if update {
		err := p.Save()
		if err != nil {
			console.ErrorLog(err)
		}
	}
}

type request struct {
	Name  string            `json:"name"`
	SKU   string            `json:"sku"`
	Skus  map[string]string `json:"skus"`
	Date  time.Time         `json:"date"`
	Link  string            `json:"link"`
	Image string            `json:"image"`
	Price float64           `json:"price"`
}
