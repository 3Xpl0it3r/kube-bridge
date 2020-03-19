package logging

import "github.com/sirupsen/logrus"

func LogDnsServerController()*logrus.Entry{
	return logrus.WithField("Controller", "dns")
}
