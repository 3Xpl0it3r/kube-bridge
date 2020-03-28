package logging

import "github.com/sirupsen/logrus"

func LogSentryController()*logrus.Entry{
	return logrus.WithField("Controller", "Sentry")
}
