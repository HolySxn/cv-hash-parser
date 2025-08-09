package service

import (
	"log/slog"
	"os"
	"strings"
	"testing"
)

type mockSMTP struct {
}

func (m *mockSMTP) Send(to, subject, body string, attachments []string) error {
	return nil
}

func (m *mockSMTP) GetEmail() string {
	return "test@gmail.com"
}

func TestService_ProcessCV(t *testing.T) {
	if err := os.MkdirAll("templates", 0755); err != nil {
		t.Fatalf("failed to create dummy templates directory: %v", err)
	}
	templateFile := "templates/email_template.html"
	templateContent := `
	<!DOCTYPE html>
	<html>
	<body>
	  <h2>Отчёт о проверке резюме (User ID: {{.UserID}})</h2>
	  <p>Здравствуйте,</p>
	  <p>Это автоматический отчёт, сгенерированный по вашему резюме, которое вы отправили со страницы <a href="{{.CvURL}}">{{.CvURL}}</a>.</p>
	  <p>В приложении вы найдёте JSON-отчёт и исходный код проекта.</p>
	  <p>С уважением,<br>Система автоматической проверки</p>
	</body>
	</html>
	`
	if err := os.WriteFile(templateFile, []byte(templateContent), 0644); err != nil {
		t.Fatalf("failed to create dummy template file: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	smtp := &mockSMTP{}

	s := NewService(logger, smtp, "test")

	cvURL := "https://example.com/resume"

	err := s.ProcessCV(cvURL)
	if err != nil {
		t.Fatalf("ProcessCV() error = %v", err)
	}

	files, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("failed to read directory: %v", err)
	}

	var reportFileName string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "report_") {
			reportFileName = file.Name()
			break
		}
	}

	if reportFileName == "" {
		t.Fatal("report file was not created")
	}

	if _, err := os.Stat(reportFileName); os.IsNotExist(err) {
		t.Errorf("report file %s was not created", reportFileName)
	}

	os.Remove(reportFileName)
	os.RemoveAll("templates")
}
