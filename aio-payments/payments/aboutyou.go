package payments

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

func (p *Payments) aboutyou() (err error) {

	println("Handling aboutyou payment", p.Id)
	ctx, c, ca := getCtx()
	defer c()
	defer ca()

	user := fmt.Sprintf("{\"isExistingCustomer\":true,\"recentlySeenProductIds\":[],\"sessionToken\":\"%s\",\"checkoutSessionData\":{\"cookieValue\":\"%s\",\"secret\":\"%s\"}}", p.data["auth"], p.data["session"], p.data["secret"])

	err = chromedp.Run(ctx,
		chromedp.Navigate("https://en.aboutyou.de/"),
		chromedp.WaitReady("body"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, exp, err := runtime.Evaluate("window.localStorage.setItem('user', '" + user + "')").Do(ctx)
			if err != nil {
				return err
			}
			if exp != nil {
				return exp
			}
			return nil
		}),
		setCookie("checkout_sid", p.data["session"], ".aboutyou.de"),
		chromedp.Reload(),
		chromedp.WaitVisible("//*[@id=\"react-root\"]/header/div[2]/section/div[3]/div/div/div[3]/span[2]", chromedp.BySearch),
		chromedp.Navigate("https://en.aboutyou.de/checkout"),
		waitURL("payment"),
		chromedp.Navigate("https://www.paypal.com/agreements/approve?ba_token="+p.data["paypal"]+"&useraction=commit"),
		waitURL("checkout/success"),
	)

	return err

}
