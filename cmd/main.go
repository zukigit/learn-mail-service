package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
)

type MailConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

func loadConfig() MailConfig {
	return MailConfig{
		Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		Port:     getEnv("SMTP_PORT", "587"),
		User:     getEnv("SMTP_USER", ""),
		Password: getEnv("SMTP_PASS", ""),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func sendMail(cfg MailConfig, to []string, subject, body string) error {
	if cfg.User == "" || cfg.Password == "" {
		return fmt.Errorf("SMTP_USER and SMTP_PASS must be set")
	}

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n%s",
		cfg.User,
		joinRecipients(to),
		subject,
		body,
	)

	auth := smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host)
	addr := cfg.Host + ":" + cfg.Port

	return smtp.SendMail(addr, auth, cfg.User, to, []byte(msg))
}

func joinRecipients(to []string) string {
	return strings.Join(to, ", ")
}

func getEnvList(key string) []string {
	val := os.Getenv(key)
	if val == "" {
		return nil
	}
	parts := strings.Split(val, ",")
	var result []string
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func main() {
	cfg := loadConfig()

	to := getEnvList("MAIL_TO")
	if len(to) == 0 {
		log.Fatal("MAIL_TO must be set (comma-separated emails)")
	}

	subject := "Hello from Go"
	body := "This is a test email sent using Gmail SMTP and Go."

	if err := sendMail(cfg, to, subject, body); err != nil {
		log.Fatalf("failed to send mail: %v", err)
	}

	fmt.Println("email sent successfully")
}
