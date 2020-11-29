package repositories

import (
	"log"

	"gopkg.in/gomail.v2"
)

func HTML(otp string) string {
	return `
		<html>
			<body>
				<div>send mail ` + otp + `</div>
			</body>
		</html>
	`
}

func SendEmail(Subject string, email string, otp string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", getEnv("MAILER_EMAIL", ""))
	m.SetHeader("To", email)
	m.SetHeader("Subject", Subject)
	html := HTML(otp)
	m.SetBody("text/html", html)

	d := gomail.NewDialer("smtp.gmail.com", 465, getEnv("MAILER_EMAIL", ""), getEnv("MAILER_PASSWORD", ""))
	log.Print(getEnv("MAILER_EMAIL", ""))
	log.Print(getEnv("MAILER_PASSWORD", ""))
	if err := d.DialAndSend(m); err != nil {
		log.Print(err)
	}
	return nil
}
