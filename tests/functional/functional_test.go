// +build functional

package functional

import (
	"context"
	"github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/trustwallet/blockatlas/api"
	"github.com/trustwallet/blockatlas/pkg/ginutils"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/trustwallet/blockatlas/config"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/blockatlas/platform"
)

func TestApis(t *testing.T) {
	_ = os.Setenv("ATLAS_GIN_MODE", "debug")
	config.LoadConfig(os.Getenv("TEST_CONFIG"))

	logger.InitLogger()
	platform.Init(viper.GetString("platform"))

	sg := sentrygin.New(sentrygin.Options{})
	p := ":8420"

	engine := gin.New()

	engine.Use(gin.Recovery())

	engine.Use(ginutils.CheckReverseProxy, sg)
	engine.Use(ginutils.CORSMiddleware())

	engine.OPTIONS("/*path", ginutils.CORSMiddleware())

	engine.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api.SetupPlatformAPI(engine)

	signalForExit := make(chan os.Signal, 1)

	signal.Notify(signalForExit,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	srv := &http.Server{
		Addr:    ":8420",
		Handler: engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Info("Application failed", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	defer func() {
		if err := srv.Shutdown(ctx); err != nil {
			logger.Info("Server Shutdown: ", err)
		}
	}()

	time.Sleep(time.Second * 2)

	var wg sync.WaitGroup
	cl := newClient(t, p)
	for _, r := range engine.Routes() {
		wg.Add(1)
		t.Run(r.Path, func(t *testing.T) {
			go cl.doTests(r.Method, r.Path, &wg)
		})
	}
	wg.Wait()
}
