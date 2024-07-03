package notifications

import (
	"database-api/database"
	"database-api/handler/websocket"
	"database-api/user"
	"encoding/json"
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

func patch(w http.ResponseWriter, r *http.Request) {
	var (
		notification, original user.Notification
		req                    Request
		ref                    *firestore.DocumentRef

		notificationID = mux.Vars(r)["notification-id"]
	)

	administrator, userID, code, err := authenticate(r)
	if err != nil {
		console.ErrorRequest(w, r, err, code)
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		console.ErrorRequest(w, r, err, http.StatusBadRequest)
		return
	} else if req.Read == nil && req.Text == "" && req.Title == "" {
		console.ErrorRequest(w, r, errors.New("no field in request body provided"), http.StatusBadRequest)
		return
	}

	if global := !strings.Contains(r.URL.Path, "users"); userID != "" && !global {
		ref = database.GetDatabase().Collection("data").Doc(userID).Collection("notifications").Doc(notificationID)

	} else if administrator && global {
		ref = database.GetDatabase().Collection("notifications").Doc(notificationID)

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

	if err := get_notification(ref, &notification); err != nil {
		if status.Code(err) == codes.NotFound {
			console.ErrorRequest(w, r, err, http.StatusNotFound)
		} else {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	if original = notification; req.Read != nil && userID != "" {
		*notification.Read = *req.Read
	}

	if administrator && !original.Global {
		if d := req.Text; d != "" {
			notification.Text = d
		}
		if d := req.Title; d != "" {
			notification.Title = d
		}
	}

	if original == notification {
		console.ErrorRequest(w, r, errors.New("no changes provided"), http.StatusNotModified)
	}

	if _, err = ref.Set(database.GetContext(), notification); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	if userID == "" {
		websocket.NotificationUpdate(notification)
	} else {
		websocket.Websocket{Action: websocket.UpdateNotification, Body: notification}.Send(userID)
	}

	responses.SendJson(notification, http.StatusOK, w, r)
}
