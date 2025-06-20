package email

import "fmt"

type SimpleEmail struct {
	To      string
	Subject string
	Body    string
}

func (e SimpleEmail) GetTo() string {
	return e.To
}

func (e SimpleEmail) GetMsg() []byte {
	return formatMessage(e.Subject, e.Body)
}

func NewSimpleEmail(to, subject, body string) SimpleEmail {
	return SimpleEmail{
		To:      to,
		Subject: subject,
		Body:    body,
	}
}

type ConfirmationEmail struct {
	To   string
	Link string
}

func (e ConfirmationEmail) GetTo() string {
	return e.To
}

func (e ConfirmationEmail) GetMsg() []byte {
	subject := "Confirm your subscription"
	body := fmt.Sprintf("Click to confirm: %s", e.Link)

	return formatMessage(subject, body)
}

func NewConfirmationEmail(to, link string) ConfirmationEmail {
	return ConfirmationEmail{
		To:   to,
		Link: link,
	}
}

func formatMessage(subject, body string) []byte {
	msg := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)
	return []byte(msg)
}

func GenerateConfirmationLink(baseURL, token string) string {
	return fmt.Sprintf("%s/api/confirm/%s", baseURL, token)
}
