package global

import (
	"github.com/Copped-Inc/aio-types/statistic/realtimedb"
)

func AddProduct() {

	ref := realtimedb.GetDatabase().NewRef("/userstats/global/monitor/products")
	var products int
	err := ref.Get(realtimedb.GetContext(), &products)
	if err != nil {
		return
	}

	err = ref.Set(realtimedb.GetContext(), products+1)
	if err != nil {
		return
	}

}
