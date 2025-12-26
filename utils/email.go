package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
	"time"
)

func SendEmailWithCSV(to string, subject string, body string, filename string, csvData []byte) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SMTP_SENDER_EMAIL")
	senderPassword := os.Getenv("SMTP_APP_PASSWORD")

	boundary := fmt.Sprintf("BOUNDARY-%d", time.Now().UnixNano())

	var msg bytes.Buffer

	msg.WriteString(fmt.Sprintf("From: %s\r\n", senderEmail))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
	msg.WriteString("\r\n")

	msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	msg.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	msg.WriteString("Content-Transfer-Encoding: 7bit\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body + "\r\n")

	msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	msg.WriteString("Content-Type: text/csv\r\n")
	msg.WriteString("Content-Transfer-Encoding: base64\r\n")
	msg.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", filename))
	msg.WriteString("\r\n")

	encodedCSV := base64.StdEncoding.EncodeToString(csvData)

	for i := 0; i < len(encodedCSV); i += 76 {
		end := i + 76
		if end > len(encodedCSV) {
			end = len(encodedCSV)
		}

		msg.WriteString(encodedCSV[i:end] + "\r\n")
	}

	msg.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	auth := smtp.PlainAuth(
		"",
		senderEmail,
		senderPassword,
		smtpHost,
	)

	return smtp.SendMail(
		smtpHost + ":" + smtpPort,
		auth,
		senderEmail,
		[]string{to},
		msg.Bytes(),
	)
}