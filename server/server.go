package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

)

func New(router http.Handler) *http.Server {
	restConfig := config.NewRestConfig()
	addr := fmt.Sprintf("%s:%s", restConfig.Host, restConfig.Port)

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}

func StartWithGracefulShutdown(s *http.Server) {
	logger.Log.Info("starting server")

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info("shutdown Server")

	timeout := config.NewAppConfig().GracefulShutdownTimeout

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		logger.Log.Info("Server Shutdown:", err)
	}

	logger.Log.Info("server exiting")
}
