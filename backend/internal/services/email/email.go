package email

import (
	"fmt"
	"log"
	"time"

	"github.com/FoxyHunter7/geocachingbrughia-backend/internal/config"
	"github.com/FoxyHunter7/geocachingbrughia-backend/internal/database"
	"gopkg.in/gomail.v2"
)

type Service struct {
	cfg    config.SMTPConfig
	dialer *gomail.Dialer
}

func New(cfg config.SMTPConfig) *Service {
	var dialer *gomail.Dialer
	if cfg.Host != "" {
		dialer = gomail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Pass)
	}

	return &Service{
		cfg:    cfg,
		dialer: dialer,
	}
}

// SendNewContactNotification sends an email when a new contact form is submitted
func (s *Service) SendNewContactNotification(fromEmail, subject, message string, submissionID int64) {
	if s.dialer == nil {
		log.Println("Email not configured, skipping notification")
		return
	}

	// Truncate message for preview
	messagePreview := message
	if len(messagePreview) > 500 {
		messagePreview = messagePreview[:500] + "..."
	}

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
    <p style="background: #fff3cd; border: 1px solid #ffc107; padding: 15px; color: #856404;">
        Automatische melding, stuur uw antwoordt naar de persoon zelf, niet deze mail.
    </p>

    <h2>Nieuw Contactformulier Ontvangen</h2>
    
    <table style="width: 100%%; margin: 20px 0;">
        <tr><td><strong>Van:</strong></td><td><a href="mailto:%s">%s</a></td></tr>
        <tr><td><strong>Onderwerp:</strong></td><td>%s</td></tr>
        <tr><td><strong>Ontvangen:</strong></td><td>%s</td></tr>
    </table>

    <h3>Bericht:</h3>
    <p style="background: #f8f9fa; padding: 15px; border: 1px solid #ddd;">%s</p>

    <hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
    <p style="font-size: 12px; color: #666;">
        Geocaching Brughia contact systeem<br>
        Submission ID: #%d
    </p>
</body>
</html>
`, fromEmail, fromEmail, subject, time.Now().Format("02-01-2006 15:04"), messagePreview, submissionID)

	plainBody := fmt.Sprintf(`!! Automatische melding, stuur uw antwoordt naar de persoon zelf, niet deze mail. !!

NIEUW CONTACTFORMULIER

Van: %s
Onderwerp: %s
Ontvangen: %s

Bericht:
%s

---
Submission ID: #%d
`, fromEmail, subject, time.Now().Format("02-01-2006 15:04"), messagePreview, submissionID)

	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.From)
	m.SetHeader("To", s.cfg.NotificationEmail)
	m.SetHeader("Subject", fmt.Sprintf("[Nieuw Contact] %s", subject))
	m.SetBody("text/plain", plainBody)
	m.AddAlternative("text/html", htmlBody)

	if err := s.dialer.DialAndSend(m); err != nil {
		log.Printf("Failed to send contact notification email: %v", err)
	} else {
		log.Printf("Contact notification email sent for submission #%d", submissionID)
	}
}

// SendReminderEmail sends a reminder for unanswered contact submissions
func (s *Service) SendReminderEmail(submissions []PendingSubmission) {
	if s.dialer == nil || len(submissions) == 0 {
		return
	}

	// Build list of pending items
	var itemsHTML, itemsPlain string
	for _, sub := range submissions {
		daysAgo := int(time.Since(sub.CreatedAt).Hours() / 24)
		itemsHTML += fmt.Sprintf(`
            <tr>
                <td style="padding: 8px; border-bottom: 1px solid #ddd;">#%d</td>
                <td style="padding: 8px; border-bottom: 1px solid #ddd;">%s</td>
                <td style="padding: 8px; border-bottom: 1px solid #ddd;">%s</td>
                <td style="padding: 8px; border-bottom: 1px solid #ddd;">%d dagen geleden</td>
            </tr>
        `, sub.ID, sub.Email, sub.Subject, daysAgo)

		itemsPlain += fmt.Sprintf("- #%d | %s | %s | %d dagen geleden\n", sub.ID, sub.Email, sub.Subject, daysAgo)
	}

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
    <p style="background: #f8d7da; border: 1px solid #f5c6cb; padding: 15px; color: #721c24;">
        <strong>Herinnering: %d onbeantwoorde berichten</strong>
    </p>

    <p>De volgende contactformulieren wachten al meer dan 3 dagen op een antwoord:</p>

    <table style="width: 100%%; border-collapse: collapse; margin: 20px 0;">
        <tr style="background: #f8f9fa;">
            <th style="padding: 10px; text-align: left; border-bottom: 2px solid #ddd;">ID</th>
            <th style="padding: 10px; text-align: left; border-bottom: 2px solid #ddd;">Van</th>
            <th style="padding: 10px; text-align: left; border-bottom: 2px solid #ddd;">Onderwerp</th>
            <th style="padding: 10px; text-align: left; border-bottom: 2px solid #ddd;">Wacht sinds</th>
        </tr>
        %s
    </table>
</body>
</html>
`, len(submissions), itemsHTML)

	plainBody := fmt.Sprintf(`HERINNERING: %d ONBEANTWOORDE BERICHTEN

%s`, len(submissions), itemsPlain)

	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.From)
	m.SetHeader("To", s.cfg.NotificationEmail)
	m.SetHeader("Subject", fmt.Sprintf("[Herinnering] %d onbeantwoorde berichten", len(submissions)))
	m.SetBody("text/plain", plainBody)
	m.AddAlternative("text/html", htmlBody)

	if err := s.dialer.DialAndSend(m); err != nil {
		log.Printf("Failed to send reminder email: %v", err)
	} else {
		log.Printf("Reminder email sent for %d pending submissions", len(submissions))
	}
}

type PendingSubmission struct {
	ID        int64
	Email     string
	Subject   string
	CreatedAt time.Time
}

// StartReminderScheduler starts a background job to send reminders
func (s *Service) StartReminderScheduler(db *database.DB, reminderDays int) {
	if s.dialer == nil {
		log.Println("Email not configured, reminder scheduler disabled")
		return
	}

	ticker := time.NewTicker(1 * time.Hour)
	log.Printf("Reminder scheduler started (checking every hour, reminder after %d days)", reminderDays)

	for range ticker.C {
		s.checkAndSendReminders(db, reminderDays)
	}
}

func (s *Service) checkAndSendReminders(db *database.DB, reminderDays int) {
	// Find submissions that:
	// 1. Have status 'new'
	// 2. Were created more than reminderDays ago
	// 3. Haven't received a reminder in the last 24 hours (or never)
	cutoffDate := time.Now().AddDate(0, 0, -reminderDays)
	reminderCutoff := time.Now().Add(-24 * time.Hour)

	rows, err := db.Query(`
		SELECT id, email, subject, created_at
		FROM contact_submissions
		WHERE status = 'new'
		  AND created_at < ?
		  AND (last_reminder_sent_at IS NULL OR last_reminder_sent_at < ?)
	`, cutoffDate, reminderCutoff)

	if err != nil {
		log.Printf("Error checking for pending submissions: %v", err)
		return
	}
	defer rows.Close()

	var pending []PendingSubmission
	var ids []int64

	for rows.Next() {
		var p PendingSubmission
		if err := rows.Scan(&p.ID, &p.Email, &p.Subject, &p.CreatedAt); err != nil {
			continue
		}
		pending = append(pending, p)
		ids = append(ids, p.ID)
	}

	if len(pending) == 0 {
		return
	}

	// Send reminder
	s.SendReminderEmail(pending)

	// Update last_reminder_sent_at for all these submissions
	for _, id := range ids {
		db.Exec("UPDATE contact_submissions SET last_reminder_sent_at = CURRENT_TIMESTAMP WHERE id = ?", id)
	}
}

// SendAdminInvitation sends an email to a newly created admin user with their credentials
func (s *Service) SendAdminInvitation(toEmail, name, tempPassword string) {
	if s.dialer == nil {
		log.Printf("Email not configured, skipping admin invitation for %s", toEmail)
		log.Printf("============================================")
		log.Printf("NEW ADMIN USER CREDENTIALS (email not sent):")
		log.Printf("  Email: %s", toEmail)
		log.Printf("  Password: %s", tempPassword)
		log.Printf("============================================")
		return
	}

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
    <h1 style="color: #2c5530;">Geocaching Brughia</h1>

    <p>Beste %s,</p>
    
    <p>U bent uitgenodigd als beheerder. Hieronder uw inloggegevens:</p>

    <table style="background: #f8f9fa; padding: 15px; margin: 20px 0; border-left: 4px solid #28a745;">
        <tr><td><strong>Email:</strong></td><td>%s</td></tr>
        <tr><td><strong>Wachtwoord:</strong></td><td>%s</td></tr>
    </table>

    <p style="background: #fff3cd; border: 1px solid #ffc107; padding: 15px; color: #856404;">
        <strong>Belangrijk:</strong> Bij uw eerste login wordt u gevraagd om uw wachtwoord te wijzigen.
    </p>

    <hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
    <p style="font-size: 12px; color: #666;">
        Automatisch verzonden door Geocaching Brughia.<br>
        Niet verwacht? Negeer deze email.
    </p>
</body>
</html>
`, name, toEmail, tempPassword)

	plainBody := fmt.Sprintf(`Beste %s,

U bent uitgenodigd als beheerder voor Geocaching Brughia.

Inloggegevens:
Email: %s
Wachtwoord: %s

BELANGRIJK: Wijzig uw wachtwoord bij eerste login.
`, name, toEmail, tempPassword)

	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.From)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Uitnodiging: Geocaching Brughia Admin Panel")
	m.SetBody("text/plain", plainBody)
	m.AddAlternative("text/html", htmlBody)

	if err := s.dialer.DialAndSend(m); err != nil {
		log.Printf("Failed to send admin invitation email to %s: %v", toEmail, err)
		log.Printf("============================================")
		log.Printf("ADMIN INVITATION EMAIL FAILED - CREDENTIALS:")
		log.Printf("  Email: %s", toEmail)
		log.Printf("  Password: %s", tempPassword)
		log.Printf("============================================")
	} else {
		log.Printf("Admin invitation email sent to %s", toEmail)
	}
}
