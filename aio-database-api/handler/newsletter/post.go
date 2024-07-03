package newsletter

import (
	"database-api/database"
	"database-api/mail"
	"database-api/user"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/responses"

	"google.golang.org/api/iterator"
)

func post(w http.ResponseWriter, r *http.Request) {

	if !helper.IsMaster(r.Header.Get("Password")) {
		console.ErrorRequest(w, r, errors.New("invalid authorization password"), http.StatusUnauthorized)
		return
	}

	var data mail.Mail
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		console.ErrorRequest(w, r, err, http.StatusBadRequest)
		return
	}

	users := database.GetDatabase().Collection("data").Where("user.email", "!=", "").Documents(database.GetContext())

	for {
		u, err := users.Next()
		if err == iterator.Done {
			break
		}
		if err == nil {
			var d user.Database
			err = u.DataTo(&d)
			if err == nil {
				go data.Send(d.User.Email)
			}
		}
	}

	responses.SendOk(w, r)
}
