package main

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"

	"github.com/setarek/pym-particle-microservice/config"
	"github.com/setarek/pym-particle-microservice/internal/redirect/cronjob"
	"github.com/setarek/pym-particle-microservice/internal/redirect/handler"
	"github.com/setarek/pym-particle-microservice/internal/redirect/queue"
	"github.com/setarek/pym-particle-microservice/internal/redirect/repository"
	"github.com/setarek/pym-particle-microservice/internal/redirect/router"
	"github.com/setarek/pym-particle-microservice/pkg/jaeger"
	"github.com/setarek/pym-particle-microservice/pkg/logger"
	"github.com/setarek/pym-particle-microservice/pkg/rabbitmq"
	"github.com/setarek/pym-particle-microservice/pkg/redis"
)

func main()  {

	config, err := config.InitConfig()
	if err != nil {
		panic(fmt.Errorf("error while initializing config: %v+", err))
	}

	logger := logger.NewLogger(config)
	logger.InitLogger()

	redisClient := redis.GetRedisClient(config)
	redisRepository := repository.NewRedirectRepository(redisClient)

	ctx := context.Background()

	rabbitmq.InitRabbitMQ(config, logger)
	go queue.ListenToLinkShortener(ctx, logger, redisRepository)

	tracer, closer, err := jaeger.InitJaeger(config)
	if err != nil {
		logger.Fatal("error while initializing Jaeger", err)
	}
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	router := router.New()
	v1 := router.Group("")
	handler := handler.NewHandler(logger, redisRepository)
	handler.Register(v1)

	go cronjob.InitializeCronJobs(ctx, logger, redisRepository)
	defer cronjob.StopCronJobs()

	router.Start(fmt.Sprintf("%s:%s", config.GetString("hostname"), config.GetString("port")))
}