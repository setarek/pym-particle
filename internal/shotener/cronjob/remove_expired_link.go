package cronjob

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/setarek/pym-particle-microservice/config"
	"github.com/setarek/pym-particle-microservice/internal/shotener/repository"
	"github.com/setarek/pym-particle-microservice/pkg/logger"
	"time"
)

func DeleteExpiredLink(ctx context.Context, config *config.Config,logger logger.Logger, repository *repository.ShortenerRepository) func() {
	return func() {
		span, ctx := opentracing.StartSpanFromContext(ctx, "redirectCronJob.QueueVisitors")
		defer span.Finish()

		now := time.Now()
		limit := config.GetInt("expire_link_limit")
		expiryTime := now.AddDate(0, 0, -limit)

		err := repository.DeleteLinks(ctx, expiryTime)
		if err != nil {
			logger.Error("error while deleting expired link", err)
		}
	}
}
