package contactform_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/ippoippo/ippoippophotography-com-functions-contact/api"
	"github.com/ippoippo/ippoippophotography-com-functions-contact/configuration"
	"github.com/ippoippo/ippoippophotography-com-functions-contact/contactform"
)

func TestNewContactFormImpl(t *testing.T) {
	_, cfg := setupValidConfiguration(t)

	cf := contactform.NewContactFormImpl(cfg, nil, nil)
	if cf == nil {
		t.Error("NewContactFormImpl() SHOULD NOT return nil")
	}
}

func TestExecuteMissingAPIKey(t *testing.T) {
	ctx := context.Background()
	cfg := configuration.NewContactFormConfiguration()

	cf := contactform.NewContactFormImpl(cfg, nil, nil)
	actual := cf.Execute(ctx, &api.EmailFormRequest{})
	expected := api.EmailFormResponse{
		StatusCode: 500,
		Headers: api.ResponseHeaders{
			ContentType: "application/json",
		},
		Body: api.ResponseBody{
			GlobalErrorMessage: "Unexpected error occurred. Please try again later.",
			Message:            "error",
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("cf.Execute() actual[%v], does not match expected[%v]", actual, expected)
	}
}

func TestExecuteNilValidator(t *testing.T) {
	ctx, cfg := setupValidConfiguration(t)

	cf := contactform.NewContactFormImpl(cfg, nil, nil)
	actual := cf.Execute(ctx, &api.EmailFormRequest{})
	expected := api.EmailFormResponse{
		StatusCode: 500,
		Headers: api.ResponseHeaders{
			ContentType: "application/json",
		},
		Body: api.ResponseBody{
			GlobalErrorMessage: "Unexpected error occurred. Please try again later.",
			Message:            "error",
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("cf.Execute() actual[%v], does not match expected[%v]", actual, expected)
	}
}

func TestExecuteInvalidRequest(t *testing.T) {
	ctx, cfg := setupValidConfiguration(t)

	mockedValidator := &MockContactFormValidator{
		ValidResult:       false,
		GlobalErrorResult: "global error from validator",
	}

	cf := contactform.NewContactFormImpl(cfg, mockedValidator, nil)
	actual := cf.Execute(ctx, &api.EmailFormRequest{})
	expected := api.EmailFormResponse{
		StatusCode: 400,
		Headers: api.ResponseHeaders{
			ContentType: "application/json",
		},
		Body: api.ResponseBody{
			GlobalErrorMessage: "global error from validator",
			Message:            "error",
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("cf.Execute() actual[%v], does not match expected[%v]", actual, expected)
	}
}

func TestExecuteFieldValidationFailure(t *testing.T) {
	ctx, cfg := setupValidConfiguration(t)

	mockedValidator := &MockContactFormValidator{
		ValidResult:       false,
		FieldErrorsResult: map[string]string{"name": "name must be between 1 and 100 characters"},
	}

	cf := contactform.NewContactFormImpl(cfg, mockedValidator, nil)
	actual := cf.Execute(ctx, &api.EmailFormRequest{})
	expected := api.EmailFormResponse{
		StatusCode: 400,
		Headers: api.ResponseHeaders{
			ContentType: "application/json",
		},
		Body: api.ResponseBody{
			FieldErrors: []api.FieldError{
				{
					Field:        "name",
					ErrorMessage: "name must be between 1 and 100 characters",
				},
			},
			Message: "error",
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("cf.Execute() actual[%v], does not match expected[%v]", actual, expected)
	}
}

func TestExecuteMailerIsNil(t *testing.T) {
	ctx, cfg := setupValidConfiguration(t)

	mockedValidator := &MockContactFormValidator{
		ValidResult: true,
	}

	cf := contactform.NewContactFormImpl(cfg, mockedValidator, nil)
	actual := cf.Execute(ctx, &api.EmailFormRequest{})
	expected := api.EmailFormResponse{
		StatusCode: 500,
		Headers: api.ResponseHeaders{
			ContentType: "application/json",
		},
		Body: api.ResponseBody{
			GlobalErrorMessage: "Unexpected error occurred. Please try again later.",
			Message:            "error",
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("cf.Execute() actual[%v], does not match expected[%v]", actual, expected)
	}
}

func TestExecuteMailerReturnsError(t *testing.T) {
	ctx, cfg := setupValidConfiguration(t)

	mockedValidator := &MockContactFormValidator{
		ValidResult: true,
	}

	mockedMailer := &MockMailer{
		SendEmailResult: errors.New("mailer error"),
	}

	cf := contactform.NewContactFormImpl(cfg, mockedValidator, mockedMailer)
	actual := cf.Execute(ctx, &api.EmailFormRequest{})
	expected := api.EmailFormResponse{
		StatusCode: 500,
		Headers: api.ResponseHeaders{
			ContentType: "application/json",
		},
		Body: api.ResponseBody{
			GlobalErrorMessage: "Unexpected error occurred. Please try again later.",
			Message:            "error",
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("cf.Execute() actual[%v], does not match expected[%v]", actual, expected)
	}
}

func TestExecuteSuccess(t *testing.T) {
	ctx, cfg := setupValidConfiguration(t)

	mockedValidator := &MockContactFormValidator{
		ValidResult: true,
	}

	mockedMailer := &MockMailer{
		SendEmailResult: nil,
	}

	cf := contactform.NewContactFormImpl(cfg, mockedValidator, mockedMailer)
	actual := cf.Execute(ctx, &api.EmailFormRequest{})
	expected := api.EmailFormResponse{
		StatusCode: 200,
		Headers: api.ResponseHeaders{
			ContentType: "application/json",
		},
		Body: api.ResponseBody{
			Message: "success",
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("cf.Execute() actual[%v], does not match expected[%v]", actual, expected)
	}
}

// Support functions

func setupValidConfiguration(t *testing.T) (context.Context, *configuration.ContactFormConfiguration) {
	t.Setenv("SENDGRID_API_KEY", "valid-api-key")
	ctx := context.Background()
	cfg := configuration.NewContactFormConfiguration()
	return ctx, cfg
}

// Mocks

type MockContactFormValidator struct {
	ValidResult       bool
	FieldErrorsResult map[string]string
	GlobalErrorResult string
}

func (v *MockContactFormValidator) Valid() bool {
	return v.ValidResult
}

func (v *MockContactFormValidator) FieldErrors() map[string]string {
	return v.FieldErrorsResult
}

func (v *MockContactFormValidator) GlobalError() string {
	return v.GlobalErrorResult
}

func (v *MockContactFormValidator) Check(_ any) {
}

type MockMailer struct {
	SendEmailResult error
}

func (m *MockMailer) SendEmail(_ *api.EmailFormRequest) error {
	return m.SendEmailResult
}
