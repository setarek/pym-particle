package handler

import (
	"github.com/setarek/pym-particle-microservice/internal/redirect/repository"
	"github.com/setarek/pym-particle-microservice/pkg/logger"
)

type Handler struct {
	logger          logger.Logger
	redisRepository *repository.RedirectRedisRepository
}

func NewHandler(logger logger.Logger, redisRepository *repository.RedirectRedisRepository) *Handler {
	return &Handler{
		logger: logger,
		redisRepository: redisRepository,
	}
}