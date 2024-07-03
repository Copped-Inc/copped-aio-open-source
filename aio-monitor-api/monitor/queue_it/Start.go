package queue_it

import (
	"monitor-api/request"
	"time"
)

func Start(url string) func() error {
	queue := Disabled
	return func() error {
		statusCode, err := getStatusCode(url)
		if err != nil {
			return err
		}

		if statusCode == 302 {
			queue = Disabled
			return err
		}

		if queue {
			return err
		}

		time.Sleep(20 * time.Second)
		statusCode, err = getStatusCode(url)
		if err != nil || statusCode == 302 {
			return err
		}

		queue = Enabled
		send(url)
		return err
	}
}

func getStatusCode(url string) (int, error) {
	res, err := request.Get(url, nil, nil, 1)
	if err != nil {
		return 0, err
	}

	res, err = request.Get(res.Header.Get("Location"), nil, nil, 1)
	return res.StatusCode, err
}
