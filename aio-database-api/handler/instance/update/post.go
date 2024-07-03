package update

import (
	"database-api/handler/websocket"
	"encoding/json"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/responses"
	"github.com/Copped-Inc/aio-types/secrets"
	"github.com/Copped-Inc/aio-types/webhook"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
)

func post(w http.ResponseWriter, r *http.Request) {

	f := mux.Vars(r)["file"]

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		console.ErrorRequest(w, r, err, http.StatusBadRequest)
		return
	}

	go func() {
		time.Sleep(time.Minute * 10)
		var err error

		if f == "payments" {
			err = DownloadPayments(r)
		} else if f == "client" {
			err = DownloadClient(r)
		}

		if err != nil {
			console.ErrorLog(err)
			return
		}

		updateReq, err := http.NewRequest(http.MethodPost, helper.ActiveInstances+"/update", nil)
		if err != nil {
			console.ErrorLog(err)
			return
		}

		updateReq.Header.Set("Password", secrets.API_Admin_PW)
		client := http.DefaultClient
		updateRes, err := client.Do(updateReq)
		if err != nil {
			console.ErrorLog(err)
			return
		}

		if updateRes.StatusCode != http.StatusOK {
			console.Log("Instances API returned status code " + string(rune(updateRes.StatusCode)) + " instead of 200 OK")
		}

		websocket.Update()
		if strings.Contains(strings.ToLower(req.Changelog[0]), "silent") {
			return
		}

		if len(req.Changelog) == 1 {
			req.Changelog = strings.Split(req.Changelog[0], ",")
			for i, log := range req.Changelog {
				req.Changelog[i] = strings.TrimSpace(log)
			}
		}

		field := *webhook.NewField("Changelog", "")
		for _, s := range req.Changelog {
			field.Value += "\n• " + s
		}

		if f == "payments" {
			go webhook.New().AddEmbed(webhook.UpdatePayments).SetFields(field).Send("") // INSERT Webhook URL here

		} else if f == "client" {
			go webhook.New().AddEmbed(webhook.UpdateClient).SetFields(field).Send("") // INSERT Webhook URL here
			go crosspost(req)
		}
	}()

	responses.SendOk(w, r)

}

type request struct {
	Changelog []string `json:"changelog"`
}
