package api

import (
	"encoding/json"
	"github.com/Copped-Inc/aio-types/console"
)

func (p InstockReq) Send() {

	j, err := json.Marshal(p)
	if err != nil {
		console.ErrorLog(err)
		return
	}

	_, err = InternalPost(j, "monitor/instock")
	if err != nil {
		console.ErrorLog(err)
		return
	}

}
