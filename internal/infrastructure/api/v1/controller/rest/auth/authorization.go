package auth

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder/interface"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	authenticator_interface "github.com/Borislavv/video-streaming/internal/domain/service/authenticator/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	response_interface "github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response/interface"
	"github.com/gorilla/mux"
	"net/http"
)

const AuthorizationPath = "/authorization"

type AuthorizationController struct {
	logger        logger_interface.Logger
	builder       builder_interface.Auth
	authenticator authenticator_interface.Authenticator
	responder     response_interface.Responder
}

func NewAuthorizationController(serviceContainer di_interface.ContainerManager) (*AuthorizationController, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	authBuilder, err := serviceContainer.GetAuthBuilder()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	authService, err := serviceContainer.GetAuthService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	responseService, err := serviceContainer.GetResponderService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &AuthorizationController{
		logger:        loggerService,
		builder:       authBuilder,
		authenticator: authService,
		responder:     responseService,
	}, nil
}

func (c *AuthorizationController) GetAccessToken(w http.ResponseWriter, r *http.Request) {
	// building an auth. request DTO
	req, err := c.builder.BuildAuthRequestDTOFromRequest(r)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	// getting access token
	token, err := c.authenticator.Auth(req)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	c.responder.Respond(w, token)
}

func (c *AuthorizationController) AddRoute(router *mux.Router) {
	router.
		Path(AuthorizationPath).
		HandlerFunc(c.GetAccessToken).
		Methods(http.MethodPost)
}
