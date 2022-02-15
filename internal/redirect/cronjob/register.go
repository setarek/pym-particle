package cronjob

import (
	"context"

	cron "github.com/robfig/cron/v3"

	"github.com/setarek/pym-particle-microservice/internal/redirect/repository"
	"github.com/setarek/pym-particle-microservice/pkg/logger"
)

var cronClient *cron.Cron

func init() {
	cronClient = cron.New()
}

func InitializeCronJobs(ctx context.Context, logger logger.Logger, repository *repository.RedirectRedisRepository) {
	cronClient.AddFunc("@every 60s", QueueVisitors(ctx, logger, repository))
	cronClient.Start()
}

func StopCronJobs() {
	cronClient.Stop()
}

