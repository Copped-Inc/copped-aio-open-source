package notifications

import (
	"database-api/database"
	"database-api/handler/websocket"
	"database-api/user"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/responses"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"google.golang.org/api/iterator"
)

func post(w http.ResponseWriter, r *http.Request) {
	var (
		req Request
		res *firestore.DocumentRef

		unread = false
	)

	if !helper.IsMaster(r.Header.Get("Password")) {
		console.ErrorRequest(w, r, errors.New("invalid authorization password"), http.StatusUnauthorized)
		return
	} else if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		console.ErrorRequest(w, r, err, http.StatusBadRequest)
		return
	}

	data := user.Notification{
		Title:     req.Title,
		Text:      req.Text,
		CreatedAt: time.Now(),
	}

	if data.Text == "" || data.Title == "" {
		console.ErrorRequest(w, r, errors.New("notification text and title mustn't be empty"), http.StatusBadRequest)
		return
	}

	add := func(ref *firestore.CollectionRef) (err error) {
		res, _, err = ref.Add(database.GetContext(), data)
		data.ID = res.ID
		return
	}

	if id, private := mux.Vars(r)["user-id"]; !private {
		if err := add(database.GetDatabase().Collection("notifications")); err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}

		queue := database.GetDatabase().Collection("data").Documents(database.GetContext())
		for {
			item, err := queue.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				console.ErrorRequest(w, r, err, http.StatusInternalServerError)
				return
			}

			if _, err = item.Ref.Collection("notifications").Doc(res.ID).Set(database.GetContext(), user.Notification{Global: true, CreatedAt: time.Now(), Read: &unread}); err != nil {
				console.ErrorRequest(w, r, err, http.StatusInternalServerError)
				return
			}
		}

		websocket.NotificationCreate(data)

	} else {

		if req.Read != nil {
			data.Read = req.Read
		} else {
			data.Read = &unread
		}

		if err := add(database.GetDatabase().Collection("data").Doc(id).Collection("notifications")); err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}

		websocket.Websocket{Action: websocket.CreateNotification, Body: data}.Send(id)
	}

	responses.SendJson(data, http.StatusCreated, w, r)
}
