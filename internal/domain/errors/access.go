package errors

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/logger"
	"net/http"
)

const (
	accessErrType         = "access"
	publicAccessErrLevel  = logger.InfoLevel
	publicAccessErrStatus = http.StatusBadRequest
)

type AccessDeniedError struct{ publicError }

func NewAccessDeniedError(msg string) *AccessDeniedError {
	baseMsg := "access denied"
	if msg != "" {
		msg = fmt.Sprintf("%v: %v", baseMsg, msg)
	} else {
		msg = baseMsg
	}

	return &AccessDeniedError{
		publicError{
			errored{
				ErrorMessage: msg,
				ErrorType:    accessErrType,
				errorLevel:   publicAccessErrLevel,
				errorStatus:  publicAccessErrStatus,
			},
		},
	}
}
