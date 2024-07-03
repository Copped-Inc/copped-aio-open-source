package statistic

import (
	"context"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/google/uuid"
	"net/http"
)

func AddError(r *http.Request, err error) {

	if !helper.RequestLog {
		return
	}

	stat := r.Context().Value("statistic").(statistic)
	stat.Err = err.Error()

	*r = *r.WithContext(context.WithValue(r.Context(), "statistic", stat))

}

func AddLog(r *http.Request, text ...any) {

	if !helper.RequestLog {
		go SaveLog(nil, uuid.New().String(), text...)
		return
	}

	stat := r.Context().Value("statistic").(statistic)
	id := uuid.New().String()
	stat.Logs = append(stat.Logs, id)

	*r = *r.WithContext(context.WithValue(r.Context(), "statistic", stat))
	go SaveLog(r, id, text...)

}
