package proxies

import (
	"errors"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/responses"
	"net/http"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	_, noUser := helper.GetClaim("id", r)
	proxystring, err := func() (string, error) {
		if noUser == nil {
			return fetch("https://cdn.discordapp.com/attachments/1042408245630877706/1201613053343572108/instances.txt")
		} else if r.Header.Get("type") == "residential" {
			return fetch("https://cdn.discordapp.com/attachments/1042408245630877706/1072974450309484615/resis.txt")
		} else {
			return fetch("https://cdn.discordapp.com/attachments/1042408245630877706/1199115602095263854/dcs.txt")
		}
	}()

	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	proxies := strings.Split(strings.ReplaceAll(strings.ReplaceAll(proxystring, " ", ""), "\r", ""), "\n")
	if len(proxies) < 1 {
		console.ErrorRequest(w, r, errors.New("no proxies found"), http.StatusInternalServerError)
		return
	}

	var proxiesArray []Proxy
	for i := 0; i < len(proxies); i++ {

		if proxies[i] == "" || len(strings.Split(proxies[i], ":")) < 4 {
			continue
		}

		proxies[i] = strings.ReplaceAll(proxies[i], "\r", "")

		proxy := Proxy{
			Ip:       strings.Split(proxies[i], ":")[0],
			Port:     strings.Split(proxies[i], ":")[1],
			Username: strings.Split(proxies[i], ":")[2],
			Password: strings.Split(proxies[i], ":")[3],
		}

		proxiesArray = append(proxiesArray, proxy)

	}

	if len(proxiesArray) < 1 {
		console.ErrorRequest(w, r, errors.New("no proxies found"), http.StatusInternalServerError)
		return
	}

	responses.SendJson(response{Proxies: proxiesArray}, http.StatusOK, w, r)

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
