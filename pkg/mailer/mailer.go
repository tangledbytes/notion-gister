package mailer

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
)

// Mailer struct
type Mailer struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
	To       []string
	Data     map[string]interface{}
}

// New returns a pointer to the instance of Mailer
func New(host, port, username, password, from string, to []string, data map[string]interface{}) *Mailer {
	return &Mailer{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
		To:       to,
		Data:     data,
	}
}

// Mail takes the subject and body and sends the mail
//
// subject: Subject template of the mail
// body: Body template of the mail
func (m *Mailer) Mail(subject string, body string) error {
	subject, body, err := m.prepareMail(subject, body)
	if err != nil {
		return err
	}

	return m.sendMail(subject, body)
}

func (m *Mailer) prepareMail(subject, body string) (string, string, error) {
	subjectTemplate, err := template.New("mail.subject").Parse(subject)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse subject template: %w", err)
	}

	bodyTemplate, err := template.New("mail.body").Parse(body)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse body template: %w", err)
	}

	var subjectBuffer strings.Builder
	if err := subjectTemplate.Execute(&subjectBuffer, m.Data); err != nil {
		return "", "", fmt.Errorf("failed to execute subject template: %w", err)
	}

	var bodyBuffer strings.Builder
	if err := bodyTemplate.Execute(&bodyBuffer, m.Data); err != nil {
		return "", "", fmt.Errorf("failed to execute body template: %w", err)
	}

	return subjectBuffer.String(), bodyBuffer.String(), nil
}

func (m *Mailer) sendMail(subject, message string) error {
	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)

	content := createContent(m.From, m.To, subject, message)

	logrus.Info("Sending:\n", content)

	conn, err := tls.Dial("tcp", m.Host+":"+m.Port, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         m.Host,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to smtp server: %w", err)
	}

	client, err := smtp.NewClient(conn, m.Host)
	if err != nil {
		return fmt.Errorf("failed to create smtp client: %w", err)
	}

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	if err := client.Mail(m.From); err != nil {
		return fmt.Errorf("failed to set from address: %w", err)
	}

	for _, to := range m.To {
		if err := client.Rcpt(to); err != nil {
			return fmt.Errorf("failed to set to address: %w", err)
		}
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = writer.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	err = client.Quit()
	if err != nil {
		return fmt.Errorf("failed to quit: %w", err)
	}

	return nil
}

func createContent(from string, to []string, subject, message string) string {
	return fmt.Sprintf(`To: %s
From: %s
Subject: %s
MIME-version: 1.0;
Content-Type: text/html; charset="UTF-8";

<html>
<body>
%s
</body>
</html>
`, strings.Join(to, ","), from, subject, message)
}
