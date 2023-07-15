package helpers

import (
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	PORT             string
	WEBSITE_URL      string
	MONGODB_URI      string
	CRYPT_SECRET     string
	JWT_SECRET       string
	SENDGRID_API_KEY string
	MAIL_USER_NAME   string
	MAIL_USER_EMAIL  string
}

func GetEnvironment() (*Environment, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Environment{
		PORT:             os.Getenv("PORT"),
		WEBSITE_URL:      os.Getenv("WEBSITE_URL"),
		MONGODB_URI:      os.Getenv("MONGODB_URI"),
		CRYPT_SECRET:     os.Getenv("CRYPT_SECRET"),
		JWT_SECRET:       os.Getenv("JWT_SECRET"),
		SENDGRID_API_KEY: os.Getenv("SENDGRID_API_KEY"),
		MAIL_USER_NAME:   os.Getenv("MAIL_USER_NAME"),
		MAIL_USER_EMAIL:  os.Getenv("MAIL_USER_EMAIL"),
	}, nil
}
