package util

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/k3a/html2text"

	"gopkg.in/gomail.v2"
)

type EmailData struct {
	FirstName string
	Subject   string
	Text      string
}

// ? Email template parser

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendEmail(emailFrom string, email string, data *EmailData) (err error) {
	// config, err := initializers.LoadConfig(".")
	// if err != nil {
	// 	log.Fatal("could not load config", err)
	// 	return
	// }

	from := emailFrom
	smtpPass := "e150afc3334f8c"
	smtpUser := "35ff331af2224f"
	to := email
	smtpHost := "sandbox.smtp.mailtrap.io"
	smtpPort := 2525
	var body bytes.Buffer

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", data.Text)
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err = d.DialAndSend(m); err != nil {
		log.Fatal("Could not send email: ", err)
		return
	}

	fmt.Println("mail sent")
	return
}

func IsValidEmail(email string) bool {
	// Regular expression pattern for email validation
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regular expression
	re := regexp.MustCompile(emailRegex)

	// Check if the email matches the pattern
	return re.MatchString(email)
}
