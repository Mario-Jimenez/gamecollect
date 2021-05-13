package service

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mario-Jimenez/gamecollect/broker/kafka"
	"github.com/Mario-Jimenez/gamecollect/collect"
	"github.com/Mario-Jimenez/gamecollect/config"
	"github.com/Mario-Jimenez/gamecollect/logger"
	"github.com/Mario-Jimenez/gamecollect/storage"
	"github.com/Mario-Jimenez/gamecollect/subscriber"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
)

// Run service
func Run(serviceName, serviceVersion string) {
	// load app configuration
	conf, err := config.NewFileConfig()
	if err != nil {
		if errors.IsNotFound(err) {
			log.WithFields(log.Fields{
				"error": errors.Details(err),
			}).Error("Configuration file not found")
			return
		}
		if errors.IsNotValid(err) {
			log.WithFields(log.Fields{
				"error": errors.Details(err),
			}).Error("Invalid configuration values")
			return
		}
		log.WithFields(log.Fields{
			"error": errors.Details(err),
		}).Error("Failed to retrieve secrets")
		return
	}

	// initialize logger
	logger.InitializeLogger(serviceName, serviceVersion, conf.Values().LogLevel)

	gamesConsumer := kafka.NewConsumer("games", "collect", conf.Values().KafkaConnection)
	scoreConsumer := kafka.NewConsumer("gamesscore", "collect", conf.Values().KafkaConnection)
	durationConsumer := kafka.NewConsumer("gamesduration", "collect", conf.Values().KafkaConnection)
	pricesConsumer := kafka.NewConsumer("gamesprices", "collect", conf.Values().KafkaConnection)

	memoryStorage := storage.NewInMemoryHandler()

	gameHandler := collect.NewGameHandler(memoryStorage)
	scoreHandler := collect.NewScoreHandler(memoryStorage)
	durationHandler := collect.NewDurationHandler(memoryStorage)
	pricesHandler := collect.NewPricesHandler(memoryStorage)

	gamesSub := subscriber.NewHandler(gamesConsumer, gameHandler.ProcessMessage)
	scoreSub := subscriber.NewHandler(scoreConsumer, scoreHandler.ProcessMessage)
	durationSub := subscriber.NewHandler(durationConsumer, durationHandler.ProcessMessage)
	pricesSub := subscriber.NewHandler(pricesConsumer, pricesHandler.ProcessMessage)

	go gamesSub.InboundMessages()
	time.Sleep(time.Second * 5)
	go scoreSub.InboundMessages()
	go durationSub.InboundMessages()
	go pricesSub.InboundMessages()

	log.Info("Running service...")

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/games", func(c *gin.Context) {
		c.JSON(http.StatusOK, memoryStorage.Get())
	})

	srv := &http.Server{
		Addr:    conf.Values().Port,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("Failed to start server")
		}
	}()

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	log.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Server forced to shutdown")
	}

	if err := gamesConsumer.Close(); err != nil {
		log.WithFields(log.Fields{
			"error": errors.Details(err),
		}).Error("Failed to close consumer")
	}

	if err := scoreConsumer.Close(); err != nil {
		log.WithFields(log.Fields{
			"error": errors.Details(err),
		}).Error("Failed to close consumer")
	}

	if err := durationConsumer.Close(); err != nil {
		log.WithFields(log.Fields{
			"error": errors.Details(err),
		}).Error("Failed to close consumer")
	}

	if err := pricesConsumer.Close(); err != nil {
		log.WithFields(log.Fields{
			"error": errors.Details(err),
		}).Error("Failed to close consumer")
	}
}
