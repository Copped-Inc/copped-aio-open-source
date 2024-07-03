package notifications

import (
	"database-api/database"
	"database-api/user"
	"errors"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/gorilla/mux"
)

func get_notification(ref *firestore.DocumentRef, notification *user.Notification) (err error) {
	doc, err := ref.Get(database.GetContext())
	if err != nil {
		return
	}

	if err = doc.DataTo(notification); err == nil {
		notification.ID = doc.Ref.ID
	}

	return
}

func authenticate(req *http.Request) (administrator bool, userID string, respCode int, err error) {
	var id interface{}

	if queryUser, ok := mux.Vars(req)["user-id"]; helper.IsMaster(req.Header.Get("Password")) {
		administrator = true
		if ok {
			userID = queryUser
		}

	} else if id, err = helper.GetClaim("id", req); err == nil {
		if userID = id.(string); queryUser != userID && ok {
			err = errors.New("request path includes a user id (" + queryUser + ") that is different from the user currently logged in (" + userID + ")\nto prevent this, provide relative \"@me\" instead of an absolute user id")
			respCode = http.StatusForbidden
		}

	} else {
		respCode = http.StatusUnauthorized
	}

	return
}
