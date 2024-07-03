package payments

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

func (p *Payments) post() error {

	j, err := json.Marshal(&p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://database.copped-inc.com/instance/payments", bytes.NewBuffer(j))
	if err != nil {
		return err
	}

	req.Header.Set("Cookie", p.auth)

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("status code not 200")
	}

	return err

}
