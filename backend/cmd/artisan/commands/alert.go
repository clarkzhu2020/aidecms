package commands

import (
	"fmt"
	"net/smtp"
	"os"
	"path/filepath"
	"time"
)

type EmailConfig struct {
	Server    string
	Port      int
	Username  string
	From      string
	Password  string
	To        []string
	Templates map[string]string
}

func SendAlertEmail(subject, body string) error {
	config := loadEmailConfig()
	if config.Server == "" {
		return fmt.Errorf("email not configured")
	}

	// 记录发送日志
	logEntry := fmt.Sprintf("[%s] Sending email: %s",
		time.Now().Format("2006-01-02 15:04:05"), subject)
	logEmailActivity(logEntry)

	auth := smtp.PlainAuth("", config.Username, config.Password, config.Server)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		config.From, config.To[0], subject, body)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", config.Server, config.Port),
		auth,
		config.From,
		config.To,
		[]byte(msg),
	)

	// 记录发送结果
	if err != nil {
		logEmailActivity(fmt.Sprintf("[ERROR] %v", err))
	} else {
		logEmailActivity("Email sent successfully")
	}

	return err
}

func logEmailActivity(message string) {
	filePath := filepath.Join("storage", "logs", "email.log")
	os.MkdirAll(filepath.Dir(filePath), 0755)

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	logEntry := fmt.Sprintf("[%s] %s\n",
		time.Now().Format("2006-01-02 15:04:05"), message)
	f.WriteString(logEntry)
}

func loadEmailConfig() EmailConfig {
	return EmailConfig{
		Server:   os.Getenv("SMTP_SERVER"),
		Port:     587,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
		To:       []string{os.Getenv("SMTP_TO")},
	}
}

func SetupEmailAlert(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: alert:setup <config_file>")
		return
	}

	config := loadEmailConfig()
	config.Templates = map[string]string{
		"error":   "templates/error.tmpl",
		"warning": "templates/warning.tmpl",
	}

	if err := validateConfig(config); err != nil {
		fmt.Printf("Invalid configuration: %v\n", err)
		return
	}

	fmt.Printf("Email alert configured for %v\n", config.To)
}

func SendTestEmail(args []string) {
	config := loadEmailConfig()
	if err := validateConfig(config); err != nil {
		fmt.Printf("Cannot send test email: %v\n", err)
		return
	}

	subject := "Artisan Alert Test"
	body := "This is a test email from Artisan CLI"
	if err := SendAlertEmail(subject, body); err != nil {
		fmt.Printf("Failed to send test email: %v\n", err)
	} else {
		fmt.Println("Test email sent successfully")
	}
}

func validateConfig(config EmailConfig) error {
	if config.Server == "" {
		return fmt.Errorf("SMTP server not configured")
	}
	if config.Username == "" || config.Password == "" {
		return fmt.Errorf("SMTP credentials not configured")
	}
	if len(config.To) == 0 {
		return fmt.Errorf("no recipients configured")
	}
	return nil
}
