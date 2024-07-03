package captcha

import (
	"errors"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/proxies"
	"github.com/Copped-Inc/aio-types/responses"
	"net/http"
	"time"

	"github.com/Copped-Inc/aio-types/modules"
	"github.com/gorilla/mux"
)

const (
	kithEUCheckout = "c1e706a3-692c-4d0b-a555-f0bd298ee8ca"
	threads        = 5
)

func get(w http.ResponseWriter, r *http.Request) {

	s := mux.Vars(r)["site"]
	switch s {
	case string(modules.Kith_EU):
		kith_eu(w, r)
	default:
		console.ErrorRequest(w, r, errors.New("no site given"), http.StatusBadRequest)
	}

}

func kith_eu(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	for {
		if time.Since(start) > time.Second*55 {
			console.ErrorRequest(w, r, errors.New("not solved"), http.StatusRequestTimeout)
			return
		}

		for i, c := range queue.get() {
			for i := 0; i < len(c); i++ {
				if c[i].Expire.Before(time.Now()) {
					c = append(c[:i], c[i+1:]...)
					i--
				}
			}
			queue.set(i, c)
		}

		if len(queue.getSite(helper.KithEUSitekey)) > 0 && len(queue.getSite(kithEUCheckout)) > 0 {
			res := captchaRes{
				Token: []string{
					queue.getSite(helper.KithEUSitekey)[0].Token,
					queue.getSite(kithEUCheckout)[0].Token,
				},
				Expire: func() time.Time {
					if queue.getSite(helper.KithEUSitekey)[0].Expire.Before(queue.getSite(kithEUCheckout)[0].Expire) {
						return queue.getSite(helper.KithEUSitekey)[0].Expire
					}
					return queue.getSite(kithEUCheckout)[0].Expire
				}(),
			}

			queue.set(helper.KithEUSitekey, queue.getSite(helper.KithEUSitekey)[1:])
			queue.set(kithEUCheckout, queue.getSite(kithEUCheckout)[1:])

			responses.SendJson(res, http.StatusOK, w, r)
			return
		}

		_ = proxies.Dcs()
		if runningSolver.len(helper.KithEUSitekey) < threads {
			requestCaptcha(helper.KithEUSitekey)
		}

		if runningSolver.len(kithEUCheckout) < threads {
			requestCaptcha(kithEUCheckout)
		}

		time.Sleep(time.Millisecond * 300)
	}

}

type captchaRes struct {
	Token  []string  `json:"token"`
	Expire time.Time `json:"expire"`
}

type Captcha struct {
	Token  string    `json:"token"`
	Expire time.Time `json:"expire"`
}
