package logging

import "github.com/sirupsen/logrus"

func LogKubeResourceController(resource string)*logrus.Entry{
	return logrus.WithField("Controller", resource)
}
