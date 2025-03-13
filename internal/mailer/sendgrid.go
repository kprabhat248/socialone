package mailer

import (
	"bytes"
	"fmt"
	"log"
	"text/template"

	"github.com/smtp2go-oss/smtp2go-go"
)

// SMTP2GoMailer represents a mailer using the SMTP2Go API.
type SMTP2GoMailer struct {
	fromEmail string
	apiKey    string
}

// NewSMTP2GoMailer creates a new SMTP2GoMailer instance.
func NewSMTP2GoMailer(apiKey, fromEmail string) *SMTP2GoMailer {
	return &SMTP2GoMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
	}
}

// Send sends an email using the SMTP2Go API with templated content.
func (m *SMTP2GoMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {
	// Read the template file
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Prepare subject and body from the template
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return fmt.Errorf("failed to execute subject template: %w", err)
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return fmt.Errorf("failed to execute body template: %w", err)
	}

	// Create an instance of the SMTP2Go email struct
	emailData := smtp2go.Email{
		From:     m.fromEmail,
		To:       []string{email},
		Subject:  subject.String(),
		TextBody: body.String(),
		HtmlBody: body.String(), // You can use a different body for HTML if needed
	}

	// Send the email using the SMTP2Go API key (API Key should be set in environment variables)
	res, err := smtp2go.Send(&emailData)
	if err != nil {
		log.Printf("An error occurred: %s", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent successfully to %v", email)
	log.Printf("Response: %v", res)
	return nil
}
