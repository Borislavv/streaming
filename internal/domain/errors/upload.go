package errors

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/logger"
	"net/http"
)

const (
	uploadErrType         = "application"
	publicUploadErrLevel  = logger.WarningLevel
	publicUploadErrStatus = http.StatusBadRequest
)

type ResourceAlreadyExistsError struct{ publicError }

func NewResourceAlreadyExistsError(name string) *ResourceAlreadyExistsError {
	return &ResourceAlreadyExistsError{
		publicError{
			errored{
				ErrorMessage: fmt.Sprintf("the resource '%v' being loaded already exists", name),
				ErrorType:    uploadErrType,
				errorStatus:  publicUploadErrStatus,
				errorLevel:   publicUploadErrLevel,
			},
		},
	}
}
