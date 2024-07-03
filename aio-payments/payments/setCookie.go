package payments

import (
	"context"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func setCookie(name, value, domain string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		return network.SetCookie(name, value).WithDomain(domain).Do(ctx)
	})
}
