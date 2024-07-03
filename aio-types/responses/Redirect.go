package responses

import (
	"github.com/Copped-Inc/aio-types/statistic"
	"net/http"
	"strconv"
)

func Redirect(w http.ResponseWriter, r *http.Request, url string) {

	http.Redirect(w, r, url, http.StatusFound)
	statistic.Response(r, strconv.Itoa(http.StatusFound)+" "+http.StatusText(http.StatusFound)+" "+url)
	statistic.Status(r, http.StatusFound)

}
