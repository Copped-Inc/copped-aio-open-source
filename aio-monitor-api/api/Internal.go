package api

import (
	"bytes"
	"errors"
	"github.com/Copped-Inc/aio-types/helper"
	"net/http"
)

var KithEuSiteKey = ""

func InternalGet(path string) (*http.Response, error) {

	req, err := http.NewRequest(http.MethodGet, helper.ActiveData+"/"+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Monitor")
	req.Header.Set("Password", "E&BWLL*2ROrv~8P19TN-f4_ZU=%c67[5")
	req.Header.Set("KithEU-Current-Sitekey", KithEuSiteKey)

	res, err := (&http.Client{}).Do(req)
	if err == nil {
		if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
			err = errors.New("request to " + res.Request.URL.String() + " failed with response " + res.Status)
		}
	}

	return res, err
}

func InternalPost(data []byte, path string) (*http.Response, error) {

	req, err := http.NewRequest(http.MethodPost, helper.ActiveData+"/"+path, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Monitor")
	req.Header.Set("Password", "E&BWLL*2ROrv~8P19TN-f4_ZU=%c67[5")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("KithEU-Current-Sitekey", KithEuSiteKey)

	res, err := (&http.Client{}).Do(req)
	if err == nil {
		if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
			err = errors.New("request to " + res.Request.URL.String() + " failed with response " + res.Status)
		}
	}

	return res, err
}
