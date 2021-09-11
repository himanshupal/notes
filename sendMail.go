package main

import (
	"net/smtp"
	"os"
)

type Provider struct {
	Host    string `json:"host"`
	Address string `json:"address"`
}

type Providers struct {
	Google Provider
}

var EmailProviders Providers = Providers{
	Google: Provider{
		Host:    "smtp.gmail.com",
		Address: "smtp.gmail.com:587",
	},
}

// Sends plaintext Email to given email address
func SendMail(provider Provider, to string, subject string, body []byte) error {
	emailAuth := smtp.PlainAuth("", os.Getenv("EMAIL_USERNAME"), os.Getenv("EMAIL_PASSWORD"), provider.Host)

	return smtp.SendMail(provider.Address, emailAuth, os.Getenv("EMAIL_FROM"), []string{to}, body)
}
