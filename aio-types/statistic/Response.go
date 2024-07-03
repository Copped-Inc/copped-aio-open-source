package statistic

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Copped-Inc/aio-types/helper"
	"net/http"
)

func Response(r *http.Request, b interface{}) {

	if !helper.RequestLog {
		return
	}

	stat := r.Context().Value("statistic").(statistic)
	switch b.(type) {
	case string:
		stat.ResponseBody = []byte(b.(string))
		break
	default:
		body, err := json.Marshal(b)
		if err != nil {
			AddError(r, err)
			break
		}

		if bytes.Contains(body, []byte("billing")) {
			stat.ResponseBody = []byte("confidential")
			break
		}

		stat.ResponseBody = body
		break
	}

	*r = *r.WithContext(context.WithValue(r.Context(), "statistic", stat))

}
