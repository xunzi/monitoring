package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const NETDEV = "/proc/net/dev"

var iface = flag.String("iface", "", "interface to check")
var stat = flag.String("stat", "", "stat to check [rxBytes, rxPackets, rxErrs rxDrop, rxFifo, rxFrame, rxCompressed  rxMulticast, txBytes, txPackets, txErrs, txDrop, txFifo, txColls txCarrier, txCompressed]")

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

func main() {
	flag.Parse()
	ifStats, ok := readNetdev(*iface)
	if !ok {
		log.Fatal("Did not find interface ", *iface)
	}
	statsMap := Stat2Map(ifStats)
	fmt.Println(statsMap[*stat])
}
