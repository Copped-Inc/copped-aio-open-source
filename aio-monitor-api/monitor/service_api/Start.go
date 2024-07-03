package service_api

import "github.com/Copped-Inc/aio-types/ping"

var Start = func() error {
	ping.Check("service.copped-inc.com")
	return nil
}
