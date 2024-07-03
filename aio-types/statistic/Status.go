package statistic

import (
	"context"
	"github.com/Copped-Inc/aio-types/helper"
	"net/http"
)

func Status(r *http.Request, status int) {

	if !helper.RequestLog {
		return
	}

	stat := r.Context().Value("statistic").(statistic)
	stat.StatusCode = status
	*r = *r.WithContext(context.WithValue(r.Context(), "statistic", stat))

}
