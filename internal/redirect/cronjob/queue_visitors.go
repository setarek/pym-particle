package cronjob

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"

	"github.com/setarek/pym-particle-microservice/internal/redirect/repository"
	"github.com/setarek/pym-particle-microservice/pkg/logger"
	"github.com/setarek/pym-particle-microservice/pkg/rabbitmq"
)

const (
	VisitedLinksSetName = "visited_links"
)

func QueueVisitors(ctx context.Context, logger logger.Logger, repository *repository.RedirectRedisRepository) func() {
	return func() {
		span, ctx := opentracing.StartSpanFromContext(ctx, "redirectCronJob.QueueVisitors")
		defer span.Finish()

		visitedLinks, err := repository.SMembers(ctx, VisitedLinksSetName)
		if err != nil {
			logger.Error("error while get visitors from redis", err)
		}

		for _, visitedLink := range visitedLinks {
			visitCount, err := repository.GetValue(ctx, fmt.Sprintf("%s:visit", visitedLink))
			if err != nil {
				logger.Error("error while getting visitor count from redis", err)
			}
			_, err = repository.DeleteKey(ctx, fmt.Sprintf("%s:visit", visitedLink))
			if err != nil {
				logger.Error("error while delete visitor count key from redis", err)
			}

			_, err = repository.SRem(ctx, VisitedLinksSetName, visitedLink)
			if err != nil {
				logger.Error("error while remove visitor from visitor set", err)
			}

			message := make(map[string]interface{}, 0)
			message["visited_link"] = visitedLink
			message["visited_count"] = visitCount

			rabbitmq.PublishDurableMessage(ctx, logger, rabbitmq.Visitors, message)
		}
	}
}
