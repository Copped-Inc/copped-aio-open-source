package payments

import (
	"context"
	"github.com/chromedp/chromedp"
	"strings"
	"time"
)

func waitURL(url string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		var currentURL string
		err := chromedp.EvaluateAsDevTools(`window.location.href`, &currentURL).Do(ctx)
		if err != nil {
			return err
		}
		for !strings.Contains(currentURL, url) {
			err := chromedp.EvaluateAsDevTools(`window.location.href`, &currentURL).Do(ctx)
			if err != nil {
				return err
			}
			time.Sleep(100 * time.Millisecond)
		}
		return nil
	})
}
