package omailer

import (
	"encoding/json"
	"net/http"
	"net/url"
	"omailer/internal/abstraction"
	"omailer/internal/dto"
	"omailer/pkg/constant"
	"omailer/pkg/general"
	"omailer/pkg/gomail"
	"omailer/pkg/util/response"
	"strconv"
)

type Service interface {
	OmailerSend(ctx *abstraction.Context, payload *dto.OmailerSend) (map[string]interface{}, error)
	OmailerSendJustMessage(ctx *abstraction.Context, payload *dto.OmailerSendJustMessage) (map[string]interface{}, error)
}

type service struct {
	Test string
}

func NewService() Service {
	return &service{Test: constant.APP}
}

type MailConfig struct {
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	AuthEmail    string `json:"auth_email"`
	AuthPassword string `json:"auth_password"`
	SenderName   string `json:"sender_name"`
	Recipient    string `json:"recipient"`
	Subject      string `json:"subject"`
	BodyHTML     string `json:"body_html"`
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

func (s *service) OmailerSendJustMessage(ctx *abstraction.Context, payload *dto.OmailerSendJustMessage) (map[string]interface{}, error) {

	key := "15042003150420031504200315042003"
	dataUrlEnc, err := general.DecryptAES(payload.Data, key)
	if err != nil {
		return nil, response.ErrorBuilder(http.StatusBadRequest, err, "error decode url encode json")
	}

	decoded, err := url.QueryUnescape(dataUrlEnc)
	if err != nil {
		return nil, response.ErrorBuilder(http.StatusBadRequest, err, "error decode url encode json")
	}

	var cfg MailConfig
	if err := json.Unmarshal([]byte(decoded), &cfg); err != nil {
		return nil, response.ErrorBuilder(http.StatusBadRequest, err, "failed to unmarshal json data")
	}

	omailerConfig := gomail.ConfigMailer{
		SmtpHost:     cfg.SMTPHost,
		SmtpPort:     cfg.SMTPPort,
		AuthEmail:    cfg.AuthEmail,
		AuthPassword: cfg.AuthPassword,
		SenderName:   cfg.SenderName,
	}

	err = omailerConfig.SendMail(cfg.Recipient, cfg.Subject, cfg.BodyHTML, nil)
	if err != nil {
		return nil, response.ErrorBuilder(http.StatusBadRequest, err, "failed sent email")
	}

	return map[string]interface{}{
		"message": "success connect & send email!",
	}, nil
}
