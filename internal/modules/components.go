package modules

import (
	"github.com/tamboto2000/dealls-dating-svc/internal/config"
	"github.com/tamboto2000/dealls-dating-svc/pkg/mail"
)

type Components struct {
	mailer *mail.Mailer
}

func NewComponents(cfg config.Config) (*Components, error) {
	smtpCfg := cfg.SMTP
	mailer := mail.NewMailer(
		smtpCfg.SenderName,
		smtpCfg.Email,
		smtpCfg.Password,
		smtpCfg.Host,
		smtpCfg.Port,
	)

	comps := Components{
		mailer: mailer,
	}

	return &comps, nil
}

func (c *Components) Mailer() *mail.Mailer {
	return c.mailer
}
