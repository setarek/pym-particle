package handler

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"

	rsErr "github.com/setarek/pym-particle-microservice/internal/error"
)

type ErrorResponse struct {
	Message    string    `json:"message"`
}

const (
	VisitedLinksSetName = "visited_links"
)

func (h *Handler) Redirect(ctx echo.Context) error {
	shorten := ctx.Param("shorten")
	if shorten == "" {
		h.logger.Error("empty query param")
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: rsErr.ErrorNoQueryParam.Error(),
		})

	}

	originalLink, err := h.redisRepository.GetValue(ctx.Request().Context(), shorten)
	if err != nil {
		h.logger.Error("error while getting original link from redis", err)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: rsErr.ErrorExpiredLink.Error(),
		})
	}

	if originalLink == "" {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: rsErr.ServerErr.Error(),
		})
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		_, err = h.redisRepository.IncrValue(fmt.Sprintf("%s:visit", shorten))
		if err != nil {
			h.logger.Error("error while increment visit count", err)
		}
		_, err  = h.redisRepository.SAdd(VisitedLinksSetName, shorten)
		if err != nil {
			h.logger.Error("error while add visited link to set", err)
		}
		wg.Done()
	}()
	wg.Wait()

	return ctx.Redirect(http.StatusPermanentRedirect, originalLink)
}
