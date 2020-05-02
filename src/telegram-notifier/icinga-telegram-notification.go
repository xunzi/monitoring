package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var botToken = flag.String("botToken", "", "Bot Token")
var notificationType = flag.String("type", "host", "notification type (host or service)")
var chatID = flag.String("chatID", "", "Chat Id")
var hostName = flag.String("hostName", "", "host name")
var state = flag.String("state", "", "State (up, down)")
var serviceName = flag.String("serviceName", "", "Service Name")
var outPut = flag.String("outPut", "", "Check command output")
var timeStamp = flag.String("timeStamp", "", "Event time stamp")
var debug = flag.Bool("debug", false, "Debug/verbose mode")
var notificationMessage string

func sendNotification(hostName string, serviceName string, state string, outPut string, timeStamp string) {
	botURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", *botToken)
	if serviceName != "" {
		msg := fmt.Sprintf("@%s: Service %s on host %s is %s (%s)", timeStamp, serviceName, hostName, state, outPut)
		notificationMessage = msg
	} else {
		msg := fmt.Sprintf("@%s: host %s is %s (%s)", timeStamp, hostName, state, outPut)
		notificationMessage = msg
	}
	requestBody, err := json.Marshal(map[string]string{
		"text":                     notificationMessage,
		"chat_id":                  *chatID,
		"parse_mode":               "HTML",
		"disable_web_page_preview": "True",
	})
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post(botURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	if *debug {
		log.Println(string(body))
	}
}

func main() {
	flag.Parse()
	hostName := *hostName
	state := *state
	serviceName := *serviceName
	outPut := *outPut
	timeStamp := *timeStamp
	sendNotification(hostName, serviceName, state, outPut, timeStamp)
}
