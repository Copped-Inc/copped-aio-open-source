package statistic

import (
	"github.com/Copped-Inc/aio-types/statistic/realtimedb"
	"strconv"
	"strings"
	"time"
)

func SetOffline(server string) error {

	year, month, day := time.Now().Date()
	globaltotalRef := realtimedb.GetDatabase().NewRef("uptime/" + strings.ReplaceAll(server, ".", "-") + "/" + strconv.Itoa(year) + "-" + month.String() + "-" + strconv.Itoa(day) + "/status")
	return globaltotalRef.Set(realtimedb.GetContext(), 1)

}
