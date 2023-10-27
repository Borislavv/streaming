package errors

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/logger"
	"net/http"
)

const (
	repositoryType = "application"

	publicRepositoryLevel  = logger.ErrorLevel
	publicRepositoryStatus = http.StatusNotFound

	internalRepositoryLevel  = logger.CriticalLevel
	internalRepositoryStatus = http.StatusInternalServerError
)

type EntityNotFoundError struct{ publicError }

func NewEntityNotFoundError(entity string, by string) *EntityNotFoundError {
	return &EntityNotFoundError{
		publicError{
			errored{
				ErrorMessage: fmt.Sprintf("'%v' not found by given '%s'", entity, by),
				ErrorType:    repositoryType,
				errorLevel:   publicRepositoryLevel,
				errorStatus:  publicRepositoryStatus,
			},
		},
	}
}

func IsEntityNotFoundError(err error) bool {
	_, ok := err.(*EntityNotFoundError)
	return ok
}

type InternalRepositoryError struct{ internalError }

func NewInternalRepositoryError(msg string) *InternalRepositoryError {
	return &InternalRepositoryError{
		internalError{
			errored{
				ErrorMessage: msg,
				ErrorType:    repositoryType,
				errorLevel:   internalRepositoryLevel,
				errorStatus:  internalRepositoryStatus,
			},
		},
	}
}
