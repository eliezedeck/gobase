package alert

import (
	"github.com/eliezedeck/gobase/logging"
	"github.com/gregdel/pushover"
	"go.uber.org/zap"
)

var (
	DefaultPushoverApp       = pushover.New("antz3gidwyf2hmqovwsip7q36hzise")
	DefaultPushoverRecipient = pushover.NewRecipient("u1ckbx8wdfhkgi9k3nj1w2baoni47z")
)

func SetDefaultPushoverTokens(app, recipient string) {
	DefaultPushoverApp = pushover.New(app)
	DefaultPushoverRecipient = pushover.NewRecipient(recipient)
}

// defaultSendAlert is an implementation that sends the message via PushOver
func defaultSendAlert(title, message string) {
	msg := pushover.NewMessageWithTitle(message, title)
	if _, err := DefaultPushoverApp.SendMessage(msg, DefaultPushoverRecipient); err != nil {
		logging.L.Error("Failed sending Alert", zap.Error(err))
	}
}
