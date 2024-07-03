package cookies

import (
	"net/http"
	"runtime"
	"time"
)

func Add(w http.ResponseWriter, key string, value string) {

	expire := time.Now().Add(72 * time.Hour)
	cookie := http.Cookie{
		Name:  key,
		Value: value,
		Domain: func() string {
			if runtime.GOOS != "windows" {
				return ".copped-inc.com"
			} else {
				return ".localhost"
			}
		}(),
		Path:    "/",
		Expires: expire,
	}
	http.SetCookie(w, &cookie)

}
