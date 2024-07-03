package instance

import (
	"database-api/database"
	"database-api/user"
	"encoding/json"
	"net/http"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"
)

func post(w http.ResponseWriter, r *http.Request) {

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusBadRequest)
		return
	}

	doc, err := database.GetDatabase().Collection("data").Where("user.code", "==", req.Code).Documents(database.GetContext()).Next()
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	db, err := user.FromId(doc.Ref.ID)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	err = db.SetCode("").Update()
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	jwt, err := db.User.Jwt()
	responses.SendJson(responseInstance{Authorization: jwt}, http.StatusOK, w, r)

}

type request struct {
	Code string `json:"code"`
}

type responseInstance struct {
	Authorization string `json:"authorization"`
}
