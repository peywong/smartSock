package smartSock

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func SetLogMode(l *zap.Logger) {
	logger = l
}
