package purchase

import (
	"database-api/link"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/responses"

	"github.com/Copped-Inc/aio-types/subscriptions"
)

func post(w http.ResponseWriter, r *http.Request) {

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusBadRequest)
		return
	}

	match := false

	for _, plan := range subscriptions.Plans {
		if match = req.Plan == plan; match {
			break
		}
	}

	if !match {
		console.ErrorRequest(w, r, errors.New("plan unavailable"), http.StatusBadRequest)
		return
	} else if req.Stock <= 0 || req.InstanceLimit <= 0 {
		console.ErrorRequest(w, r, errors.New("missing parameter"), http.StatusBadRequest)
		return
	}

	l := link.New().SetPlan(req.Plan).SetStock(req.Stock).SetInstanceLimit(req.InstanceLimit)
	err = l.Create()
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	resp := response{Link: helper.ActiveData + "/purchase/" + l.ID}
	responses.SendJson(resp, http.StatusCreated, w, r)
}

type request struct {
	Plan          subscriptions.Plan `json:"plan"`
	Stock         int                `json:"stock"`
	InstanceLimit int                `json:"instance_limit"`
}

type response struct {
	Link string `json:"link"`
}
