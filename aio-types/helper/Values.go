package helper

const (
	Localhost      = "http://localhost:90"
	webhost        = "https://aio.copped-inc.com"
	LocalData      = "http://localhost:91"
	webData        = "https://database.copped-inc.com"
	LocalService   = "http://localhost:93"
	webService     = "https://service.copped-inc.com"
	LocalInstances = "http://localhost:94"
	webInstances   = "https://instances.copped-inc.com"
	LocalCookie    = ".localhost"
	webCookie      = ".copped-inc.com"
)

var Active = webhost
var ActiveCookie = webCookie
var ActiveData = webData
var ActiveService = webService
var ActiveInstances = webInstances

var System = "linux"
var RequestLog = false
var GeneralLog = true
var Server = ""
var Webhook = ""

var (
	KithEUSitekey = "b989d9e8-0d14-41a0-870f-97b5283ba67d"
)
