package preharvest

import (
	"database-api/database"
	preharvestcache "database-api/preharvest"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"time"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"

	"github.com/Copped-Inc/aio-types/captcha/preharvest"
	"github.com/Copped-Inc/aio-types/modules"
	"github.com/gorilla/mux"
	"github.com/infinitare/disgo"
)

func post(w http.ResponseWriter, r *http.Request) {
	var task preharvest.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		console.ErrorRequest(w, r, errors.New("request body format invalid"), http.StatusBadRequest)
		return
	}

	if !func() bool {
		for _, site := range modules.Sites {
			if site == task.Site {
				return true
			}
		}
		return false
	}() {
		console.ErrorRequest(w, r, errors.New("invalid site request parameters"), http.StatusBadRequest)
		return
	} else if task.Date.Before(time.Now()) && !task.Routine {
		console.ErrorRequest(w, r, errors.New("an execution date may only be specified if it's in the future OR with routine set to \"true\""), http.StatusBadRequest)
		return
	} else if !task.Routine && task.Schedule != "" {
		console.ErrorRequest(w, r, errors.New("schedule can only be specified for repeated captcha preharvest tasks (\"routine\": true)"), http.StatusBadRequest)
		return
	} else if task.Schedule != "" && !preharvest.Schedule_Pattern.MatchString(task.Schedule) {
		console.ErrorRequest(w, r, errors.New("schedule doesn't match pattern regex"), http.StatusBadRequest)
		return
	}

	task.User_ID = disgo.Snowflake(mux.Vars(r)["user-id"])
	task.State = preharvest.Running

	if task.Schedule != "" && task.Date.Before(time.Now()) {
		task.Date.Add(preharvestcache.ParseSchedule(task.Schedule) * time.Duration(math.Ceil(float64(time.Since(task.Date))/float64(preharvestcache.ParseSchedule(task.Schedule)))))
	} else if task.Date.Before(time.Now()) {
		task.State = preharvest.Stopped
	}

	if task.Routine && task.State == preharvest.Running {
		task.Uses = 7
	}

	doc, _, err := database.GetDatabase().Collection("preharvest").Add(database.GetContext(), task)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	task.ID = doc.ID

	responses.SendJson(task, http.StatusCreated, w, r)
	preharvestcache.Add <- task
}
