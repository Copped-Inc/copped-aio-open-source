package restart

import (
	"errors"
	"monitor-api/monitor"
	"monitor-api/monitor/captcha"
	"monitor-api/monitor/kith_eu"
	"monitor-api/monitor/queue_it"
	"monitor-api/monitor/service_api"
	"monitor-api/monitor/traiding"
	"net/http"
	"time"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"
	"github.com/Copped-Inc/aio-types/webhook"
	"github.com/Copped-Inc/aio-types/worker"
	"github.com/gorilla/mux"
)

func post(w http.ResponseWriter, r *http.Request) {

	m := Type(mux.Vars(r)["site"])

	if Running[m] {
		console.ErrorRequest(w, r, errors.New("monitor already running"), http.StatusNotAcceptable)
		return
	}

	limiter := Limiter(m)

	wh := webhook.New()
	switch m {
	case Kith_EU_Captcha:
		_ = worker.Worker(monitor.Loop(captcha.Start, time.Minute*5), limiter)
	case Kith_EU:
		_ = worker.Worker(monitor.Loop(kith_eu.Start, time.Second*2), limiter)
	case Queue_it_VfB:
		limiter.Timeout = time.Minute * 5
		limiter.Interval = time.Minute * 20
		_ = worker.Worker(monitor.Loop(queue_it.Start("https://shop.vfb.de/"), time.Minute*5), limiter)
	case Service_API:
		_ = worker.Worker(monitor.Loop(service_api.Start, time.Second*10), limiter)
	case Trading_Tesla:
		limiter.Timeout = time.Minute * 1
		_ = worker.Worker(monitor.Loop(traiding.Start("TL0.F"), time.Minute*5), limiter)
	default:
		console.ErrorRequest(w, r, errors.New("invalid monitor"), http.StatusBadRequest)
		return
	}

	wh.AddEmbed(
		webhook.MonitorRestarted,
		string(m),
	)

	Running[m] = true
	_ = wh.Send("") // Insert Webhook URL here

	responses.SendOk(w, r)

}

type Type string

const (
	Kith_EU_Captcha Type = "kith_eu_captcha"
	Kith_EU         Type = "kith_eu"
	Queue_it_VfB    Type = "queue_it_vfb"
	Service_API     Type = "service_api"
	Trading_Tesla   Type = "trading_tesla"
)

var Running = map[Type]bool{
	Kith_EU_Captcha: true,
	Kith_EU:         true,
	Service_API:     true,
	Trading_Tesla:   true,
}

type Monitor struct {
	Name    string
	Running bool
}

func Limiter(t Type) worker.Limiter {
	return worker.Limiter{
		Interval: time.Minute * 10,
		Timeout:  time.Minute * 3,
		Limit:    3,
		Handler: func(err error) {
			Running[t] = false

			wh := webhook.New()
			wh.AddEmbed(
				webhook.MonitorDisabled,
				string(t),
				err.Error(),
			)

			_ = wh.Send("") // Insert Webhook URL here
		},
	}
}
