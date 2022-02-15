package queue

import (
	"context"
	"encoding/json"

	"github.com/opentracing/opentracing-go"

	"github.com/setarek/pym-particle-microservice/internal/shotener/repository"
	"github.com/setarek/pym-particle-microservice/pkg/logger"
	"github.com/setarek/pym-particle-microservice/pkg/rabbitmq"
	"github.com/setarek/pym-particle-microservice/pkg/utils"
)

func ListenToVisitedLink(ctx context.Context, logger logger.Logger, repository *repository.ShortenerRepository) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redirectQueue.ListenToLinkShortener")
	defer span.Finish()

	rabbitMQClient := rabbitmq.GetQueue()
	ch, _ := rabbitMQClient.Consume(rabbitmq.Visitors,
		"",
		false,false,false,false,nil)

	for delivery := range ch {
		var  messageMap map[string]interface{}
		err := json.Unmarshal(delivery.Body, &messageMap)
		if err != nil {
			logger.Error("error while unmarshalling consumed message", err)
		}
		err = repository.UpdateVisitInfo(ctx, utils.ParseString(messageMap["visited_link"]), utils.ParseInt64(messageMap["visited_count"]))
		if err != nil {
			logger.Error("error while update visited info ", err)
		}

		delivery.Ack(false)
	}
}
