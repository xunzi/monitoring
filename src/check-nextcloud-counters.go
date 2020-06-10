package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/olorin/nagiosplugin"
)

var hostname = flag.String("hostname", "", "hostname of nextcloud instance")
var uri = flag.String("uri", "/ocs/v2.php/apps/serverinfo/api/v1/info", "URI containing the status info")
var username = flag.String("username", "", "Nextcloud user name (admin permission reqd")
var password = flag.String("password", "", "Password to authenticate against nextcloud")
var critical = flag.Int64("critical", 0, "Critical Value")
var warning = flag.Int64("warning", 0, "Warning Value")

func fetchPerformaceInfo() []byte {
	perfURL := fmt.Sprintf("https://%s/%s?format=json", *hostname, *uri)
	req, err := http.NewRequest("GET", perfURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(*username, *password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return b
}

func main() {
	flag.Parse()
	check := nagiosplugin.NewCheck()
	defer check.Finish()
	check.AddResult(nagiosplugin.OK, "everything looks shiny, cap'n")
	perfInfo := fetchPerformaceInfo()
	fmt.Println(string(perfInfo))
}
