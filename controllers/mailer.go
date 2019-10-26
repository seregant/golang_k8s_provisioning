package controllers

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/seregant/golang_k8s_provisioning/database"
	"github.com/seregant/golang_k8s_provisioning/models"
)

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_EMAIL = "indra.mailer@gmail.com"
const CONFIG_PASSWORD = "jfezijcndhvmdgns"

func sendNotif(mailAddr []string, message string) bool {
	var db = database.DbConnect()
	defer db.Close()

	var dataUser models.Pengguna

	db.Where("pengguna_email = ?", mailAddr[0]).First(&dataUser)

	subject := "Owncloud Admin"
	err := sendMail(mailAddr, subject, message)
	if err != nil {
		log.Fatal(err.Error())
		return false
	}

	log.Println("Mail sent!")
	return true
}

func sendMail(to []string, subject, message string) error {
	body := "From: " + CONFIG_EMAIL + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		message

	auth := smtp.PlainAuth("", CONFIG_EMAIL, CONFIG_PASSWORD, CONFIG_SMTP_HOST)
	smtpAddr := fmt.Sprintf("%s:%d", CONFIG_SMTP_HOST, CONFIG_SMTP_PORT)

	err := smtp.SendMail(smtpAddr, auth, CONFIG_EMAIL, to, []byte(body))
	if err != nil {
		return err
	}

	return nil
}
