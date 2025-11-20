package omailer

import (
	"mime/multipart"
	"net/http"
	"omailer/internal/abstraction"
	"omailer/internal/dto"
	"omailer/pkg/util/response"

	"github.com/labstack/echo/v4"
)

type handler struct {
	service Service
}

func NewHandler() *handler {
	return &handler{
		service: NewService(),
	}
}

func (h *handler) OmailerSend(c echo.Context) (err error) {
	payload := new(dto.OmailerSend)
	if err = c.Bind(payload); err != nil {
		return response.ErrorBuilder(http.StatusBadRequest, err, "error bind payload").SendError(c)
	}
	if err = c.Validate(payload); err != nil {
		return response.ErrorBuilder(http.StatusBadRequest, err, "error validate payload").SendError(c)
	}
	if err := c.Request().ParseMultipartForm(64 << 20); err != nil {
		return response.ErrorBuilder(http.StatusBadRequest, err, "error bind multipart/form-data").SendError(c)
	}
	files := []*multipart.FileHeader{}
	for _, fhs := range c.Request().MultipartForm.File {
		files = append(files, fhs...)
	}
	payload.Files = files
	data, err := h.service.OmailerSend(c.(*abstraction.Context), payload)
	if err != nil {
		return response.ErrorResponse(err).SendError(c)
	}
	return response.SuccessResponse(data).SendSuccess(c)
}

func (h *handler) OmailerSendJustMessage(c echo.Context) (err error) {
	payload := new(dto.OmailerSendJustMessage)
	if err = c.Bind(payload); err != nil {
		return response.ErrorBuilder(http.StatusBadRequest, err, "error bind payload").SendError(c)
	}
	if err = c.Validate(payload); err != nil {
		return response.ErrorBuilder(http.StatusBadRequest, err, "error validate payload").SendError(c)
	}
	data, err := h.service.OmailerSendJustMessage(c.(*abstraction.Context), payload)
	if err != nil {
		return response.ErrorResponse(err).SendError(c)
	}
	return response.SuccessResponse(data).SendSuccess(c)
}
