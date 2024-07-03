package captcha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/proxies"
	"net/http"
	"time"
)

func requestCaptcha(s string) {

	clientKey := "" // Insert CapSolver URL here
	baseUrl := "https://api.capsolver.com"

	proxy := proxies.Dcs()

	c := createTask{
		ClientKey: clientKey,
		Task: task{
			Type:        "HCaptchaTurboTask",
			WebsiteURL:  "https://eu.kith.com/pages/international-checkout",
			Proxy:       fmt.Sprintf("http:%s:%s:%s:%s", proxy.Ip, proxy.Port, proxy.Username, proxy.Password),
			WebsiteKey:  s,
			UserAgent:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36",
			IsInvisible: false,
		},
		CallBackUrl: helper.ActiveData + "/captcha/" + s,
	}

	j, err := json.Marshal(c)
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/createTask", bytes.NewBuffer(j))
	if err != nil {
		console.Log("Error", err.Error())
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		console.Log("Error", err.Error())
		return
	}

	var createTaskResp createTaskResponse
	err = json.NewDecoder(res.Body).Decode(&createTaskResp)
	if err != nil {
		console.Log("Error", err.Error())
		return
	}

	if createTaskResp.ErrorId == 1 {
		return
	}

	runningSolver.add(s, createTaskResp.TaskId)
	go func() {
		time.Sleep(time.Minute)
		runningSolver.remove(s, createTaskResp.TaskId)
	}()

}

type task struct {
	Type        string `json:"type"`
	WebsiteURL  string `json:"websiteURL"`
	WebsiteKey  string `json:"websiteKey"`
	Proxy       string `json:"proxy"`
	UserAgent   string `json:"userAgent"`
	IsInvisible bool   `json:"isInvisible"`
}

type createTask struct {
	ClientKey   string `json:"clientKey"`
	Task        task   `json:"task"`
	CallBackUrl string `json:"callBackUrl"`
}

type createTaskResponse struct {
	ErrorId int    `json:"errorId"`
	TaskId  string `json:"taskId"`
	Status  string `json:"status"`
}
