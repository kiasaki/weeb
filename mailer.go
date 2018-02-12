package weeb

import (
	"errors"
	"fmt"
	"net/smtp"
	"strings"
)

type Mailer interface {
	Send(from, to, subject, body string, headers ...map[string]string) error
}

type MailerConsole struct {
	log *Logger
}

func NewMailerConsole(log *Logger) *MailerConsole {
	return &MailerConsole{log: log}
}

func (m *MailerConsole) Send(from, to, subject, body string, headers ...map[string]string) error {
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

func (m *MailerSMTP) Send(from, to, subject, body string, headers ...map[string]string) error {
	host := m.config.MustGet("mailHost")
	if from == "" {
		return errors.New("Mailer: Missing mail host")
	}
	from = m.config.Get("mailDefaultFrom", from)
	if from == "" {
		return errors.New("Mailer: Missing from")
	}

	headerValues := map[string]string{}
	if len(headers) > 0 {
		headerValues = headers[0]
	}
	headerValues["Form"] = from
	headerValues["To"] = to
	headerValues["Subject"] = subject

	msg := ""
	for k, v := range headerValues {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	formattedBody := strings.Replace(body, "\n", "\r\n", -1)
	msg += fmt.Sprintf("\r\n%s\r\n", formattedBody)

	m.log.Info("sending smtp email", L{"from": from, "to": to, "subject": subject, "len": len(body)})
	auth := smtp.CRAMMD5Auth(m.config.MustGet("mailUsername"), m.config.MustGet("mailPassword"))
	return smtp.SendMail(host, auth, from, []string{to}, []byte(msg))
}
