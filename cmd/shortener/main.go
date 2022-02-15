package main

import (
	"context"
	"fmt"
	"github.com/setarek/pym-particle-microservice/internal/shotener/cronjob"
	"github.com/setarek/pym-particle-microservice/internal/shotener/queue"

	"github.com/opentracing/opentracing-go"

	"github.com/setarek/pym-particle-microservice/config"
	"github.com/setarek/pym-particle-microservice/internal/shotener/handler"
	"github.com/setarek/pym-particle-microservice/internal/shotener/repository"
	"github.com/setarek/pym-particle-microservice/internal/shotener/router"
	"github.com/setarek/pym-particle-microservice/pkg/jaeger"
	"github.com/setarek/pym-particle-microservice/pkg/logger"
	"github.com/setarek/pym-particle-microservice/pkg/mongo"
	"github.com/setarek/pym-particle-microservice/pkg/rabbitmq"
)

func main()  {

	ctx := context.Background()

	config, err := config.InitConfig()
	if err != nil {
		panic(fmt.Errorf("error while initializing config: %v+", err))
	}

	logger := logger.NewLogger(config)
	logger.InitLogger()

	db, err := mongo.InitDB(config)
	if err != nil {
		panic(fmt.Errorf("error while initializing mongo: %v+", err))
	}

	repository := repository.NewShortenerRepository(db)

	rabbitmq.InitRabbitMQ(config, logger)

	tracer, closer, err := jaeger.InitJaeger(config)
	if err != nil {
		logger.Fatal("error while initializing Jaeger", err)
	}

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	router := router.New()
	v1 := router.Group("/api/v1")
	handler := handler.NewHandler(config, logger, repository)
	handler.Register(v1)

	go queue.ListenToVisitedLink(ctx, logger, repository)
	go cronjob.InitializeCronJobs(ctx, config, logger, repository)
	defer cronjob.StopCronJobs()

	router.Start(fmt.Sprintf("%s:%s", config.GetString("hostname"), config.GetString("port")))
}