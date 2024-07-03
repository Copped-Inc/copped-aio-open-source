package handler

import (
	"errors"
	"fmt"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"
	"github.com/gorilla/mux"
	"net/http"
	"os/exec"
)

func post(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	if id == "" {
		console.ErrorRequest(w, r, errors.New(http.StatusText(http.StatusBadRequest)), http.StatusBadRequest)
		return
	}

	console.RequestLog(r, "Creating Instance", id)
	console.RequestLog(r, "Run Container")

	cmdStr := fmt.Sprintf("sudo docker run -d -e INSTANCE_ID=%s --name client-%s --restart on-failure client", id, id)
	o, err := exec.Command("/bin/sh", "-c", cmdStr).Output()
	if err != nil {
		console.ErrorRequest(w, r, errors.New(string(o)), http.StatusInternalServerError)
		return
	}

	responses.SendOk(w, r)

}
