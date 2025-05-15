package auth

import (
	"fmt"
	"log"
	"os"
)

// SendConfirmationEmail logs or sends a confirmation email with a token link.
func SendConfirmationEmail(toEmail, token string) error {
	baseURL := os.Getenv("BACKEND_URL") // e.g., http://localhost:8080
	link := fmt.Sprintf("%s/api/v1/auth/confirm/%s", baseURL, token)

	// Just log it for now — wire to SMTP later
	log.Printf("📧 Confirmation email to %s: %s", toEmail, link)

	// Optional real SMTP (setup later)
	// smtp.SendMail(...)

	return nil
}
