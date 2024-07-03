package captcha

type response struct {
	Success                  bool   `json:"Success"`
	IsCaptcha                bool   `json:"IsCaptcha"`
	CaptchaSdkUrl            string `json:"CaptchaSdkUrl"`
	HtmlInject               string `json:"HtmlInject"`
	ResultCallbackMethodName string `json:"ResultCallbackMethodName"`
}
