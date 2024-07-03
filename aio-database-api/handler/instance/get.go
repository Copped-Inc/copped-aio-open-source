package instance

import (
	"database-api/user"
	"errors"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"
	"math/rand"
	"net/http"
	"time"
)

func get(w http.ResponseWriter, r *http.Request, database *user.Database) {

	if database.User.InstanceLimit <= len(database.Instances) {
		console.ErrorRequest(w, r, errors.New("instance limit reached"), http.StatusForbidden)
		return
	}

	if time.Since(database.User.CodeExpire) < 0 && database.User.Code != "" {
		responses.SendJson(response{Code: database.User.Code}, http.StatusOK, w, r)
		return
	}

	code := RandCode()
	err := database.SetCode(code).Update()
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	responses.SendJson(response{Code: code}, http.StatusOK, w, r)

}

var letterRunes = []rune("0123456789")

func RandCode() string {
	b := make([]rune, 6)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type response struct {
	Code string `json:"code"`
}
