package alert

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/eliezedeck/gobase/logging"
	"go.uber.org/zap"
)

var (
	SendAlertImplementation   = defaultSendAlert
	SendAlertDebounceDuration = 60 * time.Minute
	SendAlertDisabled         = false
	NoAlertHosts              []string
)

var (
	tagsDebounceTracking      = make(map[string]time.Time, 4)
	tagsDebounceTrackingMutex = &sync.Mutex{}
)

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	cnh := strings.ToLower(strings.TrimSpace(hostname))
	for _, h := range NoAlertHosts {
		nh := strings.ToLower(strings.TrimSpace(h))
		if nh == cnh {
			SendAlertDisabled = true
			logging.L.Info("Alerting disabled", zap.String("hostname", hostname))
			break
		}
	}
}

// Send uses the specified SendAlertImplementation function to asynchronously send an alert. You should NOT start
// a goroutine when calling this.
//
// What is exactly done in order to deliver the message is really up to the implementation. So, it will allow you to do
// absolutely anything.
func Send(tag, title, message string, syncMode ...bool) {
	if SendAlertDisabled {
		return
	}

	tagsDebounceTrackingMutex.Lock()
	defer tagsDebounceTrackingMutex.Unlock()

	// Check if this tag has recently been used to send a message
	now := time.Now()
	if t, found := tagsDebounceTracking[tag]; tag != "" && found {
		if t.Add(SendAlertDebounceDuration).After(now) {
			return
		}
	}

	// Record the send time of this tag
	tagsDebounceTracking[tag] = now

	synchronous := false
	if len(syncMode) > 0 && syncMode[0] {
		synchronous = true
	}

	if SendAlertImplementation != nil {
		if synchronous {
			SendAlertImplementation(title, message)
		} else {
			go SendAlertImplementation(title, message)
		}
	} else {
		if synchronous {
			defaultSendAlert(title, message)
		} else {
			go defaultSendAlert(title, message)
		}
	}
}
