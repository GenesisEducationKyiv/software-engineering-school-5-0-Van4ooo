package services

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendConfirmationEmail(email, baseURL, token string) error {
	auth := authAccountInSmtp()
	err := smtp.SendMail(
		os.Getenv("SMTP_ADDR"),
		auth,
		os.Getenv("SMTP_NAME"),
		[]string{email},
		[]byte(generateMessage(baseURL, token)),
	)
	return err
}

func generateLink(baseURL, token string) string {
	return fmt.Sprintf("%s/api/confirm/%s", baseURL, token)
}

func generateMessage(baseURL, token string) string {
	return fmt.Sprintf(
		"Subject: Confirm your subscription\r\n\r\n"+
			"Click to confirm your subscription: %s", generateLink(baseURL, token),
	)
}

func authAccountInSmtp() smtp.Auth {
	return smtp.PlainAuth(
		"",
		os.Getenv("SMTP_NAME"),
		os.Getenv("SMTP_PASS"),
		os.Getenv("SMTP_HOST"),
	)
}

func SendEmail(to, subject, body string) error {
	auth := authAccountInSmtp()
	msg := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)
	err := smtp.SendMail(
		os.Getenv("SMTP_ADDR"),
		auth,
		os.Getenv("SMTP_NAME"),
		[]string{to},
		[]byte(msg),
	)
	return err
}
