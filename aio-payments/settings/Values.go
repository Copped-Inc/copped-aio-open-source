package settings

type Settings struct {
	Authorization string  `json:"authorization"`
	Id            string  `json:"id"`
	Price         float64 `json:"price"`
	Provider      string  `json:"provider"`
	TaskMax       int     `json:"task_max"`
	Region        string  `json:"region"`
}

type request struct {
	Code string `json:"code"`
}

type response struct {
	Authorization string `json:"authorization"`
}
