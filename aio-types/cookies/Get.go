package cookies

import "net/http"

func Get(r *http.Request, key string) *http.Cookie {

	cookies := r.Cookies()
	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == key {
			return cookies[i]
		}
	}
	return nil

}
