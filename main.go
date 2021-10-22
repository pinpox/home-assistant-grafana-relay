package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Long-lived access token from home-assistant
var AUTH_TOKEN string

// Home-assistant api uri
var HM_SERVICE_URI string

// Port to listen for webhooks from grafana
var LISTEN_PORT string

// Host/ip to listen for webhooks from grafana
var LISTEN_HOST string

func init() {
	AUTH_TOKEN = os.Getenv("AUTH_TOKEN")
	HM_SERVICE_URI = os.Getenv("HM_SERVICE_URI")
	LISTEN_PORT = os.Getenv("LISTEN_PORT")
	LISTEN_HOST = os.Getenv("LISTEN_HOST")
}

type GrafanaJson struct {
	Title       string   `json:"title"`
	RuleID      int64    `json:"ruleId"`
	RuleName    string   `json:"ruleName"`
	State       string   `json:"state"`
	RuleURL     string   `json:"ruleUrl"`
	ImageURL    string   `json:"imageUrl"`
	Message     string   `json:"message"`
	OrgID       int      `json:"orgId"`
	DashboardID int      `json:"dashboardId"`
	PanelID     int      `json:"panelId"`
	Tags        struct{} `json:"tags"`
	EvalMatches []struct {
		Value  int         `json:"value"`
		Metric string      `json:"metric"`
		Tags   interface{} `json:"tags"`
	} `json:"evalMatches"`
}

func receiveHook(rw http.ResponseWriter, req *http.Request) {

	// Parse webhook data
	var j GrafanaJson
	body, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(body, &j)

	// send notification
	log.Println(string(body))
	notify(j)
}

func notify(hookData GrafanaJson) {

	// Encode the data
	postBody, err := json.Marshal(map[string]interface{}{
		"message": hookData.Message,
		"title":   hookData.Title,
		"data": map[string]string{
			"image": hookData.ImageURL,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	request, err := http.NewRequest("POST", HM_SERVICE_URI, bytes.NewBuffer(postBody))
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+AUTH_TOKEN)

	// Don't care about the response
	_, err = client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	log.Println("Listening for webooks on: " + LISTEN_HOST + ":" + LISTEN_PORT)
	log.Println("Using home-assistant at: " + HM_SERVICE_URI)
	http.HandleFunc("/", receiveHook)
	log.Fatal(http.ListenAndServe(LISTEN_HOST+":"+LISTEN_PORT, nil))
}
