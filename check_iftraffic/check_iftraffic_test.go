package main

import (
	"strings"
	"testing"
)

func TestParseNetDev(t *testing.T) {
	sampleOutput := "wlp2s0: 100154002   76522    0 3341    0     0          0         0  3220546   26904    0    0    0     0       0          0"
	//statsLine := strings.Trim(sampleOutput)
	statsLine := strings.Fields(sampleOutput)
	statsMap := Stat2Map(statsLine)
	if statsMap["rxBytes"] != 100154002 {
		t.Errorf("rxBytes: expected 100154002, got %d", statsMap["rxBytes"])
	}
}
