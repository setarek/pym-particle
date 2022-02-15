package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	rsErr "github.com/setarek/pym-particle-microservice/internal/error"
	"github.com/setarek/pym-particle-microservice/internal/shotener/model"
	"github.com/setarek/pym-particle-microservice/pkg/rabbitmq"
	"github.com/setarek/pym-particle-microservice/pkg/utils"
)

type ShortenerRequest struct {
	OriginalLink    string    `json:"original_link" validate:"required"`
}

type ShortenerResponse struct {
	ShortLink    string    `json:"short_link"`
}

type ErrorResponse struct {
	Message    string    `json:"message"`
}

func (h *Handler) Shortener(ctx echo.Context) error {
	var request ShortenerRequest
	if err := ctx.Bind(&request); err != nil {
		h.logger.Error("error while binding body request", err)
		return ctx.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Message: rsErr.EmptyBodyRequest.Error(),
		})
	}

	isUrlValid := utils.CheckUrlValidation(request.OriginalLink)
	if !isUrlValid {
		h.logger.Error("invalid url", request.OriginalLink)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: rsErr.InvalidUrl.Error(),
		})
	}
	randString := utils.GenerateRandString(h.config.GetInt("shortener_length"))

	// todo improve creating short link
	shortLink := fmt.Sprintf("localhost:9005/%s", randString)

	go func() {
		message := make(map[string]interface{}, 0)
		message["original_url"] = request.OriginalLink
		message["shorten"] = randString
		rabbitmq.PublishDurableMessage(ctx.Request().Context(), h.logger, rabbitmq.ShortenerQueue, message)
	}()

	// todo: improve saving in database, goroutine or something else
	linkInfo := model.Link{
		OriginalLink: request.OriginalLink,
		Shorten: randString,
	}

	if err := h.repository.CreateLinkInfo(ctx.Request().Context(), linkInfo); err != nil {
		h.logger.Error("error while creating link info", err)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: rsErr.ServerErr.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, ShortenerResponse{
		ShortLink: shortLink,
	})
}
