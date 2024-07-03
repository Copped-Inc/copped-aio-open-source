package statistic

import (
	"context"
	"github.com/Copped-Inc/aio-types/helper"
	"net/http"
)

func User(r *http.Request, user string) {

	if !helper.RequestLog {
		return
	}

	stat := r.Context().Value("statistic").(statistic)
	if stat.User == "" {
		stat.User = user
	}

	*r = *r.WithContext(context.WithValue(r.Context(), "statistic", stat))

}
