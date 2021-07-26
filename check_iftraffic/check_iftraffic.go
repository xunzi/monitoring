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

var stats = []string{"rxBytes", "txBytes", "rxErrs", "txErrs", "rxDrop", "txDrop", "bytes", "packets", "errs", "drop"}

var iface = flag.String("iface", "", "interface to check")
var stat = flag.String("stat", "", fmt.Sprintf("stat to check %s", stats))
var cacheFile = flag.String("cacheFile", "/var/tmp/check_iftraffic_cache.json", "cache file to save values from last run")
var warning = flag.Float64("warning", 0, "Warning")
var critical = flag.Float64("critical", 0, "Critical")
var perfdata = flag.Bool("perfdata", false, "output perfdata")

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
	statMap["packets"] = statMap["rxPackets"] + statMap["txPackets"]
	statMap["bytes"] = statMap["rxBytes"] + statMap["txBytes"]
	statMap["errs"] = statMap["rxErrs"] + statMap["txErrs"]
	statMap["drop"] = statMap["rxDrop"] + statMap["txDrop"]
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

func readStatsfile() statsWithTimestamp {
	var statsWithTimestamp statsWithTimestamp
	jsonFile, err := os.Open(*cacheFile)
	if err != nil {
		log.Fatalf("Could not open file %s: %v", NETDEV, err.Error())
	}
	defer jsonFile.Close()
	bArray, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(bArray, &statsWithTimestamp)
	return statsWithTimestamp
}

func main() {
	flag.Parse()
	check := nagiosplugin.NewCheck()
	defer check.Finish()
	ifStats, ok := readNetdev(*iface)
	if !ok {
		check.AddResult(nagiosplugin.UNKNOWN, fmt.Sprintf("Interface %s not found", *iface))
	}
	//read old data
	oldStatsMap := readStatsfile()
	statsMap := Stat2Map(ifStats)
	saveStats2Json(statsMap)
	oldStatValue := float64(oldStatsMap.Statsmap[*stat])
	oldStatTstamp := oldStatsMap.Tstamp
	statName := fmt.Sprintf("%s-%s", *iface, *stat)
	statValue := float64(statsMap[*stat]) - oldStatValue
	timeDiff := time.Since(oldStatTstamp) / time.Second
	if *perfdata {
		check.AddPerfDatum(statName, "", statValue, 0.0, math.Inf(1), *warning, *critical)
	}
	outputSuffix := fmt.Sprintf("(%d %s in %d seconds)", int(statValue), *stat, timeDiff)
	switch {
	case statValue < *warning:
		check.AddResult(nagiosplugin.OK, fmt.Sprintf("Interface stats %s ok, %s", statName, outputSuffix))
	case statValue > *warning && statValue < *critical:
		check.AddResult(nagiosplugin.WARNING, fmt.Sprintf("Interface stats %s warning, %s", statName, outputSuffix))
	case statValue > *critical:
		check.AddResult(nagiosplugin.CRITICAL, fmt.Sprintf("Interface stats %s critical, %s", statName, outputSuffix))
	default:
		check.AddResult(nagiosplugin.UNKNOWN, fmt.Sprintf("Interface stats %s unknown", statName))
	}

	//fmt.Printf("%v\n", oldStatTstamp)
}
