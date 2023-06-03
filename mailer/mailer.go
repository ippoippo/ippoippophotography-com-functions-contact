package mailer

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/ippoippo/ippoippophotography-com-functions-contact/api"
	"github.com/ippoippo/ippoippophotography-com-functions-contact/configuration"
)

var (
	acceptedSendStatusCodes = []int{200, 202}
	contactEmailAddress     = "contact@ippoippophotography.com"
	websiteUrl              = "https://ippoippophotography.com"
)

type Mailer interface {
	SendEmail(request *api.EmailFormRequest) error
}

type SendGridMailer struct {
	configuration *configuration.ContactFormConfiguration
}

func NewSendGridMailer(cfg *configuration.ContactFormConfiguration) *SendGridMailer {
	return &SendGridMailer{configuration: cfg}
}

func (m *SendGridMailer) SendEmail(request *api.EmailFormRequest) error {
	message := buildMessage(request)
	client := sendgrid.NewSendClient(m.configuration.SendGridApiKey)
	emailSendResponse, err := client.Send(message)
	if err != nil {
		fmt.Printf("Error sending email: %v", err)
		return err
	}
	if !m.isAcceptedStatusCode(emailSendResponse.StatusCode) {
		fmt.Printf("Error sending email: %v", emailSendResponse.Body)
		return fmt.Errorf("error sending email: %v", emailSendResponse.Body)
	}
	return nil
}

func buildMessage(request *api.EmailFormRequest) *mail.SGMailV3 {
	from := mail.NewEmail(fmt.Sprintf("%s Contact Form", websiteUrl), contactEmailAddress)
	subject := fmt.Sprintf("Contact Message from %s", websiteUrl)
	to := mail.NewEmail("ippoippo Photography", contactEmailAddress)
	plainTextContent := request.Message
	message := mail.NewSingleEmailPlainText(from, subject, to, plainTextContent)
	message = message.SetReplyTo(mail.NewEmail(request.Name, request.Email))
	return message
}

func (m *SendGridMailer) isAcceptedStatusCode(statusCode int) bool {
	for _, acceptedStatusCode := range acceptedSendStatusCodes {
		if statusCode == acceptedStatusCode {
			return true
		}
	}
	return false
}
