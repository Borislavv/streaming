package stdout

import (
	"context"
	abstractLogger "github.com/Borislavv/video-streaming/internal/infrastructure/service/logger"
	"os"
)

type Logger struct {
	*abstractLogger.Logger
}

func NewLogger(ctx context.Context, errBuff int, reqBuff int) (logger *Logger, closeFunc func()) {
	l, closeFunc := logger.NewLogger(ctx, os.Stdout, errBuff, reqBuff)
	return &Logger{Logger: l}, closeFunc
}
