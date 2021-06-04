package alert

import (
	"github.com/eliezedeck/gobase/logging"
	"github.com/gregdel/pushover"
	"go.uber.org/zap"
)

var (
	defaultPushoverApp       *pushover.Pushover
	defaultPushoverRecipient *pushover.Recipient
)

func SetDefaultPushoverTokens(app, recipient string) {
	defaultPushoverApp = pushover.New(app)
	defaultPushoverRecipient = pushover.NewRecipient(recipient)
	SendAlertImplementation = defaultSendAlert
}

// defaultSendAlert is an implementation that sends the message via PushOver
func defaultSendAlert(title, message string) {
	if defaultPushoverApp == nil || defaultPushoverRecipient == nil {
		return
	}
	msg := pushover.NewMessageWithTitle(message, title)
	if _, err := defaultPushoverApp.SendMessage(msg, defaultPushoverRecipient); err != nil {
		logging.L.Error("Failed sending Alert", zap.Error(err))
	}
}
