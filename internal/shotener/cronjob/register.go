package cronjob

import (
	"context"
	"github.com/robfig/cron/v3"
	"github.com/setarek/pym-particle-microservice/config"
	"github.com/setarek/pym-particle-microservice/internal/shotener/repository"
	"github.com/setarek/pym-particle-microservice/pkg/logger"
)

var cronClient *cron.Cron

func init() {
	cronClient = cron.New()
}
func InitializeCronJobs(ctx context.Context, config *config.Config, logger logger.Logger, repository *repository.ShortenerRepository) {

	//cronClient.AddFunc("@every 60s", DeleteExpiredLink(ctx, logger, repository))
	cronClient.AddFunc("* 05 * * *", DeleteExpiredLink(ctx, config, logger, repository))
	cronClient.Start()
}

func StopCronJobs() {
	cronClient.Stop()
}
