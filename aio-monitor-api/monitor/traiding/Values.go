package traiding

type ApiData struct {
	Series map[string]Series `json:"Time Series (Daily)"`
}

type Series struct {
	Open  string `json:"1. open"`
	Close string `json:"4. close"`
}

type ExpData struct {
	Date      string  `json:"date"`
	DayWeek   float64 `json:"day_week"`
	DayMonth  float64 `json:"day_month"`
	Month     float64 `json:"month"`
	Year      float64 `json:"year"`
	FiveDiff  float64 `json:"five_diff"`
	FourDiff  float64 `json:"four_diff"`
	ThreeDiff float64 `json:"three_diff"`
	TwoDiff   float64 `json:"two_diff"`
	OneDiff   float64 `json:"one_diff"`
	Open      float64 `json:"open"`
}

type YahooData struct {
	Chart struct {
		Result []struct {
			Indicators struct {
				Quote []struct {
					Open  []float64 `json:"open"`
					Close []float64 `json:"close"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
	} `json:"chart"`
}
