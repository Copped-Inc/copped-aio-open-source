package notifications

import (
	"database-api/database"
	"database-api/handler/websocket"
	"database-api/user"
	"errors"
	"net/http"
	"strings"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
)

func delete(w http.ResponseWriter, r *http.Request) {
	var (
		queue *firestore.CollectionRef
		ref   *firestore.DocumentRef

		notificationID = mux.Vars(r)["notification-id"]
	)

	administrator, userID, code, err := authenticate(r)
	if err != nil {
		console.ErrorRequest(w, r, err, code)
		return
	}

	if global := !strings.Contains(r.URL.Path, "users"); userID != "" && notificationID != "" && !global {
		ref = database.GetDatabase().Collection("data").Doc(userID).Collection("notifications").Doc(notificationID)

	} else if userID != "" && !global {
		queue = database.GetDatabase().Collection("data").Doc(userID).Collection("notifications")

	} else if notificationID != "" && administrator {
		ref = database.GetDatabase().Collection("notifications").Doc(notificationID)

	} else if administrator && global {
		queue = database.GetDatabase().Collection("notifications")

	} else {
		console.ErrorRequest(w, r, func() error {
			if administrator {
				return errors.New("requesting @me endpoint with admin pw")
			} else {
				return nil
			}
		}(), http.StatusForbidden)
		return
	}

	if ref != nil {
		if _, err = ref.Delete(database.GetContext()); err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}

		responses.SendJson(nil, http.StatusNoContent, w, r)
		if userID != "" {
			websocket.Websocket{Action: websocket.DeleteNotification, Body: ref.ID}.Send(userID)
		} else {
			websocket.NotificationDelete(user.Notification{ID: notificationID})
		}

	} else if refs, err := queue.DocumentRefs(database.GetContext()).GetAll(); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)

	} else {
		for _, ref := range refs {
			if _, err = ref.Delete(database.GetContext()); err != nil {
				console.ErrorRequest(w, r, err, http.StatusInternalServerError)
				return
			}
		}

		responses.SendJson(nil, http.StatusNoContent, w, r)

		if len(refs) > 0 {
			if userID != "" {
				websocket.Websocket{Action: websocket.DeleteNotification}.Send(userID)
			} else {
				websocket.NotificationDelete(user.Notification{})
			}
		}
	}
}
