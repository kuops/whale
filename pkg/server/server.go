package server

import (
	"context"
	"errors"
	"fmt"
	limit "github.com/aviddiviner/gin-limit"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"time"
	"whale/pkg/client/cri"
	"whale/pkg/collector"
	"whale/pkg/config"
)

func setupRouter(config *config.Config) *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	pprof.Register(r)
	r.Use(limit.MaxAllowed(config.MaxRequests))
	r.GET("/", index())
	r.GET("/metrics",metrics())
	return r
}

func Run(ctx context.Context,config *config.Config) error {
	runtimeClient,runtimeConn,err  := cri.NewClient(config.SocketPath)
	defer func(runtimeConn *grpc.ClientConn) {
		_ = runtimeConn.Close()
	}(runtimeConn)

	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx,"containerdClient",runtimeClient)
	ctx = context.WithValue(ctx,"namespace",config.NameSpace)
	ctx = context.WithValue(ctx,"mountPaths",config.MountPaths)
	ctx = context.WithValue(ctx,"allNamespaces",config.AllNamespaces)
	ctx = context.WithValue(ctx, "nodeIP",config.NodeIP)
	collector.Register(ctx)

	router := setupRouter(config)
	srv := &http.Server{
		Addr:    config.ListenAddress,
		Handler: router,
	}

	errCh := make(chan error, 1)
	go func() {
		<-ctx.Done()

		log.Println("server: context closed")
		shutdownCtx, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()

		log.Println("server: shutting down")
		errCh <- srv.Shutdown(shutdownCtx)
	}()

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err,http.ErrServerClosed) {
		return fmt.Errorf("failed to serve: %w", err)
	}

	if err := <-errCh;err != nil {
		return fmt.Errorf("failed to shutdown: %w", err)
	}

	return nil
}
