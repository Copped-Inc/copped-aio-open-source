package cookies

import (
	"net/http"
	"time"
)

func Remove(w http.ResponseWriter, key string) {

	c := &http.Cookie{
		Name:    key,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}

	http.SetCookie(w, c)

}
