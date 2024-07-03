package webhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func (b *Body) Send(to string) error {

	body, err := json.Marshal(b)
	if err != nil {
		return err
	}

	res, err := http.Post(to, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if res.StatusCode != 204 {
		body, _ := io.ReadAll(res.Body)
		return errors.New("request to " + res.Request.URL.Host + res.Request.URL.Path + " failed with " + res.Status + ": " + string(body))
	}

	return err

}

func (b *Body) Update(url string) error {

	body, err := json.Marshal(b)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		return errors.New("request to " + res.Request.URL.Host + res.Request.URL.Path + " failed with " + res.Status + ": " + string(body))
	}

	return err

}

func (b *Body) SendMultiple(to []string) {

	for _, e := range to {
		_ = b.Send(e)
	}

}
