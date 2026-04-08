package mail

import (
	"fmt"
	"net/smtp"
	"os"
)

type Mailer interface {
	SendOTP(to, code string) error
	SendForgotPasswordOTP(to, code string) error
	SendStaffResetRequest(ownerEmail, staffEmail string) error
}

type smtpMailer struct {
	host       string
	port       string
	username   string
	password   string
	from       string
	senderName string
}

func NewSMTPMailer() Mailer {
	return &smtpMailer{
		host:       os.Getenv("SMTP_HOST"),
		port:       os.Getenv("SMTP_PORT"),
		username:   os.Getenv("SMTP_USER"),
		password:   os.Getenv("SMTP_PASS"),
		from:       os.Getenv("SMTP_SENDER_EMAIL"),
		senderName: os.Getenv("SMTP_SENDER_NAME"),
	}
}

func (m *smtpMailer) SendOTP(to, code string) error {
	subject := "Subject: Kode Verifikasi Registrasi - POS App\r\n"
	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: auto; padding: 20px; border: 1px solid #eee; border-radius: 10px;">
			<h2 style="color: #4F46E5; text-align: center;">Verifikasi Email Anda</h2>
			<p>Halo,</p>
			<p>Terima kasih telah mendaftar di POS App. Gunakan kode OTP di bawah ini untuk memverifikasi akun Anda:</p>
			<div style="background: #F3F4F6; padding: 20px; text-align: center; font-size: 32px; font-weight: bold; letter-spacing: 5px; color: #111827; border-radius: 8px; margin: 20px 0;">
				%s
			</div>
			<p style="color: #6B7280; font-size: 14px;">Kode ini hanya berlaku selama 5 menit. Jangan bagikan kode ini kepada siapapun.</p>
			<hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;">
			<p style="font-size: 12px; color: #9CA3AF; text-align: center;">&copy; 2024 POS App. All rights reserved.</p>
		</div>
	`, code)

	return m.send(to, subject, body)
}

func (m *smtpMailer) SendForgotPasswordOTP(to, code string) error {
	subject := "Subject: Reset Password - POS App\r\n"
	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: auto; padding: 20px; border: 1px solid #eee; border-radius: 10px;">
			<h2 style="color: #4F46E5; text-align: center;">Lupa Password?</h2>
			<p>Halo,</p>
			<p>Kami menerima permintaan untuk mereset password akun Anda. Gunakan kode OTP di bawah ini untuk melanjutkan:</p>
			<div style="background: #F3F4F6; padding: 20px; text-align: center; font-size: 32px; font-weight: bold; letter-spacing: 5px; color: #111827; border-radius: 8px; margin: 20px 0;">
				%s
			</div>
			<p style="color: #6B7280; font-size: 14px;">Kode ini hanya berlaku selama 5 menit. Jika Anda tidak merasa melakukan permintaan ini, abaikan email ini.</p>
			<hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;">
			<p style="font-size: 12px; color: #9CA3AF; text-align: center;">&copy; 2024 POS App. All rights reserved.</p>
		</div>
	`, code)

	return m.send(to, subject, body)
}

func (m *smtpMailer) SendStaffResetRequest(ownerEmail, staffEmail string) error {
	subject := "Subject: Permintaan Reset Password Staff - POS App\r\n"
	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: auto; padding: 20px; border: 1px solid #eee; border-radius: 10px;">
			<h2 style="color: #4F46E5; text-align: center;">Permintaan Reset Password</h2>
			<p>Halo Owner,</p>
			<p>Staff dengan email <b>%s</b> telah mengirimkan permintaan untuk melakukan reset password.</p>
			<p>Silakan masuk ke aplikasi untuk membantu Staff Anda mengganti password mereka melalui menu Manajemen Staff.</p>
			<hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;">
			<p style="font-size: 12px; color: #9CA3AF; text-align: center;">&copy; 2024 POS App. All rights reserved.</p>
		</div>
	`, staffEmail)

	return m.send(ownerEmail, subject, body)
}

func (m *smtpMailer) send(to, subject, body string) error {
	if m.host == "" {
		fmt.Printf("DEBUG: MOCK EMAIL to %s\nSubject: %s\nBody: %s\n", to, subject, body)
		return nil
	}

	fromHeader := fmt.Sprintf("From: %s <%s>\r\n", m.senderName, m.from)
	toHeader := fmt.Sprintf("To: %s\r\n", to)
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

	msg := []byte(fromHeader + toHeader + subject + mime + body)

	auth := smtp.PlainAuth("", m.username, m.password, m.host)
	addr := fmt.Sprintf("%s:%s", m.host, m.port)

	return smtp.SendMail(addr, auth, m.from, []string{to}, msg)
}
