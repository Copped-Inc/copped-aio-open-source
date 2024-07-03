package global

import (
	"github.com/Copped-Inc/aio-types/statistic/realtimedb"
)

func AddPing() {

	ref := realtimedb.GetDatabase().NewRef("/userstats/global/monitor/pings")
	var pings int
	err := ref.Get(realtimedb.GetContext(), &pings)
	if err != nil {
		return
	}

	err = ref.Set(realtimedb.GetContext(), pings+1)
	if err != nil {
		return
	}

}
