package logging

import "github.com/sirupsen/logrus"

func LogSentryController()*logrus.Entry{
	logrus.WithField("Controller", "Sentry")
}
