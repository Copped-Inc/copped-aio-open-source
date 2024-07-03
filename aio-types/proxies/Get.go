package proxies

import "math/rand"

func Residential() Proxy {

	p := residentials()
	proxy := p[rand.Intn(len(p))]
	return proxy

}

func Dcs() Proxy {

	p := dcs()
	proxy := p[rand.Intn(len(p))]
	return proxy

}
