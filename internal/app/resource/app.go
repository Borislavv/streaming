package resource

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	domainauth "github.com/Borislavv/video-streaming/internal/domain/service/authenticator"
	domainresource "github.com/Borislavv/video-streaming/internal/domain/service/resource"
	domainuploader "github.com/Borislavv/video-streaming/internal/domain/service/uploader"
	domainuser "github.com/Borislavv/video-streaming/internal/domain/service/user"
	domainvideo "github.com/Borislavv/video-streaming/internal/domain/service/video"
	"github.com/Borislavv/video-streaming/internal/domain/validator"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/render"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/rest/audio"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/rest/auth"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/rest/resource"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/rest/user"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/rest/video"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/static"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/request"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/mongodb"
	"github.com/Borislavv/video-streaming/internal/infrastructure/server/http"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/logger/stdout"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/storager"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/tokenizer"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/uploader"
	_ "github.com/Borislavv/video-streaming/internal/infrastructure/service/uploader"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/uploader/file"
	"github.com/caarlos0/env/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type ResourcesApp struct {
	cfg config
}

func NewResourcesApp() *ResourcesApp {
	return &ResourcesApp{cfg: config{}}
}

// Run is method which running the REST API part of app
func (app *ResourcesApp) Run(mWg *sync.WaitGroup) {
	defer mWg.Done()
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	// init. loggerService and close func.
	loggerService, cls := stdout.NewLogger(ctx, app.cfg.LoggerErrorsBufferCap, app.cfg.LoggerRequestsBufferCap)
	defer func() {
		cancel()
		wg.Wait()
		cls()
	}()

	// parse env. config
	if err := env.Parse(&app.cfg); err != nil {
		loggerService.Critical(err)
		return
	}

	// init. mongodb client
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(app.cfg.MongoUri))
	if err != nil {
		loggerService.Critical(err)
		return
	}
	defer func() { _ = mongoClient.Disconnect(ctx) }()

	// ping mongodb
	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		loggerService.Critical(err)
		return
	}

	// connect to target mongodb database
	db := mongoClient.Database(app.cfg.MongoDb)

	// request param. resolver
	reqParamsExtractor := request.NewParametersExtractor()

	// response service
	responseService := response.NewResponseService(loggerService)

	// video repository
	videoRepository := mongodb.NewVideoRepository(db, loggerService, time.Minute)

	// resource repository
	resourceRepository := mongodb.NewResourceRepository(db, loggerService, time.Minute)

	// resource validator
	resourceValidator := validator.NewResourceValidator(ctx, resourceRepository, app.cfg.MaxFilesizeThreshold)

	// video validator
	videoValidator := validator.NewVideoValidator(
		ctx, loggerService, resourceValidator, videoRepository, resourceRepository,
	)

	// video builder
	videoBuilder := builder.NewVideoBuilder(
		ctx, loggerService, reqParamsExtractor, videoRepository, resourceRepository,
	)

	// filesystem storage
	filesystemStorage := storager.NewFilesystemStorage(ctx, loggerService)

	// filename computer
	filenameComputerService := file.NewNameService()

	// resource builder
	resourceBuilder := builder.NewResourceBuilder(
		loggerService, app.cfg.ResourceFormFilename, app.cfg.InMemoryFileSizeThreshold,
	)

	// user repository
	userRepository := mongodb.NewUserRepository(db, loggerService, time.Minute)

	// user builder
	userBuilder := builder.NewUserBuilder(ctx, loggerService, reqParamsExtractor, userRepository)

	// user validator
	userValidator := validator.NewUserValidator(ctx, loggerService, userRepository, app.cfg.AdminContactEmail)

	var uploaderStrategy domainuploader.Uploader
	if app.cfg.UploadingStrategy == uploader.MultipartFormUploadingType {
		// used parsing of full form into RAM
		uploaderStrategy =
			uploader.NewNativeUploader(
				loggerService,
				filesystemStorage,
				filenameComputerService,
				app.cfg.ResourceFormFilename,
				app.cfg.InMemoryFileSizeThreshold,
			)
	} else if app.cfg.UploadingStrategy == uploader.MultipartPartUploadingType {
		// used partial reading from multipart.Part
		uploaderStrategy =
			uploader.NewPartsUploader(
				loggerService,
				filesystemStorage,
				filenameComputerService,
			)
	}

	// auth. validator
	authValidator := validator.NewAuthValidator(loggerService, app.cfg.AdminContactEmail)

	// auth. builder
	authBuilder := builder.NewAuthBuilder(loggerService)

	/**
	 * CRUD services.
	 */
	resourceService := domainresource.NewResourceService(
		ctx, loggerService, uploaderStrategy, resourceValidator, resourceBuilder, resourceRepository, filesystemStorage,
	)
	videoService := domainvideo.NewCRUDService(
		ctx, loggerService, videoBuilder, videoValidator, videoRepository, resourceService,
	)
	userService := domainuser.NewCRUDService(
		ctx, loggerService, userBuilder, userValidator, userRepository, videoService,
	)
	tokenService := tokenizer.NewJwtService(
		loggerService, app.cfg.JwtSecretSalt,
	)
	authService := domainauth.NewAuthService(
		loggerService, userService, authValidator, tokenService,
	)

	wg.Add(1)
	go http.NewHttpServer(
		ctx,
		app.cfg.Host,
		app.cfg.Port,
		app.cfg.Transport,
		app.cfg.ApiVersionPrefix,
		app.cfg.RenderVersionPrefix,
		app.cfg.StaticVersionPrefix,
		app.InitRestApiControllers(
			loggerService,
			responseService,
			resourceBuilder,
			resourceService,
			videoBuilder,
			videoService,
			userBuilder,
			userService,
			authBuilder,
			authService,
		),
		app.InitNativeRenderingControllers(
			loggerService,
			responseService,
		),
		app.InitStaticServingControllers(
			loggerService,
			responseService,
		),
		loggerService,
		reqParamsExtractor,
	).Listen(ctx, wg)

	<-app.shutdown()
}

func (app *ResourcesApp) shutdown() chan os.Signal {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	return stopCh
}

func (app *ResourcesApp) InitRestApiControllers(
	loggerService *stdout.Logger,
	responseService response.Responder,
	// resource deps.
	resourceBuilder builder.Resource,
	resourceService domainresource.CRUD,
	// video deps.
	videoBuilder builder.Video,
	videoService domainvideo.CRUD,
	// user. deps.
	userBuilder builder.User,
	userService domainuser.CRUD,
	// auth. deps.
	authBuilder builder.Auth,
	authService domainauth.Authenticator,
) []controller.Controller {
	return []controller.Controller{
		// resource
		resource.NewUploadController(
			loggerService,
			resourceBuilder,
			resourceService,
			responseService,
		),
		// video
		video.NewCreateController(
			loggerService,
			videoBuilder,
			videoService,
			responseService,
		),
		video.NewDeleteController(
			loggerService,
			videoBuilder,
			videoService,
			responseService,
		),
		video.NewGetController(
			loggerService,
			videoBuilder,
			videoService,
			responseService,
		),
		video.NewListController(
			loggerService,
			videoBuilder,
			videoService,
			responseService,
		),
		video.NewUpdateController(
			loggerService,
			videoBuilder,
			videoService,
			responseService,
		),
		// audio
		audio.NewCreateController(),
		audio.NewDeleteController(),
		audio.NewGetController(),
		audio.NewListController(),
		audio.NewUpdateController(),
		// user
		user.NewCreateController(
			loggerService,
			userBuilder,
			userService,
			responseService,
		),
		user.NewUpdateUserController(
			loggerService,
			userBuilder,
			userService,
			responseService,
		),
		user.NewGetController(
			loggerService,
			userBuilder,
			userService,
			responseService,
		),
		user.NewDeleteController(
			loggerService,
			userBuilder,
			userService,
			responseService,
		),
		// auth
		auth.NewAuthorizationController(
			loggerService,
			authBuilder,
			authService,
			responseService,
		),
		auth.NewRegistrationController(
			loggerService,
			userBuilder,
			userService,
			responseService,
		),
	}
}

func (app *ResourcesApp) InitNativeRenderingControllers(
	loggerService *stdout.Logger,
	responseService response.Responder,
) []controller.Controller {
	return []controller.Controller{
		render.NewIndexController(loggerService, responseService),
	}
}

func (app *ResourcesApp) InitStaticServingControllers(
	loggerService *stdout.Logger,
	responseService response.Responder,
) []controller.Controller {
	return []controller.Controller{
		static.NewFilesController(loggerService, responseService),
	}
}
