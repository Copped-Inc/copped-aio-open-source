package helper

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func GetIP(w http.ResponseWriter, r *http.Request) (string, error) {

	var userIP string
	if len(r.Header.Get("CF-Connecting-IP")) > 1 {
		userIP = r.Header.Get("CF-Connecting-IP")
	} else if len(r.Header.Get("X-Forwarded-For")) > 1 {
		userIP = r.Header.Get("X-Forwarded-For")
	} else if len(r.Header.Get("X-Real-IP")) > 1 {
		userIP = r.Header.Get("X-Real-IP")
	} else {
		userIP = r.RemoteAddr
		if strings.Contains(userIP, ":") {
			userIP = strings.Split(userIP, ":")[0]
		}
	}

	if userIP == "" {
		return "", fmt.Errorf("could not get ip")
	}

	ip := net.ParseIP(userIP)
	if ip == nil {
		return "", fmt.Errorf("could not get ip")
	}

	return ip.String(), nil

}
