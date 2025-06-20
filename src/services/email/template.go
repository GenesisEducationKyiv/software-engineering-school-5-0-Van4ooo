package email

import "fmt"

type SimpleMail struct {
	To      string
	Subject string
	Body    string
}

func (e SimpleMail) GetTo() string {
	return e.To
}

func (e SimpleMail) GetMsg() []byte {
	return formatMessage(e.Subject, e.Body)
}

func NewSimpleMail(to, subject, body string) SimpleMail {
	return SimpleMail{
		To:      to,
		Subject: subject,
		Body:    body,
	}
}

type ConfirmationMail struct {
	To   string
	Link string
}

func (e ConfirmationMail) GetTo() string {
	return e.To
}

func (e ConfirmationMail) GetMsg() []byte {
	subject := "Confirm your subscription"
	body := fmt.Sprintf("Click to confirm: %s", e.Link)

	return formatMessage(subject, body)
}

func NewConfirmationMail(to, link string) ConfirmationMail {
	return ConfirmationMail{
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
