package configuration

import (
	"os"
	"strings"
	"unicode/utf8"
)

type ContactFormConfiguration struct {
	SendGridApiKey string
}

func NewContactFormConfiguration() *ContactFormConfiguration {
	return &ContactFormConfiguration{SendGridApiKey: os.Getenv("SENDGRID_API_KEY")}
}

func (c *ContactFormConfiguration) Valid() bool {
	return utf8.RuneCountInString(strings.TrimSpace(c.SendGridApiKey)) != 0
}
