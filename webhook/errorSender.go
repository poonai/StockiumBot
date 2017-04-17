package webhook

import "github.com/sch00lb0y/StockiumBot/fb"
import "github.com/Sirupsen/logrus"

func errorSend(senderID string, errMsg string) {

	err := fb.Send(senderID, errMsg)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ERROR": err.Error(),
		}).Warn("WEBHOOK ERROR HANDELER")
	}
}
