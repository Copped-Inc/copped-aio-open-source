package user

import (
	"database-api/user"
	"net/http"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"
	"github.com/Copped-Inc/aio-types/subscriptions"
	"github.com/gorilla/mux"
)

func get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id != "" {
		u, err := user.FromId(id)
		if err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}

		checkouts, err := u.GetCheckouts()
		if err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}

		pUser := struct {
			user.Database
			Checkouts []user.Product `json:"checkouts"`
		}{
			Database:  u,
			Checkouts: checkouts,
		}

		responses.SendJson(pUser, http.StatusOK, w, r)
		return
	}

	u, err := user.GetAll()
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	var data []response
	for _, d := range u {
		checkouts, err := d.GetCheckouts()
		if err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}

		res := response{
			ID:        d.User.ID,
			Name:      d.User.Name,
			Plan:      d.User.Subscription.Plan,
			Checkouts: checkouts,
			Instances: d.Instances,
			Picture:   d.User.Picture,
		}

		if d.Session != nil {
			res.Status = d.Session.Status
		}

		data = append(data, res)
	}

	if len(data) == 0 {
		console.ErrorRequest(w, r, err, http.StatusNotFound)
		return
	}

	responses.SendJson(data, http.StatusOK, w, r)
}

type response struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Plan      subscriptions.Plan `json:"plan"`
	Status    string             `json:"state,omitempty"`
	Checkouts []user.Product     `json:"checkouts,omitempty"`
	Instances []user.Instance    `json:"instances,omitempty"`
	Picture   string             `json:"picture"`
}
