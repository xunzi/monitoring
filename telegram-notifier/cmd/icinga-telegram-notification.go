package main

import (
	"flag"
	telegramnnotifier "telegramnotifier"
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

func main() {
	flag.Parse()
	telegramnnotifier.Debug = *debug
	notificationMessage = telegramnnotifier.GenerateNotification(*objectType, *notificationType, *hostName,
		*serviceName, *state, *outPut, *timeStamp)
	telegramnnotifier.SendNotification(*botToken, *chatID, notificationMessage)
}
