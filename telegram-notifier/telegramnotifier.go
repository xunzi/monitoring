package telegramnotifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var tgBotApiUrl = "https://api.telegram.org"

var Debug bool

func GenerateNotification(objectType string, notificationType string, hostName string, serviceName string, state string,
	outPut string, timeStamp string) (notificationMessage string) {
	//var notificationMessage string
	if objectType == "host" {
		notificationMessage = fmt.Sprintf("%s:\nHost %s is %s\n(<pre>%s</pre>)\n @%s",
			notificationType, hostName, state, outPut, timeStamp)
	} else {
		notificationMessage = fmt.Sprintf("%s:\n Service %s on host %s is %s\n(<pre>%s</pre>)\n@%s",
			notificationType, serviceName, hostName, state, outPut, timeStamp)
	}
	return notificationMessage
}

func SendNotification(token string, chat string, notificationMessage string) {

	botURL := fmt.Sprintf("%s/bot%s/sendMessage", tgBotApiUrl, token)
	if Debug {
		log.Printf("%s\n", botURL)
	}

	requestBody, err := json.Marshal(map[string]string{
		"text":                     notificationMessage,
		"chat_id":                  chat,
		"parse_mode":               "HTML",
		"disable_web_page_preview": "True",
	})
	if err != nil {
		log.Fatal(err)
	}

	if Debug {
		log.Printf("Notification message:\n%s\n", notificationMessage)
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
	if Debug {
		log.Println(string(body))
	}
}
