package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/setarek/pym-particle-microservice/internal/redirect/repository"
	"github.com/setarek/pym-particle-microservice/pkg/logger"
	"github.com/setarek/pym-particle-microservice/pkg/rabbitmq"
	"github.com/setarek/pym-particle-microservice/pkg/utils"
)

const ShortLinkExpiryTTL = 60

func ListenToLinkShortener(ctx context.Context, logger logger.Logger, repository *repository.RedirectRedisRepository) {
	logger.Info("listen to new shorted link")
	span, ctx := opentracing.StartSpanFromContext(ctx, "redirectQueue.ListenToLinkShortener")
	defer span.Finish()

	rabbitMQClient := rabbitmq.GetQueue()
	ch, _ := rabbitMQClient.Consume(rabbitmq.ShortenerQueue,
		"",
		false,false,false,false,nil)

	for delivery := range ch {
		var messageMap map[string]interface{}
		err := json.Unmarshal(delivery.Body, &messageMap)
		if err != nil {
			logger.Error("error while unmarshalling consumed message", err)
		}
		_, err = repository.SetValue(ctx, utils.ParseString(messageMap["shorten"]), utils.ParseString(messageMap["original_url"]),  ShortLinkExpiryTTL * time.Second)
		if err != nil {
			logger.Error("error while setting short link in redis", err)
		}

		delivery.Ack(false)
	}

}
