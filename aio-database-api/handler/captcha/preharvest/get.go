package preharvest

import (
	"database-api/database"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"

	"cloud.google.com/go/firestore"
	"github.com/Copped-Inc/aio-types/captcha/preharvest"
	"github.com/gorilla/mux"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func get(w http.ResponseWriter, r *http.Request) {
	var (
		ref   *firestore.DocumentRef
		task  preharvest.Task
		tasks []preharvest.Task
		err   error

		order = firestore.Desc
		query = database.GetDatabase().Collection("preharvest").Query
	)

	getAll := func() error {
		iter := query.Documents(database.GetContext())
		for {
			item, err := iter.Next()
			if err != nil {
				if err == iterator.Done {
					break
				}
				return err
			}

			var task preharvest.Task

			if err = item.DataTo(&task); err != nil {
				return err
			}

			task.ID = item.Ref.ID
			tasks = append(tasks, task)
		}

		return nil
	}

	getOne := func() error {
		doc, err := ref.Get(database.GetContext())
		if err != nil {
			return err
		}

		if err = doc.DataTo(&task); err != nil {
			return err
		}

		task.ID = doc.Ref.ID

		return nil
	}

	if task_id, ok := mux.Vars(r)["task-id"]; ok {
		ref = database.GetDatabase().Doc("preharvest/" + task_id)

	} else {
		if user_id, ok := mux.Vars(r)["user-id"]; ok {
			query = query.Where("user_id", "==", user_id)
		}

		// check if any querystring params are present and if so, apply them to the firestore query
		if params := r.URL.Query(); len(params) != 0 {
			var start_after time.Time

			if after, ok := params["after"]; ok {
				after, err := strconv.ParseInt(after[len(after)-1], 10, 64)
				if err != nil {
					console.ErrorRequest(w, r, err, http.StatusInternalServerError)
					return
				}
				start_after = time.Unix(after, 0)
				query = query.StartAfter(start_after)
			}

			if before, ok := params["before"]; ok {
				before, err := strconv.ParseInt(before[len(before)-1], 10, 64)
				if err != nil {
					console.ErrorRequest(w, r, err, http.StatusInternalServerError)
					return
				}
				if end_before := time.Unix(before, 0); !end_before.After(start_after) && !start_after.IsZero() {
					console.ErrorRequest(w, r, errors.New("before querystring param unix timestamp mustn't be before after timestamp"), http.StatusBadRequest)
					return
				} else {
					query = query.EndBefore(end_before)
				}
			}

			if limit, ok := params["limit"]; ok {
				limit, err := strconv.ParseInt(limit[len(limit)-1], 10, 64)
				if err != nil {
					console.ErrorRequest(w, r, err, http.StatusInternalServerError)
					return
				} else if limit < 1 {
					console.ErrorRequest(w, r, errors.New("limit query parameter must be > 0"), http.StatusBadRequest)
					return
				}
				query = query.Limit(int(limit))
			}

			if sort, ok := params["sort"]; ok {
				sort, err := strconv.ParseInt(sort[len(sort)-1], 10, 64)
				if err != nil {
					console.ErrorRequest(w, r, err, http.StatusInternalServerError)
					return
				}
				order = firestore.Direction(sort)
			}
		}
	}

	if query = query.OrderBy("date", order); ref != nil {
		err = getOne()
	} else {
		err = getAll()
	}
	if err != nil {
		if status.Code(err) == codes.NotFound {
			console.ErrorRequest(w, r, err, http.StatusNotFound)
		} else {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	responses.SendJson(
		func() interface{} {
			if len(tasks) != 0 {
				return tasks
			} else if task.ID != "" {
				return task
			}
			return nil
		}(),
		func() int {
			if len(tasks) != 0 || task.ID != "" {
				return http.StatusOK
			}
			return http.StatusNoContent
		}(),
		w, r,
	)
}
