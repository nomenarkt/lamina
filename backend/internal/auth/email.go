package auth

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"
)

//go:embed email_templates/confirmation_email.html
var confirmationTemplateFS embed.FS

// SendConfirmationEmail renders and logs the email body with a confirmation link.
func SendConfirmationEmail(toEmail, token string, isResend bool) error {
	baseURL := os.Getenv("FRONTEND_CONFIRM_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000/confirm"
	}
	link := fmt.Sprintf("%s/%s", baseURL, token)

	tmpl, err := template.ParseFS(confirmationTemplateFS, "email_templates/confirmation_email.html")
	if err != nil {
		return fmt.Errorf("❌ failed to parse embedded email template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, struct {
		Token    string
		Link     string
		Year     int
		IsResend bool
	}{
		Token:    token,
		Link:     link,
		Year:     time.Now().Year(),
		IsResend: isResend,
	}); err != nil {
		return fmt.Errorf("❌ failed to render email template: %w", err)
	}

	log.Printf("📧 Confirmation email to %s:\n%s", toEmail, body.String())
	return nil
}
