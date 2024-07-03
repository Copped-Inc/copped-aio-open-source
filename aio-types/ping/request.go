package ping

import (
	"errors"
	"net/http"
)

func request(domain string) error {

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, "https://"+domain+"/ping", nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("Ping failed, status code: " + res.Status)
	}

	return err

}
