package proxies

import (
	"encoding/json"
	"errors"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"net/http"
	"time"

	"github.com/Copped-Inc/aio-types/secrets"
)

var proxies []Proxy
var resis []Proxy
var updated time.Time

func getProxies(t string) ([]Proxy, error) {

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, helper.ActiveService+"/proxies", nil)
	req.Header.Set("type", t)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Password", secrets.API_Admin_PW)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var resBody response
	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return nil, err
	}

	return resBody.Proxies, err

}

func handle(t string) []Proxy {
	if len(proxies) < 1 || time.Now().Sub(updated).Minutes() > 30 {
		var err error
		p, err := getProxies(t)
		if err != nil {
			console.ErrorLog(errors.New("failed to get proxies: " + err.Error()))
			if len(proxies) < 1 {
				panic(err)
			}
		} else {
			proxies = p
		}
		updated = time.Now()
	}

	return proxies
}

func dcs() []Proxy {
	return handle("dcs")
}

func residentials() []Proxy {
	return handle("residential")
}

type response struct {
	Proxies []Proxy `json:"proxies"`
}

type Proxy struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}
