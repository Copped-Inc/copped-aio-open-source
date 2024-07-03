package console

import (
	"fmt"
	"github.com/Copped-Inc/aio-types/statistic"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

func Log(text ...any) {

	log(text...)
	go statistic.SaveLog(nil, uuid.New().String(), text...)

}

func RequestLog(r *http.Request, text ...any) {

	log(text...)
	statistic.AddLog(r, text...)

}

func log(text ...any) {

	currentTime := time.Now().String()[:27]
	if len(strings.Split(currentTime, " ")) > 2 {

		currentTime = strings.Split(currentTime, " ")[0] + " " + strings.Split(currentTime, " ")[1]
		for len(currentTime) < 27 {
			currentTime = currentTime + "0"
		}

	}

	var parsedText []interface{}
	parsedText = append(parsedText, "["+currentTime+"]")
	for i, t := range text {
		if i != len(text)-1 {
			parsedText = append(parsedText, "["+fmt.Sprint(t)+"]")
		}
	}

	parsedText = append(parsedText, func() any {
		if len(text) > 0 {
			return text[len(text)-1]
		} else {
			return ""
		}
	}())

	fmt.Println(parsedText...)
	LoopText = append(LoopText, parsedText)

}
