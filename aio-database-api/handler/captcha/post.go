package captcha

import (
	"encoding/json"
	"errors"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func post(w http.ResponseWriter, r *http.Request) {

	s := mux.Vars(r)["site"]
	if s == "" {
		console.ErrorRequest(w, r, errors.New("sitekey is empty"), http.StatusBadRequest)
		return
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusBadRequest)
		return
	}

	runningSolver.remove(s, req.TaskId)

	if req.ErrorId != 0 {
		console.ErrorRequest(w, r, errors.New("not solved"), http.StatusBadRequest)
		return
	}

	c := Captcha{
		Token:  req.Solution.GRecaptchaResponse,
		Expire: time.Now().Add(time.Minute * 1),
	}

	if r.Header.Get("Browser") == "" {
		queue.set(s, append(queue.getSite(s), c))
	} else {
		queue.set(s, append([]Captcha{c}, queue.getSite(s)...))
	}

	responses.SendOk(w, r)

}

type request struct {
	ErrorId  int    `json:"errorId"`
	TaskId   string `json:"taskId"`
	Status   string `json:"status"`
	Solution struct {
		GRecaptchaResponse string `json:"gRecaptchaResponse"`
	} `json:"solution"`
}
