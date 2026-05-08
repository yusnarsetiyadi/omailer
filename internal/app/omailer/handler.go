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

// OmailerSend
// @Summary Send email with attachments
// @Description Send email using SMTP with optional file attachments (multipart/form-data)
// @Tags Omailer
// @Accept mpfd
// @Produce json
// @Param smtp_host formData string true "SMTP Host"
// @Param smtp_port formData string true "SMTP Port"
// @Param auth_email formData string true "Auth Email"
// @Param auth_password formData string true "Auth Password"
// @Param sender_name formData string true "Sender Name"
// @Param recipient formData string true "Recipient Email"
// @Param subject formData string true "Email Subject"
// @Param body_html formData string true "Email Body (HTML)"
// @Param files formData file false "Attachment Files"
// @Success 200 {object} response.MetaSuccess
// @Failure 400 {object} response.MetaError
// @Router /send [post]
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

// OmailerSendJustMessage
// @Summary Send email from URL-encoded JSON
// @Description Send email by passing URL-encoded JSON config in the `data` query parameter
// @Tags Omailer
// @Accept json
// @Produce json
// @Param data query string true "URL-encoded JSON containing smtp_host, smtp_port, auth_email, auth_password, sender_name, recipient, subject, body_html"
// @Success 200 {object} response.MetaSuccess
// @Failure 400 {object} response.MetaError
// @Router /send/just-message [get]
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
