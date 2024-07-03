package preharvest

import (
	"database-api/database"
	preharvestcache "database-api/preharvest"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"

	"github.com/Copped-Inc/aio-types/captcha/preharvest"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func patch(w http.ResponseWriter, r *http.Request) {
	var (
		task, original preharvest.Task
		changes        preharvest.Task_Edit
	)

	if err := json.NewDecoder(r.Body).Decode(&changes); err != nil {
		console.ErrorRequest(w, r, errors.New(http.StatusText(http.StatusBadRequest)), http.StatusBadRequest)
		return
	} else if changes.Uses == 0 && changes.State == 0 && changes.Schedule == "" && changes.Date.IsZero() {
		console.ErrorRequest(w, r, errors.New("no request payload was provided"), http.StatusNotModified)
		return
	} else if changes.Schedule != "" && !preharvest.Schedule_Pattern.MatchString(changes.Schedule) {
		console.ErrorRequest(w, r, errors.New("schedule didn't match pattern regex"), http.StatusBadRequest)
		return
	} else if changes.Date.Before(time.Now()) && !changes.Date.IsZero() {
		console.ErrorRequest(w, r, errors.New("date mustn't be changed to a value in the past"), http.StatusBadRequest)
		return
	}

	doc, err := database.GetDatabase().Doc("preharvest/" + mux.Vars(r)["task-id"]).Get(database.GetContext())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			console.ErrorRequest(w, r, err, http.StatusNotFound)
		} else {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	if err = doc.DataTo(&task); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	if original = task; changes.State > 0 {
		task.State = changes.State
	}
	if changes.Uses > 0 {
		if changes.Uses > 7 {
			task.Uses = 7
		} else {
			task.Uses = changes.Uses
		}
	}
	if task.State == preharvest.Stopped {
		task.Uses = 0
	} else if task.Uses < 1 {
		task.State = preharvest.Stopped
	}
	if changes.Schedule != "" {
		task.Schedule = changes.Schedule
	}
	if !changes.Date.IsZero() {
		task.Date = changes.Date
	}

	if task == original {
		console.ErrorRequest(w, r, errors.New("no changes provided"), http.StatusNotModified)
		return
	}

	task.ID = doc.Ref.ID

	if _, err = doc.Ref.Set(database.GetContext(), task); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	responses.SendJson(task, http.StatusOK, w, r)
	preharvestcache.Update <- task
}
