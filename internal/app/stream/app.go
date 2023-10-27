package stream

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/mongodb"
	"github.com/Borislavv/video-streaming/internal/infrastructure/server/ws"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/detector"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/logger/stdout"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/reader"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/handler"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/handler/strategy"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/listener"
	ws2 "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/proto/ws"
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

type StreamingApp struct {
	cfg config
}

func NewStreamingApp() *StreamingApp {
	return &StreamingApp{cfg: config{}}
}

// Run is method which running the streaming part of app
func (app *StreamingApp) Run(mWg *sync.WaitGroup) {
	defer mWg.Done()
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	// init. logger and close func.
	loggerService, cls := stdout.NewLogger(ctx, 10, 10)
	defer func() {
		cancel()
		wg.Wait()
		cls()
	}()

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

	// video repository
	videoRepository := mongodb.NewVideoRepository(db, loggerService, time.Minute)

	// resource reader service
	readerService := reader.NewFileReaderService(ctx, loggerService, app.cfg.ChunkSize)

	// custom websocket communication protocol
	wsCommunicator := ws2.NewWebSocketCommunicator(loggerService)

	// resource codecs determiner
	codecsDetector := detector.NewResourceCodecInfo(ctx, loggerService)

	// websocket actions listener
	actionsListener := listener.NewWebSocketActionsListener(loggerService, wsCommunicator)

	// websocket actions handler
	actionsHandler := handler.NewWebSocketActionsHandler(
		ctx,
		loggerService,
		[]strategy.ActionStrategy{
			strategy.NewStreamByIDActionStrategy(
				ctx, loggerService, videoRepository, readerService, codecsDetector, wsCommunicator,
			),
		},
	)

	// resource streaming service
	streamingService := streamer.NewStreamingService(loggerService, actionsListener, actionsHandler)

	wg.Add(1)
	go ws.NewWebSocketServer( // websocket server
		app.cfg.Host, app.cfg.Port, app.cfg.Transport, streamingService, loggerService,
	).Listen(ctx, wg)

	<-app.shutdown()
}

func (app *StreamingApp) shutdown() chan os.Signal {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	return stopCh
}
