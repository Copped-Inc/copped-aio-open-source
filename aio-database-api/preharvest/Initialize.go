package preharvest

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/worker"

	"github.com/Copped-Inc/aio-types/captcha/preharvest"
	"github.com/Copped-Inc/aio-types/secrets"
)

func Initialize() error {
	worker.Worker(manager)

	time.Sleep(time.Second * 30)

	req, err := http.NewRequest(http.MethodGet, helper.ActiveData+"/captcha/preharvest?sort="+url.QueryEscape(strconv.Itoa(int(firestore.Asc))), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Password", secrets.API_Admin_PW)

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		return errors.New("request to " + res.Request.URL.String() + " failed with response " + res.Status)
	}

	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	var tasks []preharvest.Task

	if err = json.NewDecoder(res.Body).Decode(&tasks); err != nil {
		return err
	}

	for _, task := range tasks {
		Add <- task
	}

	return nil
}
