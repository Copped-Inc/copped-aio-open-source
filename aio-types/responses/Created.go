package responses

import (
	"net/http"
)

func SendCreated(w http.ResponseWriter, r *http.Request) {

	SendJson(struct {
		Status string `json:"status"`
	}{
		Status: "created",
	}, http.StatusCreated, w, r)

}
