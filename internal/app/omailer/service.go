package omailer

import (
	"net/http"
	"omailer/internal/abstraction"
	"omailer/internal/dto"
	"omailer/pkg/constant"
	"omailer/pkg/gomail"
	"omailer/pkg/util/response"
	"strconv"
)

type Service interface {
	OmailerSend(ctx *abstraction.Context, payload *dto.OmailerSend) (map[string]interface{}, error)
}

type service struct {
	Test string
}

func NewService() Service {
	return &service{Test: constant.APP}
}

func (s *service) OmailerSend(ctx *abstraction.Context, payload *dto.OmailerSend) (map[string]interface{}, error) {
	smtpPort, err := strconv.Atoi(payload.SmtpPort)
	if err != nil {
		return nil, response.ErrorBuilder(http.StatusBadRequest, err, "error convert smtp port")
	}
	omailerConfig := gomail.ConfigMailer{
		SmtpHost:     payload.SmtpHost,
		SmtpPort:     smtpPort,
		AuthEmail:    payload.AuthEmail,
		AuthPassword: payload.AuthPassword,
		SenderName:   payload.SenderName,
	}

	err = omailerConfig.SendMail(payload.Recipient, payload.Subject, payload.BodyHtml, payload.Files)
	if err != nil {
		return nil, response.ErrorBuilder(http.StatusInternalServerError, err, "server_error")
	}

	return map[string]interface{}{
		"message": "success connect & send email!",
	}, nil
}
