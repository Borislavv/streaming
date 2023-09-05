package video

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/service"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/gorilla/mux"
	"net/http"
)

const CreatePath = "/video"

type CreateVideoController struct {
	logger   service.Logger
	builder  builder.Video
	service  service.Video
	response response.Responder
}

func NewCreateController(
	logger service.Logger,
	builder builder.Video,
	service service.Video,
	response response.Responder,
) *CreateVideoController {
	return &CreateVideoController{
		logger:   logger,
		builder:  builder,
		service:  service,
		response: response,
	}
}

func (c *CreateVideoController) Create(w http.ResponseWriter, r *http.Request) {
	videoDto, err := c.builder.BuildCreateRequestDtoFromRequest(r)
	if err != nil {
		c.logger.Error(err)
		c.response.RespondError(w, err)
		return
	}

	id, err := c.service.Create(videoDto)
	if err != nil {
		c.logger.Error(err)
		c.response.RespondError(w, err)
		return
	}

	c.response.RespondData(w, id)
	w.WriteHeader(http.StatusCreated)
	return
}

func (c *CreateVideoController) AddRoute(router *mux.Router) {
	router.
		Path(CreatePath).
		HandlerFunc(c.Create).
		Methods(http.MethodPost)
}