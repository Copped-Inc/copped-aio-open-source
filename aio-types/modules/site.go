package modules

import (
	"math/big"

	"golang.org/x/exp/slices"
)

type (
	Site string

	siteData struct {
		Name            string
		Monitored       bool
		Runable         bool
		CaptchaRequired bool
		URL             string
	}
)

const (
	Kith_EU  Site = "kith_eu"
	Queue_it Site = "queue_it"
)

func (site Site) GetData() siteData {
	return sites[site]
}

func (site Site) Parse() *big.Int {
	return new(big.Int).Lsh(big.NewInt(1), uint(slices.Index(Sites, site)))
}

var (
	Sites = []Site{Kith_EU, Queue_it}

	sites = map[Site]siteData{
		Kith_EU: {
			Name:            "Kith EU",
			Monitored:       true,
			Runable:         true,
			CaptchaRequired: true,
			URL:             "https://eu.kith.com",
		},
		Queue_it: {
			Name:            "Queue-it",
			Monitored:       false,
			Runable:         false,
			CaptchaRequired: false,
			URL:             "https://queue-it.com",
		},
	}
)
