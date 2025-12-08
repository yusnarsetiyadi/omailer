package scheduler

import (
	"context"
	"omailer/pkg/whatsapp"
	"time"

	"github.com/sirupsen/logrus"
)

func WaitUntilWAReadyThenRun(name string, fn func()) {
	go func() {
		logrus.Infof("Waiting WhatsApp ready before running %s...", name)

		ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
		defer cancel()

		if err := whatsapp.WaitUntilReady(ctx); err != nil {
			logrus.Errorf("Failed waiting WA ready for %s: %v", name, err)
			return
		}

		logrus.Infof("WhatsApp ready → running %s", name)
		fn()
	}()
}

func AutomatedMessage(mainText string) string {
	footer := "_Automated message._\n_Please, buy me a coffee ☕_"
	return mainText + "\n\n\n" + footer
}
