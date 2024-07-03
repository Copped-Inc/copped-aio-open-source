package request

import (
	"compress/flate"
	"compress/gzip"
	"crypto/tls"
	"encoding/base64"
	"github.com/Copped-Inc/aio-types/proxies"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func Get(u string, body io.Reader, headers map[string]string, stopRedirect ...any) (*http.Response, error) {

	proxy := func() proxies.Proxy {
		if strings.Contains(u, "stockx") {
			return proxies.Residential()
		} else {
			return proxies.Dcs()
		}
	}()

	proxyURL, _ := url.Parse("http://" + proxy.Ip + ":" + proxy.Port)
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(proxy.Username+":"+proxy.Password))

	hdr := http.Header{}
	hdr.Add("Proxy-Authorization", basicAuth)
	transport := &http.Transport{
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		Proxy:              http.ProxyURL(proxyURL),
		ProxyConnectHeader: hdr,
	}

	c := &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}

	if len(stopRedirect) != 0 {
		c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	req, err := http.NewRequest(http.MethodGet, u, body)
	if err != nil {
		return &http.Response{}, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.101 Safari/537.36")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Connection", "keep-alive")

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

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
