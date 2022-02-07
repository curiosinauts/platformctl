package mailutil

import (
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Send sends email using SendGrid
func Send(fromName string, fromEmail string, subject string, toName string, toEmail string, plainTextContent string, htmlContent string, apiKey string) (*rest.Response, error) {
	from := mail.NewEmail(fromName, fromEmail)
	to := mail.NewEmail(toName, toEmail)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(apiKey)
	return client.Send(message)
}
