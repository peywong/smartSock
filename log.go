package smartSock

import (
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func SetLogMode(l *logrus.Logger) {
	logger = l
}
