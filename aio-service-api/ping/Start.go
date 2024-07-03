package ping

import (
	"github.com/Copped-Inc/aio-types/ping"
	"time"
)

func Start() {

	for {
		ping.Check("aio.copped-inc.com")
		ping.Check("database.copped-inc.com")
		ping.Check("monitor.copped-inc.com")
		ping.Check("instances.copped-inc.com")

		time.Sleep(time.Second * 10)
	}

}
