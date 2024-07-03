package helper

import (
	"github.com/Copped-Inc/aio-types/statistic/realtimedb"
	"runtime"
)

func Set(server string, opt ...any) {

	if runtime.GOOS == "windows" && len(opt) == 0 {
		Active = Localhost
		ActiveCookie = LocalCookie
		ActiveData = LocalData
		ActiveService = LocalService
		ActiveInstances = LocalInstances
		System = "windows"
	}

	var err error
	switch server {
	case "aio.copped-inc.com":
		Webhook = "" // Insert Webhook URL here
		err = realtimedb.Init(Active)
	case "database.copped-inc.com":
		RequestLog = true
		Webhook = "" // Insert Webhook URL here
		err = realtimedb.Init(ActiveData)
	case "monitor.copped-inc.com":
		Webhook = "" // Insert Webhook URL here
		err = realtimedb.Init("https://monitor.copped-inc.com")
	case "service.copped-inc.com":
		Webhook = "" // Insert Webhook URL here
		err = realtimedb.Init(ActiveService)
	case "instances.copped-inc.com":
		Webhook = "" // Insert Webhook URL here
		err = realtimedb.Init(ActiveInstances)
	}

	if err != nil {
		panic(err)
	}

	Server = server

}
