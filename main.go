package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loanengine/internal/handler/rest"
	"github.com/loanengine/internal/repo"
	"github.com/loanengine/internal/service"
	"github.com/loanengine/pkg/middleware"
	"github.com/loanengine/pkg/server"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {
	config := server.NewConfig()
	app := server.NewApp(config)
	g := &errgroup.Group{}
	g.Go(func() error {
		gin.SetMode(gin.DebugMode)
		engine := gin.New()

		endpoint := engine.Group("")

		app.RestHandler = SetupHandler()

		app.SetupRoutesAndMiddleware(endpoint, app.RestHandler)

		engine.HandleMethodNotAllowed = true
		engine.NoRoute(middleware.NoRoute)
		engine.NoMethod(middleware.NoMethod)
		engine.MaxMultipartMemory = 8 << 20 // 8 MiB
		log.Info("starting loan engine server host:", config.Host, " port:", config.Port)

		err := startServer(engine, config.Host, config.Port)
		if err != nil {
			log.Panic("Fatal: service failed to start. ", err.Error())
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Println(err.Error())
		return
	}
}

func SetupHandler() rest.RestHandler {
	repo := repo.NewLoanRepo()
	notificationService := service.NewNotificationService(repo)
	agreementGenerator := service.NewAgreementGenerator(repo)
	loanService := service.NewLoanService(repo, agreementGenerator, notificationService)
	handler := rest.NewRestHandler(loanService)
	return handler
}

func startServer(engine *gin.Engine, host string, port int) error {
	httpServer, err := createHttpServer(engine, host, port)
	if err != nil {
		return err
	}
	err = httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal(err.Error())
		return err
	}
	return nil
}

// Keeping an option to add error in case mTLS required in future
func createHttpServer(engine *gin.Engine, host string, port int) (*http.Server, error) {
	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Handler:      engine,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}, nil
}
