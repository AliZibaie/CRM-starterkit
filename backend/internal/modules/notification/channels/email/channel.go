package email

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"crm-starterkit/config"
	"crm-starterkit/internal/modules/notification/domain"
)

type Channel struct {
	cfg config.EmailConfig
}

func New(cfg config.EmailConfig) *Channel { return &Channel{cfg: cfg} }

func (c *Channel) Name() string { return "email" }

func (c *Channel) Send(ctx context.Context, payload domain.Payload) error {
	_ = ctx
	if c.cfg.SMTPHost == "" || payload.To == "" {
		return nil
	}
	from := c.cfg.FromAddress
	if from == "" {
		from = c.cfg.SMTPUser
	}
	addr := fmt.Sprintf("%s:%d", c.cfg.SMTPHost, c.cfg.SMTPPort)
	msg := []byte("To: " + payload.To + "\r\nSubject: " + payload.Title + "\r\n\r\n" + payload.Body)
	var auth smtp.Auth
	if c.cfg.SMTPUser != "" {
		auth = smtp.PlainAuth("", c.cfg.SMTPUser, c.cfg.SMTPPassword, c.cfg.SMTPHost)
	}
	return smtp.SendMail(addr, auth, from, strings.Split(payload.To, ","), msg)
}
