package gomail

import (
	"errors"
	"io"
	"mime/multipart"
	"omailer/pkg/general"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type ConfigMailer struct {
	SmtpHost     string
	SmtpPort     int
	AuthEmail    string
	AuthPassword string
	SenderName   string
}

func (c *ConfigMailer) SendMail(recipient, subject, bodyHtml string, files []*multipart.FileHeader) error {
	if bodyHtml == "" {
		return errors.New("error parsing body html")
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", c.AuthEmail, c.SenderName)
	mailer.SetHeader("To", recipient)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/plain", general.ParseTemplateEmailToPlainText(bodyHtml))
	mailer.AddAlternative("text/html", bodyHtml)

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		mailer.Attach(fileHeader.Filename,
			gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := io.Copy(w, file)
				return err
			}))
	}

	dialer := gomail.NewDialer(
		c.SmtpHost,
		c.SmtpPort,
		c.AuthEmail,
		c.AuthPassword,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	logrus.Info("Mail sent!")
	return nil
}
