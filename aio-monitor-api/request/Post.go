package request

import (
	"compress/flate"
	"compress/gzip"
	"encoding/base64"
	"github.com/Copped-Inc/aio-types/proxies"
	"io"
	"net/http"
	"net/url"
	"time"
)

func PostForm(u string, body io.Reader) (*http.Response, error) {

	proxy := proxies.Dcs()

	proxyURL, _ := url.Parse("http://" + proxy.Ip + ":" + proxy.Port)
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(proxy.Username+":"+proxy.Password))

	hdr := http.Header{}
	hdr.Add("Proxy-Authorization", basicAuth)
	transport := &http.Transport{
		TLSHandshakeTimeout:   time.Second * 10,
		ResponseHeaderTimeout: time.Second * 10,
		Proxy:                 http.ProxyURL(proxyURL),
		ProxyConnectHeader:    hdr,
	}

	c := &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}

	req, err := http.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return &http.Response{}, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.Do(req)
	if err != nil {
		return &http.Response{}, err
	}

	var reader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(res.Body)
		defer reader.Close()
		break
	case "deflate":
		reader = flate.NewReader(res.Body)
		defer reader.Close()
		break
	default:
		reader = res.Body
	}

	res.Body = reader
	return res, err
}
