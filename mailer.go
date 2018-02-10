package weeb

import (
	"errors"
	"fmt"
	"net/smtp"
	"strings"
)

type Mailer interface {
	Send(from, to, subject, body string) error
}

type MailerConsole struct {
	log *Logger
}

func NewMailerConsole(log *Logger) *MailerConsole {
	return &MailerConsole{log: log}
}

func (m *MailerConsole) Send(from, to, subject, body string) error {
	m.log.Info("sending email", L{"from": from, "to": to, "subject": subject, "body": body})
	return nil
}

type MailerSMTP struct {
	config *Config
	log    *Logger
}

func NewMailerSMTP(log *Logger, config *Config) *MailerSMTP {
	return &MailerSMTP{log: log, config: config}
}

func (m *MailerSMTP) Send(from, to, subject, body string) error {
	from = m.config.Get("mailDefaultFrom", from)
	if from == "" {
		return errors.New("Mailer: Missing from")
	}
	host := m.config.MustGet("mailHost")
	auth := smtp.CRAMMD5Auth(m.config.MustGet("mailUsername"), m.config.MustGet("mailPassword"))
	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", to, subject, strings.Replace(body, "\n", "\r\n", -1))

	m.log.Info("sending smtp email", L{"from": from, "to": to, "subject": subject, "len": len(body)})

	return smtp.SendMail(host, auth, from, []string{to}, []byte(msg))
}
