package contactform

import (
	"context"
	"fmt"

	"github.com/ippoippo/ippoippophotography-com-functions-contact/api"
	"github.com/ippoippo/ippoippophotography-com-functions-contact/configuration"
	"github.com/ippoippo/ippoippophotography-com-functions-contact/mailer"
	"github.com/ippoippo/ippoippophotography-com-functions-contact/validation"
)

type ContactForm interface {
	Execute(ctx context.Context, emailFormReq api.EmailFormRequest) api.EmailFormResponse
}

type ContactFormImpl struct {
	configuration *configuration.ContactFormConfiguration
	validator     validation.Validator
	mailer        mailer.Mailer
}

func NewContactFormImpl(
	configuration *configuration.ContactFormConfiguration,
	validator validation.Validator,
	mailer mailer.Mailer) *ContactFormImpl {
	return &ContactFormImpl{
		configuration: configuration,
		validator:     validator,
		mailer:        mailer,
	}
}

func (cf *ContactFormImpl) Execute(ctx context.Context, emailFormReq *api.EmailFormRequest) api.EmailFormResponse {
	if !cf.configuration.Valid() {
		fmt.Printf("Configuration: [%v]", cf.configuration)
		return api.InternalFailureResponse("configuration is invalid")
	}

	if cf.validator == nil {
		fmt.Println("Validator was nil")
		return api.InternalFailureResponse("validator is invalid")
	}

	cf.validator.Check(emailFormReq)
	if !cf.validator.Valid() {
		return api.ValidationFailureResponse(cf.validator.GlobalError(), cf.validator.FieldErrors())
	}

	if cf.mailer == nil {
		fmt.Println("Mailer was nil")
		return api.InternalFailureResponse("mailer is invalid")
	}

	err := cf.mailer.SendEmail(emailFormReq)
	if err != nil {
		return api.InternalFailureResponse(err.Error())
	}

	return api.SuccessResponse()
}
