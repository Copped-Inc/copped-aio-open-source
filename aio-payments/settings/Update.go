package settings

import (
	"github.com/inconshreveable/go-update"
	"net/http"
	"os"
)

func (s *Settings) Update() error {

	_ = os.Remove(".payments.exe.old")
	req, err := http.NewRequest(http.MethodGet, "https://database.copped-inc.com/instance/update/payments", nil)
	if err != nil {
		return err
	}

	req.Header.Add("cookie", "authorization="+s.Authorization)
	req.Header.Add("version", "1.2.0")

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusOK {
		return err
	}
	defer res.Body.Close()

	err = update.Apply(res.Body, update.Options{})
	return err

}
