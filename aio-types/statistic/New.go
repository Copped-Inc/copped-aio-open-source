package statistic

import (
	"bytes"
	"context"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
	"time"
)

func New(r *http.Request) *http.Request {

	if !helper.RequestLog {
		return r
	}

	var body []byte
	if r.Body != nil {
		if strings.Contains(r.URL.Path, "billing") || strings.Contains(r.URL.Path, "shipping") {
			body = []byte("confidential")
		} else {
			body, _ = io.ReadAll(r.Body)
			r.Body.Close()
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}
	}

	id, err := helper.GetClaim("id", r)
	if err != nil {
		if helper.IsMaster(r.Header.Get("Password")) {
			id = "internal"
		} else {
			id = "unknown"
		}
	}

	*r = *r.WithContext(context.WithValue(r.Context(), "statistic", statistic{
		Id:          uuid.New().String(),
		Path:        r.URL.Path,
		Method:      r.Method,
		Start:       time.Now(),
		RequestBody: body,
		User:        id.(string),
	}))

	go func() {
		_ = <-r.Context().Done()
		stat := r.Context().Value("statistic").(statistic)

		if stat.StatusCode >= 500 {
			stat.ClosingState = stateErr
		} else if stat.StatusCode != 0 {
			stat.ClosingState = stateOk
		} else {
			stat.ClosingState = stateTimeout
		}

		err := stat.SaveRequest()
		if err != nil {
			panic(err)
		}
	}()

	return r

}
