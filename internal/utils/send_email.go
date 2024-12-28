package utils

import (
	"errors"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"strconv"
)

func SendMail(toParam string, MailData string, Subject string) error {
	err := godotenv.Load(".env")
	if err != nil {
		return errors.New("error loading .env file")
	}

	MailAddress := os.Getenv("MAIL_ADDRESS")
	if MailAddress == "" {
		return errors.New("MailAddress is empty")
	}

	MailPassword := os.Getenv("MAIL_PASSWORD")
	if MailPassword == "" {
		return errors.New("MailPassword is empty")
	}

	MailServiceAddress := os.Getenv("MAIL_SERVICE_ADDRESS")
	if MailServiceAddress == "" {
		return errors.New("MailServiceAddress is empty")
	}

	MailServicePort := os.Getenv("MAIL_SERVICE_PORT")
	if MailServicePort == "" {
		return errors.New("MailServicePort is empty")
	}

	from := MailAddress
	password := MailPassword
	host := MailServiceAddress
	port := MailServicePort

	intPort, err := strconv.Atoi(port)
	if err != nil {
		return errors.New("port is empty")
	}

	subject := Subject
	body := MailData

	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", toParam)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)
	n := gomail.NewDialer(host, intPort, from, password)
	if err = n.DialAndSend(msg); err != nil {
		log.Printf(err.Error())
	}

	return nil
}
