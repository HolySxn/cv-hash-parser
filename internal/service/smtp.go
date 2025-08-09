package service

import (
	"fmt"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

type GomailSender struct {
	dialer *gomail.Dialer
}

func NewGomailSender(host, port, username, password string) (*GomailSender, error) {
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP port: %w", err)
	}
	dialer := gomail.NewDialer(host, portInt, username, password)
	return &GomailSender{
		dialer: dialer,
	}, nil
}

func (s *GomailSender) Send(to, subject, body string, attachments []string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.dialer.Username)
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

func (s *GomailSender) GetEmail() string {
	return s.dialer.Username
}
