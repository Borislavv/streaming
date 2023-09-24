package ws

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	host           string // example: "0.0.0.0"
	port           string // example: "9988"
	transportProto string // example: "tcp"

	streamer streamer.Streamer
	logger   logger.Logger
}

func NewWebSocketServer(
	host string,
	port string,
	transportProto string,
	streamer streamer.Streamer,
	logger logger.Logger,
) *Server {
	return &Server{
		host:           host,
		port:           port,
		transportProto: transportProto,
		streamer:       streamer,
		logger:         logger,
	}
}

// Listen is method which running a websocket server
func (s *Server) Listen(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	addr, err := net.ResolveTCPAddr(s.transportProto, net.JoinHostPort(s.host, s.port))
	if err != nil {
		s.logger.Error(err)
		return
	}

	server := &http.Server{
		Addr:    addr.String(),
		Handler: http.HandlerFunc(s.handleConnection),
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer s.logger.Info("stopped")
		if lsErr := server.ListenAndServe(); lsErr != nil && lsErr != http.ErrServerClosed {
			s.logger.Error(lsErr)
			return
		}
	}()

	s.logger.Info("running...")
	<-ctx.Done()
	s.logger.Info("shutting down...")

	serverCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if sdErr := server.Shutdown(serverCtx); sdErr != nil && sdErr != context.Canceled {
		s.logger.Error(sdErr)
		return
	}
}

// handleConnection is method which handle each websocket connection
func (s *Server) handleConnection(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error(err)
		return
	}
	defer func() {
		if err = conn.Close(); err != nil {
			s.logger.Error(err)
			return
		}
	}()

	s.logger.Info(fmt.Sprintf("[%v]: accpted a new connection", conn.RemoteAddr()))

	s.streamer.HandleConn(conn)
}