package update

import (
	"bytes"
	"database-api/handler/notifications"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"

	"github.com/Copped-Inc/aio-types/secrets"
)

func crosspost(body request) {
	version, err := os.ReadFile("version-client-local-v3")
	if err != nil {
		console.Log(err.Error())
		return
	}

	payload := notifications.Request{Title: "Version " + string(version) + " live.", Text: "That's new:"}

	for _, s := range body.Changelog {
		payload.Text += "\n• " + s
	}

	data, err := json.Marshal(payload)
	if err != nil {
		console.Log(err.Error())
		return
	}

	req, err := http.NewRequest(http.MethodPost, helper.ActiveData+"/notifications", bytes.NewBuffer(data))
	if err != nil {
		console.Log(err.Error())
		return
	}

	req.Header.Set("Password", secrets.API_Admin_PW)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	for i := 0; i < 4; i++ {

		res, err := client.Do(req)
		if res.StatusCode != http.StatusCreated || err != nil {
			time.Sleep(time.Second * 10 * time.Duration(i))
			continue
		}

		return
	}

	console.Log("Failed to send update notification.", err.Error())
}
