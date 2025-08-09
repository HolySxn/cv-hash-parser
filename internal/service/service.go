package service

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	randomChars   = "abcdefghijklmnopqrstuvwxyz0123456789"
	sourceCodeZip = "source_code.zip"
)

type EmailSender interface {
	Send(to, subject, body string, attacments []string) error
	GetEmail() string
}

type Service struct {
	logger    *slog.Logger
	smtp      EmailSender
	recipient string

	wg sync.WaitGroup
}

func NewService(logger *slog.Logger, smtp EmailSender, recipient string) *Service {
	return &Service{
		logger:    logger,
		smtp:      smtp,
		recipient: recipient,
	}
}

func (s *Service) Wait() {
	s.wg.Wait()
}

func (s *Service) ProcessCV(url string) error {
	hash := s.calculateSHA256(url)
	userID := s.generateUserID(hash)

	report := &Report{
		CvURL:     url,
		Hash:      hash,
		UserID:    userID,
		Email:     s.smtp.GetEmail(),
		Timestamp: time.Now(),
	}

	reportFileName, err := s.saveReport(report)
	if err != nil {
		return fmt.Errorf("failed to save report: %w", err)
	}

	htmlBody, err := s.parseAndPrepareHTMLBody(report)
	if err != nil {
		return fmt.Errorf("failed to parse and prepare HTML body: %w", err)
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		err = s.smtp.Send(
			s.recipient,
			fmt.Sprintf("Golang Test – %s", userID),
			htmlBody,
			[]string{reportFileName, sourceCodeZip},
		)
		if err != nil {
			s.logger.Error("failed to send email", "error", err)
		} else {
			s.logger.Info("email sent successfully in background", "user_id", userID)
		}
	}()

	return nil
}

func (s *Service) parseAndPrepareHTMLBody(report *Report) (string, error) {
	tmpl, err := template.ParseFiles("templates/email_template.html")
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, report); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return body.String(), nil
}

func (s *Service) saveReport(report *Report) (string, error) {
	fileName := fmt.Sprintf("report_%s.json", report.UserID)
	file, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create report file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		return "", fmt.Errorf("failed to encode report: %w", err)
	}

	return fileName, nil
}

func (s *Service) calculateSHA256(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (s *Service) generateUserID(hash string) string {
	b := make([]byte, 4)
	for i := range b {
		b[i] = randomChars[rand.Intn(len(randomChars))]
	}
	return fmt.Sprintf("%s-%s", hash[:8], string(b))
}
