package dto

import "mime/multipart"

type OmailerSend struct {
	SmtpHost     string                `json:"smtp_host" form:"smtp_host" validate:"required" example:"smtp.gmail.com"`
	SmtpPort     string                `json:"smtp_port" form:"smtp_port" validate:"required" example:"587"`
	AuthEmail    string                `json:"auth_email" form:"auth_email" validate:"required" example:"user@gmail.com"`
	AuthPassword string                `json:"auth_password" form:"auth_password" validate:"required" example:"app-password"`
	SenderName   string                `json:"sender_name" form:"sender_name" validate:"required" example:"John Doe"`
	Recipient    string                `json:"recipient" form:"recipient" validate:"required" example:"recipient@example.com"`
	Subject      string                `json:"subject" form:"subject" validate:"required" example:"Hello World"`
	BodyHtml     string                `json:"body_html" form:"body_html" validate:"required" example:"<h1>Hello</h1><p>Email body here</p>"`
	Files        []*multipart.FileHeader
}

type OmailerSendJustMessage struct {
	Data string `query:"data" example:"%7B%22smtp_host%22...%7D"`
}
