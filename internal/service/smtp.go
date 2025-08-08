package service

import (
	"time"

	"gopkg.in/gomail.v2"
)

type GomailSender struct {
	dialer *gomail.Dialer
	from   string
}

func NewGomailSender(host string, port int, username, password, from string) *GomailSender {
	dialer := gomail.NewDialer(host, port, username, password)
	return &GomailSender{
		dialer: dialer,
		from:   from,
	}
}

func (s *GomailSender) Send(to, subject, body string, attachments []string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	for _, attachment := range attachments {
		m.Attach(attachment)
	}

	for i := 0; i < 3; i++ {
		if err := s.dialer.DialAndSend(m); err == nil {
			return nil
		}
		time.Sleep(time.Second * 5)
	}

	return nil
}
