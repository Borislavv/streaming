package errors

import (
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/logger"
	"net/http"
)

const (
	applicationType               = "application"
	internalServerErrorMessage    = "internal server error"
	internalServerErrorStatusCode = http.StatusInternalServerError
	internalServerErrorLevel      = logger.ErrorLevel
)

type InternalServerError struct{ publicError }

func NewInternalServerError(msg string) InternalServerError {
	if msg == "" {
		msg = internalServerErrorMessage
	}
	return InternalServerError{
		publicError{
			errored{
				ErrorMessage: msg,
				ErrorType:    applicationType,
				errorStatus:  internalServerErrorStatusCode,
				errorLevel:   internalServerErrorLevel,
			},
		},
	}
}
