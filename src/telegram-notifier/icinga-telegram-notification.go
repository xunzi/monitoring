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
var objectType = flag.String("objectType", "service", "Object type (host or service)")
var chatID = flag.String("chatID", "", "Chat Id")
var notificationType = flag.String("notificationType", "", "Notification type")
var hostName = flag.String("hostName", "", "host name")
var state = flag.String("state", "", "State (up, down)")
var serviceName = flag.String("serviceName", "", "Service Name")
var outPut = flag.String("outPut", "", "Check command output")
var timeStamp = flag.String("timeStamp", "", "Event time stamp")
var debug = flag.Bool("debug", false, "Debug/verbose mode")
var notificationMessage string

func sendNotification(objectType string, notificationType string, hostName string, serviceName string, state string, outPut string, timeStamp string) {
	botURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", *botToken)
	if objectType == "host" {
		notificationMessage = fmt.Sprintf("%s:\nHost %s is %s\n(<pre>%s</pre>)\n @%s", notificationType, hostName, state, outPut, timeStamp)
	} else {
		notificationMessage = fmt.Sprintf("%s:\n Service %s on host %s is %s\n(<pre>%s</pre>)\n@%s",
			notificationType, serviceName, hostName, state, outPut, timeStamp)
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
	sendNotification(*objectType, *notificationType, *hostName, *serviceName, *state, *outPut, *timeStamp)
}
