package api

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"
)

type WorkType struct {
	key   string `json:"key"`
	value string `json:"value"`
}

type TempoWorklogPostRequestBody struct {
	authorAccountId  string     `json:"authorAccountId"`
	description      string     `json:"description"`
	issueId          int        `json:"issueId"`
	startDate        string     `json:"startDate"`
	startTime        string     `json:"startTime"`
	timeSpentSeconds int        `json:"timeSpentSeconds"`
	attributes       []WorkType `json:"attributes"`
}

var DevWorkType = WorkType{
	key:   "_WorkType_",
	value: "14C",
}
var PtoWorkType = WorkType{
	key:   "_WorkType_",
	value: "20E",
}
var MeetingWorkType = WorkType{
	key:   "_WorkType_",
	value: "12E",
}
var SupportWorkType = WorkType{
	key:   "_WorkType_",
	value: "11E",
}

func SendWorklog(workType WorkType, hours int, day time.Time, bearerToken string, accountId string, issueId string) {
	for logForDay := 1; logForDay <= hours && logForDay <= 5; logForDay++ {
		var timeSpentHours int
		if hours < 5 {
			timeSpentHours = 1
		} else {
			timeSpentHours = hours/5 + (int(math.Mod(float64(hours), 5))+5-logForDay)/5
		}

		fmt.Print("Logging ", timeSpentHours, " hours for ")
		fmt.Print(day.AddDate(0, 0, logForDay).Format(time.DateOnly), " \n")

		tempoWorklogRequest := fmt.Sprintf("{ \"authorAccountId\":  \"%s\", \"description\":  \"dpctl tempo\", \"issueId\": %s,\"startDate\": \"%s\", \"timeSpentSeconds\": %d, \"startTime\": \"09:00:00\", \"attributes\": [ {\"key\":   \"_WorkType_\", \"value\": \"%s\" }]}", accountId, issueId, day.AddDate(0, 0, logForDay-1).Format(time.DateOnly), timeSpentHours*3600, workType.value)

		client := &http.Client{}
		req, err := http.NewRequest("POST", "https://api.tempo.io/4/worklogs", bytes.NewBuffer([]byte(tempoWorklogRequest)))
		if err != nil {
			log.Fatal(err)
			return
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

		resp, err := client.Do(req)

		if err != nil {
			log.Fatal(err)
			return
		} else if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			log.Fatal(string(bodyBytes))
			return
		}
	}
}
