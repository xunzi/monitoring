package main

import (
	"flag"
	"fmt"
	"ncperflib"
	"time"
)

var hostname = flag.String("hostname", "", "hostname of nextcloud instance")
var uri = flag.String("uri", "/ocs/v2.php/apps/serverinfo/api/v1/info", "URI containing the status info")
var username = flag.String("username", "", "Nextcloud user name (admin permission reqd")
var password = flag.String("password", "", "Password to authenticate against nextcloud")
var counter = flag.String("counter", "", "Counter to be monitored [AppUdatesAvailable|FreeSpace|NumShares|ActiveUsers5Min|DbSize]")
var critical = flag.Int64("critical", 0, "Critical Value")
var warning = flag.Int64("warning", 0, "Warning Value")
var debug = flag.Bool("debug", false, "show debugging output")
var perfdata = flag.Bool("perfdata", false, "output perfdata")

func main() {

	flag.Parse()
	ncperflib.Debug = *debug
	ncperflib.CheckArguments(*counter, *warning, *critical)
	startTime := time.Now()
	var perfData = ncperflib.FetchPerformanceInfo(*counter, *hostname, *uri, *username, *password)
	var perfInfo = ncperflib.ParsePerfData(perfData, *counter)
	endTime := time.Now()
	runtime := endTime.Sub(startTime)
	result := fmt.Sprintf("%s: %s", *counter, fmt.Sprintf("%d", perfInfo))
	if *perfdata {
		result = fmt.Sprintf("%s | %s=%d,runtime=%s", result, *counter, perfInfo, runtime)
	}
	if perfInfo == -1 {
		ncperflib.NagiosResult(3, fmt.Sprintf("Unknown value for %s", *counter))
	}
	if perfInfo < *warning {
		ncperflib.NagiosResult(0, result)
	}
	if perfInfo >= *warning {
		ncperflib.NagiosResult(1, result)
	}
	if perfInfo >= *critical {
		ncperflib.NagiosResult(2, result)
	}
	//Debugprint(fmt.Sprintf("Total runtime: %s", runtime))
}
