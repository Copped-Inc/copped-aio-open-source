package preharvest

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Copped-Inc/aio-types/helper"

	"github.com/Copped-Inc/aio-types/captcha/preharvest"
	"github.com/Copped-Inc/aio-types/secrets"
)

func delete(id string) error {
	var client = new(http.Client)

	req, err := http.NewRequest(http.MethodDelete, helper.ActiveData+"/captcha/preharvest/"+id, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Password", secrets.API_Admin_PW)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusNotFound {
		return errors.New("request to " + res.Request.URL.String() + " failed with response " + res.Status)
	}

	return nil
}

func patch(id string, payload preharvest.Task_Edit) error {
	var client = new(http.Client)

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, helper.ActiveData+"/captcha/preharvest/"+id, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Add("Password", secrets.API_Admin_PW)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotModified && res.StatusCode != http.StatusNotFound {
		return errors.New("request to " + res.Request.URL.String() + " failed with response " + res.Status)
	}

	return nil
}
