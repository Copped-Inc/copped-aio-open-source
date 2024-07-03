package settings

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func Login(code string) (*Settings, error) {

	var r request
	r.Code = code

	j, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://database.copped-inc.com/instance", bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	var resp response
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}

	s := &Settings{
		Authorization: resp.Authorization,
		Id:            code,
		Price:         0,
		Provider:      "Payments",
		TaskMax:       0,
		Region:        "Unavailable",
	}

	err = s.Save()
	return s, err

}
