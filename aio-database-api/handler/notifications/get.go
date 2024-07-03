package notifications

import (
	"database-api/database"
	"database-api/user"
	"errors"
	"net/http"
	"strings"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func get(w http.ResponseWriter, r *http.Request) {
	var (
		notification  user.Notification
		notifications []user.Notification
		queue         *firestore.CollectionRef
		ref           *firestore.DocumentRef

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

	} else if notificationID != "" {
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
		if err := get_notification(ref, &notification); err != nil {
			if status.Code(err) == codes.NotFound {
				console.ErrorRequest(w, r, err, http.StatusNotFound)
			} else {
				console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			}
			return
		}

		responses.SendJson(notification, http.StatusOK, w, r)

	} else if docs, err := queue.OrderBy("created_at", firestore.Desc).Documents(database.GetContext()).GetAll(); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)

	} else if len(docs) == 0 {
		responses.SendJson(nil, http.StatusNoContent, w, r)

	} else {
		for _, doc := range docs {
			var notification user.Notification
			if err = doc.DataTo(&notification); err != nil {
				console.ErrorRequest(w, r, err, http.StatusInternalServerError)
				return
			}

			notification.ID = doc.Ref.ID
			notifications = append(notifications, notification)
		}

		responses.SendJson(notifications, http.StatusOK, w, r)
	}
}
