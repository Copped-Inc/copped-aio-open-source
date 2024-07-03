package traiding

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/webhook"
	"monitor-api/request"
	"strconv"
	"time"
)

func Start(symbol string) func() error {
	return func() error {
		loc, _ := time.LoadLocation("Europe/Berlin")
		t := time.Now().In(loc)
		if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
			return nil
		}

		if !(t.Hour() == 8 && t.Minute() >= 51 && t.Minute() <= 55) {
			if !(t.Hour() == 21 && t.Minute() >= 51 && t.Minute() <= 55) {
				return nil
			}
		}

		res, err := request.Get("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol="+symbol+"&apikey=UCBSSQP29GLMLBCM&outputsize=full", nil, nil)
		if err != nil {
			return err
		}

		var data ApiData
		err = json.NewDecoder(res.Body).Decode(&data)
		if err != nil {
			return errors.New(fmt.Sprintf("%s error: %s", symbol, err.Error()))
		}

		res, err = request.Get("https://query1.finance.yahoo.com/v8/finance/chart/"+symbol, nil, nil)
		if err != nil {
			return err
		}

		var yData YahooData
		err = json.NewDecoder(res.Body).Decode(&yData)
		if err != nil {
			return errors.New(fmt.Sprintf("%s error: %s", symbol, err.Error()))
		}

		data.Series[t.Format("2006-01-02")] = Series{
			Open:  strconv.FormatFloat(yData.Chart.Result[0].Indicators.Quote[0].Open[0], 'f', 2, 64),
			Close: strconv.FormatFloat(yData.Chart.Result[0].Indicators.Quote[0].Close[len(yData.Chart.Result[0].Indicators.Quote[0].Close)-1], 'f', 2, 64),
		}

		var p, n float64
		var prompt string
		if t.Hour() == 8 {
			p, n, prompt, err = predictIntra(data.Series, t, symbol)
		} else {
			p, n, prompt, err = predictDiff(data.Series, t, symbol)
		}
		if err != nil {
			return err
		}

		console.Log(fmt.Sprintf("%s: pos. %f, neg. %f, input: %s", symbol, p, n, prompt))
		if t.Hour() == 8 {
			err = webhook.New().AddEmbed(webhook.AiPredictIntra,
				symbol,
				strconv.FormatFloat(yData.Chart.Result[0].Indicators.Quote[0].Open[0], 'f', 2, 64),
				func() string {
					if p >= n {
						return "📈"
					} else {
						return "📉"
					}
				}(),
			).AddEmbed(webhook.ISINList).Send("") // Insert Webhook URL here
		} else {
			err = webhook.New().AddEmbed(webhook.AiPredictDiff,
				symbol,
				strconv.FormatFloat(yData.Chart.Result[0].Indicators.Quote[0].Close[len(yData.Chart.Result[0].Indicators.Quote[0].Close)-1], 'f', 2, 64),
				func() string {
					if p >= n {
						return "📈"
					} else {
						return "📉"
					}
				}(),
			).AddEmbed(webhook.ISINList).Send("") // Insert Webhook URL here
		}
		return err

	}
}
