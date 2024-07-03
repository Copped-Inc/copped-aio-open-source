package insta

import (
	"database-api/handler/instance"
	"database-api/user"
	"errors"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/responses"
	"github.com/Copped-Inc/aio-types/secrets"
	"net/http"
	"time"
)

func post(w http.ResponseWriter, r *http.Request, database *user.Database) {

	if database.User.InstanceLimit <= len(database.Instances) {
		console.ErrorRequest(w, r, errors.New("instance limit reached"), http.StatusForbidden)
		return
	}

	if time.Since(database.User.CodeExpire) > 0 || database.User.Code == "" {
		code := instance.RandCode()
		err := database.SetCode(code).Update()
		if err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}
	}

	req, err := http.NewRequest(http.MethodPost, helper.ActiveInstances+"/instance/"+database.User.Code, nil)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	req.Header.Set("Password", secrets.API_Admin_PW)
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	if res.StatusCode != http.StatusOK {
		console.ErrorRequest(w, r, errors.New(res.Status), res.StatusCode)
		return
	}

	responses.SendOk(w, r)

}
