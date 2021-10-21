package telegramnotifier_test

import (
	telegramnnotifier "telegramnotifier"
	"testing"
)

func TestGenerateNotification(t *testing.T) {
	t.Parallel()
	var got string
	want := "PROBLEM:\nHost samplehost is down\n(<pre>not reachable</pre>)\n @2021-10-20:21:33"
	got = telegramnnotifier.GenerateNotification("host", "PROBLEM", "samplehost", "", "down", "not reachable", "2021-10-20:21:33")
	if want != got {
		t.Errorf("Want: %q, got %q", want, got)
	}
}
