package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

const (
	gmailSMTPHost = "smtp.gmail.com"
	gmailSMTPPort = "587"
)

type Sender struct {
	from        string
	appPassword string
	timeout     time.Duration
}

func NewSender(from, appPassword string) *Sender {
	return &Sender{
		from:        strings.TrimSpace(from),
		appPassword: strings.TrimSpace(appPassword),
		timeout:     10 * time.Second,
	}
}

func (s *Sender) SendRegistrationThankYou(ctx context.Context, recipient, name string) error {
	if s.from == "" || s.appPassword == "" {
		return fmt.Errorf("Gmail SMTP configuration is incomplete; set EMAIL_FROM and EMAIL_APP_PASSWORD")
	}

	fromAddress, err := parseEmailAddress(s.from)
	if err != nil {
		return fmt.Errorf("invalid EMAIL_FROM: %w", err)
	}
	recipientAddress, err := parseEmailAddress(recipient)
	if err != nil {
		return fmt.Errorf("invalid recipient email: %w", err)
	}

	address := net.JoinHostPort(gmailSMTPHost, gmailSMTPPort)
	dialer := net.Dialer{Timeout: s.timeout}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("connect to Gmail SMTP: %w", err)
	}
	defer conn.Close()

	if deadline, ok := ctx.Deadline(); ok {
		if err := conn.SetDeadline(deadline); err != nil {
			return fmt.Errorf("set SMTP deadline: %w", err)
		}
	} else if err := conn.SetDeadline(time.Now().Add(s.timeout)); err != nil {
		return fmt.Errorf("set SMTP deadline: %w", err)
	}

	client, err := smtp.NewClient(conn, gmailSMTPHost)
	if err != nil {
		return fmt.Errorf("start Gmail SMTP session: %w", err)
	}
	defer client.Close()

	if ok, _ := client.Extension("STARTTLS"); !ok {
		return fmt.Errorf("Gmail SMTP server does not support STARTTLS")
	}
	if err := client.StartTLS(&tls.Config{
		ServerName: gmailSMTPHost,
		MinVersion: tls.VersionTLS12,
	}); err != nil {
		return fmt.Errorf("start TLS with Gmail SMTP: %w", err)
	}

	auth := smtp.PlainAuth("", fromAddress, s.appPassword, gmailSMTPHost)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("authenticate with Gmail SMTP: %w", err)
	}
	if err := client.Mail(fromAddress); err != nil {
		return fmt.Errorf("set Gmail sender: %w", err)
	}
	if err := client.Rcpt(recipientAddress); err != nil {
		return fmt.Errorf("set registration email recipient: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("start registration email body: %w", err)
	}
	message := fmt.Sprintf(
		"From: Matthew E-Wallet <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: Selamat datang di GoWallet!\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/plain; charset=UTF-8\r\n"+
			"\r\n"+
			"Halo %s! 👋\r\n"+
			"Selamat datang di GoWallet.\r\n"+
			"Registrasi akun kamu berhasil. Sekarang kamu bisa mulai menggunakan fitur dompet digital untuk melakukan top up saldo simulasi dan transfer antar pengguna.\r\n"+
			"Terima kasih telah bergabung dengan GoWallet! 💚\r\n",
		fromAddress,
		recipientAddress,
		strings.TrimSpace(name),
	)
	if _, err := writer.Write([]byte(message)); err != nil {
		_ = writer.Close()
		return fmt.Errorf("write registration email: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("send registration email: %w", err)
	}
	if err := client.Quit(); err != nil {
		return fmt.Errorf("close Gmail SMTP session: %w", err)
	}

	return nil
}

func parseEmailAddress(address string) (string, error) {
	address = strings.TrimSpace(address)
	if address == "" || strings.ContainsAny(address, "\r\n<> ") || !strings.Contains(address, "@") {
		return "", fmt.Errorf("must be a plain email address")
	}
	return address, nil
}
