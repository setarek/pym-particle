package handler

import (
	"github.com/setarek/pym-particle-microservice/config"
	"github.com/setarek/pym-particle-microservice/internal/shotener/repository"
	"github.com/setarek/pym-particle-microservice/pkg/logger"
)


type Handler struct {
	config *config.Config
	logger logger.Logger
	repository *repository.ShortenerRepository

}

func NewHandler(config *config.Config, logger logger.Logger, repository *repository.ShortenerRepository) *Handler {
	return &Handler{
		repository: repository,
		logger: logger,
		config: config,
	}
}