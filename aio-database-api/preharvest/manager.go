package preharvest

import (
	"errors"
	"math"
	"net/http"
	"sort"
	"time"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"

	"github.com/Copped-Inc/aio-types/captcha/preharvest"
	"github.com/Copped-Inc/aio-types/secrets"
)

var (
	cache   Tasks
	stop    chan any
	current preharvest.Task

	Add    = make(chan preharvest.Task)
	Remove = make(chan string)
	Update = make(chan preharvest.Task)
)

// manages list of active preharvest tasks
func manager() {
	for {
		if len(cache) != 0 {
			current = cache[len(cache)-1]
		}

		//
		// wait for any actions to be made to the task list
		//
		select {

		//
		// if a new captcha preharvest task is added, check whether it needs to be added to the list of active / executing tasks
		//
		case task := <-Add:
			newTask := task

			// delete tasks that should've been already run in the past and shouldn't repeat
			if !newTask.Routine && newTask.Date.Before(time.Now()) {

				go func() {
					if err := delete(newTask.ID); err != nil {
						console.ErrorLog(err)
					}
				}()
				continue

			} else
			// stop tasks without any uses remaining
			if newTask.Routine && newTask.Uses < 1 && newTask.State == preharvest.Running {

				go func() {
					if err := patch(newTask.ID, preharvest.Task_Edit{State: preharvest.Stopped}); err != nil {
						console.ErrorLog(err)
					}
				}()
				continue

			} else
			// update tasks that ran in the past to change their execution date to the next date following the schedule, if possible
			// otherwise stop the task
			if changes := (preharvest.Task_Edit{}); newTask.Date.Before(time.Now()) && newTask.State == preharvest.Running {

				if newTask.Schedule != "" {
					changes.Date = newTask.Date.Add(ParseSchedule(newTask.Schedule) * time.Duration(math.Ceil(float64(time.Since(newTask.Date))/float64(ParseSchedule(newTask.Schedule)))))
				} else {
					changes.State = preharvest.Stopped
				}

				go func() {
					if err := patch(newTask.ID, changes); err != nil {
						console.ErrorLog(err)
					}
				}()
				continue

			} else
			// add running tasks to the list of running tasks
			if newTask.State == preharvest.Running {
				cache = append(cache, newTask)
			}

		//
		// remove a preharvest task from the list of active preharvest tasks
		//
		case id := <-Remove:
			match := false

			for i, task := range cache {
				if match = task.ID == id; match {
					cache = append(cache[:i], cache[:i+1]...)
					break
				}
			}

			if !match {
				continue
			}

		//
		// update a preharvest task
		//
		case task := <-Update:
			changedTask := task

			if changedTask.State == preharvest.Stopped {
				go func() { Remove <- changedTask.ID }()
				continue
			}

			match := false

			for i, task := range cache {
				if match = task.ID == changedTask.ID; match {
					cache[i] = changedTask
					break
				}
			}

			if !match {
				go func() { Add <- changedTask }()
				continue
			}
		}

		//
		// check whether sorting of the tasks needs to be updated
		//
		if !sort.IsSorted(cache) {
			sort.Sort(cache)
		}

		//
		// if there are preharvest tasks and the currently executing preharvest task changed, restart manager
		// else stop manager, in case there is still one running
		//
		if len(cache) != 0 {
			if current != cache[len(cache)-1] {
				if stop != nil {
					stop <- nil
					stop = nil
				}

				go execute(cache[len(cache)-1])
			}

		} else if stop != nil {
			stop <- nil
			stop = nil
		}
	}
}

// manager executes the last task from the task list, meaning the one to be executed next
func execute(task preharvest.Task) {
	stop = make(chan any)
	timer := time.NewTimer(time.Until(task.Date))
	defer timer.Stop()

	select {
	case <-timer.C:
		// request captcha
		go func() {
			req, err := http.NewRequest(http.MethodGet, helper.ActiveData+"/captcha/"+string(task.Site), nil)
			if err != nil {
				console.ErrorLog(err)
				return
			}

			req.Header.Add("Password", secrets.API_Admin_PW)

			res, err := (&http.Client{}).Do(req)
			if err != nil {
				console.ErrorLog(err)
				return
			}

			if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotModified {
				console.ErrorLog(errors.New("request to " + res.Request.URL.String() + " failed with response " + res.Status))
			}
		}()

		// tasks that shouldn't be repeated are deleted after executing once
		if !task.Routine {

			go func() {
				if err := delete(task.ID); err != nil {
					console.ErrorLog(err)
				}
			}()

		} else
		// for scheduled tasks, decrease uses and  if necessary stop the task
		// also update schedule if possible, otherwise stop task too
		{

			if task.Uses--; task.Uses < 1 || task.Schedule == "" {
				task.State = preharvest.Stopped
			}

			go func() {
				if err := patch(task.ID, preharvest.Task_Edit{
					Uses: task.Uses,
					Date: func() (date time.Time) {
						if task.Schedule != "" {
							date = task.Date.Add(ParseSchedule(task.Schedule) * time.Duration(math.Ceil(float64(time.Since(task.Date))/float64(ParseSchedule(task.Schedule)))))
						}
						return
					}(),
					State: task.State},
				); err != nil {
					console.ErrorLog(err)
				}
			}()
		}

	case <-stop:
	}
}
