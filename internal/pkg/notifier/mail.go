package sender

import (
	"encoding/json"
	"net/smtp"

	"github.com/avraam311/improved-calendar-service/internal/models"
)

type Mail struct {
	host     string
	port     string
	auth     smtp.Auth
	from     string
	password string
}

func NewMail(host, port, user, from, password string) *Mail {
	auth := smtp.PlainAuth("", user, password, host)
	return &Mail{
		host:     host,
		port:     port,
		auth:     auth,
		from:     from,
		password: password,
	}
}

func (m *Mail) SendMessage(msg []byte) error {
	ev := &models.EventCreate{}
	err := json.Unmarshal(msg, ev)
	if err != nil {
		return err
	}
	to := []string{ev.Mail}
	msgToSend := []byte("Subject: Notifying about event\r\n" +
		"\r\n" +
		"You have event planned in an hour: " + ev.Event + "\r\n")

	err = smtp.SendMail(m.host+":"+m.port, m.auth, m.from, to, msgToSend)
	if err != nil {
		return err
	}
	return nil
}
