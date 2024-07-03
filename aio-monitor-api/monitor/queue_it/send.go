package queue_it

import (
	"monitor-api/api"
	"strings"
	"time"
)

func send(url string) {

	api.InstockReq{
		Name: strings.Split(url, "/")[2],
		Sku:  strings.Split(url, "/")[2],
		Skus: map[string]string{
			"OS": "OS",
		},
		Date:  time.Now(),
		Link:  url,
		Image: "https://cdn.discordapp.com/attachments/901123672809566258/1142381279359275038/317600223_637593141491709_1867642227124396667_n.jpg",
		Price: 0.0,
	}.Send()

}
