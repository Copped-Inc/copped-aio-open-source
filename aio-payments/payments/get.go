package payments

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
)

func getCtx() (context.Context, context.CancelFunc, context.CancelFunc) {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
	)

	allocCtx, c := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, ca := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	return ctx, c, ca

}
