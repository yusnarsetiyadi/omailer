package dto

import "mime/multipart"

type OmailerSend struct {
	// config
	SmtpHost     string `json:"smtp_host" form:"smtp_host" validate:"required"`
	SmtpPort     string `json:"smtp_port" form:"smtp_port" validate:"required"`
	AuthEmail    string `json:"auth_email" form:"auth_email" validate:"required"`
	AuthPassword string `json:"auth_password" form:"auth_password" validate:"required"`
	SenderName   string `json:"sender_name" form:"sender_name" validate:"required"`

	// data
	Recipient string `json:"recipient" form:"recipient" validate:"required"`
	Subject   string `json:"subject" form:"subject" validate:"required"`
	BodyHtml  string `json:"body_html" form:"body_html" validate:"required"`
	Files     []*multipart.FileHeader
}

type OmailerSendJustMessage struct {
	// json encode
	Data string `query:"data"`
}
