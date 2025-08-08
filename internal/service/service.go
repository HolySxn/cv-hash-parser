package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"
)

const (
	randomChars = "abcdefghijklmnopqrstuvwxyz0123456789"
)

type EmailSender interface {
	Send(to, subject, body string, attacments []string) error
}

type Service struct {
	// logger *slog.Logger
	smtp EmailSender
}

func NewService(logger *slog.Logger, smtp EmailSender) *Service {
	return &Service{
		// logger: logger,
		smtp: smtp,
	}
}

func (s *Service) ProcessCV(url string) error {
	hash := s.calculateSHA256(url)
	userID := s.generateUseID(hash)

	report := &Report{
		CvURL:     url,
		Hash:      hash,
		UserID:    userID,
		Email:     "", // This should be set based on your application logic
		Timestamp: time.Now(),
	}

	fileName, err := s.saveReport(report)
	if err != nil {
		return fmt.Errorf("failed to save report: %w", err)
	}

	return nil
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

func (s *Service) generateUseID(hash string) string {
	b := make([]byte, 4)
	for i := range b {
		b[i] = randomChars[rand.Intn(len(randomChars))]
	}
	return fmt.Sprintf("%s-%s", hash[:8], string(b))
}
