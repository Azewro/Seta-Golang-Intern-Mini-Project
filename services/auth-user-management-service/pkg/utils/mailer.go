package utils

import (
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

// SMTPConfig stores SMTP transport configuration.
type SMTPConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	FromEmail string
	FromName  string
}

// LoadSMTPConfig reads SMTP settings from environment variables.
func LoadSMTPConfig() SMTPConfig {
	port, err := strconv.Atoi(strings.TrimSpace(os.Getenv("SMTP_PORT")))
	if err != nil || port <= 0 {
		port = 587
	}

	return SMTPConfig{
		Host:      strings.TrimSpace(os.Getenv("SMTP_HOST")),
		Port:      port,
		Username:  strings.TrimSpace(os.Getenv("SMTP_USERNAME")),
		Password:  os.Getenv("SMTP_PASSWORD"),
		FromEmail: strings.TrimSpace(os.Getenv("SMTP_FROM_EMAIL")),
		FromName:  strings.TrimSpace(os.Getenv("SMTP_FROM_NAME")),
	}
}

// SendVerificationEmail sends account verification email with a confirmation link.
func SendVerificationEmail(cfg SMTPConfig, recipientEmail string, verifyURL string) error {
	from := cfg.FromEmail
	if from == "" {
		from = cfg.Username
	}

	subject := "Verify your account"
	body := fmt.Sprintf("Please verify your account within 5 minutes by opening this link:\n%s\n", verifyURL)

	headers := []string{
		fmt.Sprintf("From: %s <%s>", fallback(cfg.FromName, "Seta App"), from),
		fmt.Sprintf("To: %s", recipientEmail),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
	}
	message := strings.Join(headers, "\r\n") + "\r\n\r\n" + body

	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	return smtp.SendMail(address, auth, from, []string{recipientEmail}, []byte(message))
}

func fallback(value string, defaultValue string) string {
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
}
