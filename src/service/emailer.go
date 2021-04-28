package service

import (
	"fmt"
	"gopkg.in/gomail.v2"
	cnf "../../configuration"
)

func SendEmail(
	conf *cnf.Configuration,
	from string,
	to map[string]string,
	cc map[string]string,
	bcc map[string]string,
	subject string,
	message string,
	isHTML bool,
	files map[int]string) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	for _, toEmail := range to {
		m.SetHeader("To", toEmail)
	}
	for _, ccEmail := range cc {
		m.SetHeader("Cc", ccEmail)
	}
	for _, bccEmail := range bcc {
		m.SetHeader("Bcc", bccEmail)
	}

	m.SetHeader("Subject", subject)

	if isHTML {
		m.SetBody("text/html", message)
	} else {
		m.SetBody("text/plain", message)
	}

	for _, file := range files {
		m.Attach(file)
	}

	//go get gopkg.in/gomail.v2
	d := gomail.NewPlainDialer(conf.SMTP_Server, conf.SMTP_Port, conf.SMTP_User, conf.SMTP_Pass)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err);
		return err
	}

	return nil
}
