package notifications

import (
	"database-api/database"
	"database-api/handler/websocket"
	"encoding/json"
	"errors"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"
	"google.golang.org/api/iterator"
)

func put(w http.ResponseWriter, r *http.Request) {
	var (
		req     Request
		changes bool
	)

	_, userID, code, err := authenticate(r)
	if err != nil {
		console.ErrorRequest(w, r, err, code)
		return
	} else if userID == "" {
		console.ErrorRequest(w, r, errors.New("requesting @me endpoint with admin pw"), http.StatusForbidden)
		return
	} else if err = json.NewDecoder(r.Body).Decode(&req); err != nil || req.Read == nil {
		console.ErrorRequest(w, r, err, http.StatusBadRequest)
		return
	}

	for docs := database.GetDatabase().Collection("data").Doc(userID).Collection("notifications").Where("read", "!=", *req.Read).Documents(database.GetContext()); ; changes = true {
		doc, err := docs.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}

		if _, err := doc.Ref.Update(database.GetContext(), []firestore.Update{{Path: "read", Value: *req.Read}}); err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}
	}

	if changes {
		responses.SendOk(w, r)
		websocket.Websocket{Action: websocket.UpdateNotificationReadstate, Body: *req.Read}.Send(userID)

	} else {
		console.ErrorRequest(w, r, errors.New("no notification read state was changed"), http.StatusNotModified)
	}
}
