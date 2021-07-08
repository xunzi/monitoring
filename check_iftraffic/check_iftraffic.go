package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olorin/nagiosplugin"
)

const NETDEV = "/proc/net/dev"

type statsWithTimestamp struct {
	Tstamp   time.Time
	Statsmap map[string]int
}

var iface = flag.String("iface", "", "interface to check")
var stat = flag.String("stat", "", "stat to check [rxBytes, rxPackets, rxErrs rxDrop, rxFifo, rxFrame, rxCompressed  rxMulticast, txBytes, txPackets, txErrs, txDrop, txFifo, txColls txCarrier, txCompressed]")
var cacheFile = flag.String("cacheFile", "/var/tmp/check_iftraffic_cache.json", "cache file to save values from last run")
var warning = flag.Float64("warning", 0, "Warning")
var critical = flag.Float64("critical", 0, "Critical")

//reads NETDEV and returns stats for iface
func readNetdev(iface string) ([]string, bool) {
	fileHandle, err := os.Open(NETDEV)
	if err != nil {
		log.Fatal(err)
	}
	defer fileHandle.Close()
	scanner := bufio.NewScanner(fileHandle)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, iface) {
			return strings.Fields(line), true
		} else {
			continue
		}
	}
	// return empty string slice and false if if not found
	return []string{}, false
}

func Stat2Map(statsLine []string) map[string]int {
	statMap := make(map[string]int)
	statMap["rxBytes"], _ = strconv.Atoi(statsLine[1])
	statMap["rxPackets"], _ = strconv.Atoi(statsLine[2])
	statMap["rxErrs"], _ = strconv.Atoi(statsLine[3])
	statMap["rxDrop"], _ = strconv.Atoi(statsLine[4])
	statMap["rxFifo"], _ = strconv.Atoi(statsLine[5])
	statMap["rxFrame"], _ = strconv.Atoi(statsLine[6])
	statMap["rxCompressed"], _ = strconv.Atoi(statsLine[7])
	statMap["rxMulticast"], _ = strconv.Atoi(statsLine[8])
	statMap["txBytes"], _ = strconv.Atoi(statsLine[9])
	statMap["txPackets"], _ = strconv.Atoi(statsLine[10])
	statMap["txErrs"], _ = strconv.Atoi(statsLine[11])
	statMap["txDrop"], _ = strconv.Atoi(statsLine[12])
	statMap["txFifo"], _ = strconv.Atoi(statsLine[13])
	statMap["txFrame"], _ = strconv.Atoi(statsLine[14])
	statMap["txCompressed"], _ = strconv.Atoi(statsLine[15])

	return statMap
}

//writes json data to file, returns Error ok
func saveStats2Json(statsMap map[string]int) {
	jsonStats := statsWithTimestamp{}
	jsonStats.Tstamp = time.Now()
	jsonStats.Statsmap = statsMap
	jsonMarshalled, err := json.Marshal(jsonStats)
	if err != nil {
		log.Fatalf("Error marshalling %v", jsonStats.Statsmap)
	}
	ok := ioutil.WriteFile(*cacheFile, jsonMarshalled, 0644)
	if ok != nil {
		log.Fatalf("Could not write json data to file %s: %s", *cacheFile, ok.Error())
	}
}

func main() {
	flag.Parse()
	check := nagiosplugin.NewCheck()
	defer check.Finish()
	ifStats, ok := readNetdev(*iface)
	if !ok {
		check.AddResult(nagiosplugin.UNKNOWN, fmt.Sprintf("Interface %s not found", *iface))
	}
	statsMap := Stat2Map(ifStats)
	saveStats2Json(statsMap)
	statName := fmt.Sprintf("%s-%s", *iface, *stat)
	statValue := float64(statsMap[*stat])
	check.AddPerfDatum(statName, "", statValue, 0.0, math.Inf(1), *warning, *critical)
	switch {
	case statValue < *warning:
		check.AddResult(nagiosplugin.OK, fmt.Sprintf("Interface stats %s ok", statName))
	case statValue > *warning && statValue < *critical:
		check.AddResult(nagiosplugin.WARNING, fmt.Sprintf("Interface stats %s warning", statName))
	case statValue > *critical:
		check.AddResult(nagiosplugin.CRITICAL, fmt.Sprintf("Interface stats %s critical", statName))
	default:
		check.AddResult(nagiosplugin.UNKNOWN, fmt.Sprintf("Interface stats %s unknown", statName))
	}
	if statValue < *warning {

	}
	//fmt.Println(statsMap[*stat])
}
