package email

import (
	"fmt"

	"gopkg.in/mail.v2"
)

// SendDatabaseEmail sends an email with the database backup file attached
func SendDatabaseEmail(backupFilePath string) error {
	// Define email details
	emailSender := "gauravyadav00729@gmail.com"
	emailPass := "aiif rfjh jwdd iptp"
	recipient := "yadav.gaurav00729@gmail.com"
	subject := "Database Backup"
	body := "Please find the attached database backup file."

	// Create a new dialer
	d := mail.NewDialer("smtp.gmail.com", 587, emailSender, emailPass)
	d.TLSConfig = nil // Disable TLS verification

	// Create a new email message
	m := mail.NewMessage()
	m.SetHeader("From", emailSender)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Attach the backup file
	m.Attach(backupFilePath)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	fmt.Printf("Email sent successfully with attachment: %s\n", backupFilePath)
	return nil
}
