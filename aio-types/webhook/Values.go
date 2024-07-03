package webhook

import "github.com/infinitare/disgo"

type webhook int

const (
	Test webhook = iota + 1
	DataRequest
	UpdateClient
	UpdatePayments
	DataReceived
	NewProduct
	Log
	NewCheckout
	NewCheckoutLink
	PingFailed
	PingSuccess
	Whitelist
	MonitorDisabled
	MonitorRestarted
	ErrorLog
	InstanceLogout
	AiPredictIntra
	AiPredictDiff
	ISINList
)

type Body struct {
	Username  string        `json:"username"`
	AvatarUrl string        `json:"avatar_url"`
	Embeds    []disgo.Embed `json:"embeds"`
}
