package main

import (
	"github.com/Copped-Inc/aio-types/worker"
	"math/rand"
	"monitor-api/handler"
	"monitor-api/handler/restart"
	"monitor-api/monitor"
	"monitor-api/monitor/captcha"
	"monitor-api/monitor/kith_eu"
	"monitor-api/monitor/queue_it"
	"monitor-api/monitor/service_api"
	"monitor-api/monitor/traiding"
	"net/http"
	"time"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/proxies"
)

var port = "92"

func main() {
	rand.Seed(time.Now().UnixNano())

	console.Log("Initialize", "Url", "monitor.copped-inc.com")
	console.Log("Initialize", "Port", port)
	helper.Set("monitor.copped-inc.com", 1)
	go console.Loop()

	proxies.Dcs()
	proxies.Residential()

	_ = worker.Worker(monitor.Loop(captcha.Start, time.Minute*5), restart.Limiter(restart.Kith_EU_Captcha))
	_ = worker.Worker(monitor.Loop(kith_eu.Start, time.Second*2), restart.Limiter(restart.Kith_EU))
	_ = worker.Worker(monitor.Loop(service_api.Start, time.Second*10), restart.Limiter(restart.Service_API))

	qLimiter := restart.Limiter(restart.Queue_it_VfB)
	qLimiter.Timeout = time.Minute * 5
	qLimiter.Interval = time.Minute * 20
	_ = worker.Worker(monitor.Loop(queue_it.Start("https://shop.vfb.de/"), time.Minute*5), qLimiter)

	tLimiter := restart.Limiter(restart.Trading_Tesla)
	tLimiter.Timeout = time.Minute * 1
	_ = worker.Worker(monitor.Loop(traiding.Start("TL0.F"), time.Minute*5), tLimiter)

	router := handler.Add()
	console.Log("Initialize", "Finished", "Listen and Serve")
	console.Log((&http.Server{Handler: router, Addr: ":" + port, WriteTimeout: 2 * time.Minute, ReadTimeout: 2 * time.Minute}).ListenAndServe())
}
