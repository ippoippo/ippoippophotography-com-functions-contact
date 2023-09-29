package validation_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/ippoippo/ippoippophotography-com-functions-contact/api"
	"github.com/ippoippo/ippoippophotography-com-functions-contact/validation"
)

func TestNewContactFormValidator(t *testing.T) {
	actual := validation.NewContactFormValidator()
	if actual == nil {
		t.Error("NewContactFormValidator() should result in non-nil")
	}
}

func TestContactFormValidAndFieldErrorsCount(t *testing.T) {
	type testSpec struct {
		request                  api.EmailFormRequest
		expectedValid            bool
		expectedFieldErrorsCount int
	}

	testSpecs := []testSpec{
		{
			request: api.EmailFormRequest{
				Name:    "",
				Email:   "",
				Message: "",
			},
			expectedValid:            false,
			expectedFieldErrorsCount: 3,
		},
		{
			request: api.EmailFormRequest{

				Name:    "Gavin Thomas",
				Email:   "",
				Message: "",
			},
			expectedValid:            false,
			expectedFieldErrorsCount: 2,
		},
		{
			request: api.EmailFormRequest{

				Name:    "Gavin Thomas",
				Email:   "",
				Message: "This is a test message.",
			},
			expectedValid:            false,
			expectedFieldErrorsCount: 1,
		},
		{
			request: api.EmailFormRequest{

				Name:    "Gavin Thomas",
				Email:   "test@example.com",
				Message: "",
			},
			expectedValid:            false,
			expectedFieldErrorsCount: 1,
		},
		{
			request: api.EmailFormRequest{

				Name:    "Gavin Thomas",
				Email:   "test@example.com",
				Message: " ",
			},
			expectedValid:            false,
			expectedFieldErrorsCount: 1,
		},
		{
			request: api.EmailFormRequest{

				Name:    "Gavin Thomas",
				Email:   "test@example.com",
				Message: "This is a test message.",
			},
			expectedValid:            true,
			expectedFieldErrorsCount: 0,
		},
	}

	for _, test := range testSpecs {
		validator := validation.ContactFormValidator{}
		validator.Check(&test.request)
		if actual := validator.Valid(); actual != test.expectedValid {
			t.Errorf("Valid() output [%v] not equal to expected [%v]", actual, test.expectedValid)
		}
		if actual := len(validator.FieldErrors()); actual != test.expectedFieldErrorsCount {
			t.Errorf("len(validator.FieldErrors()) output [%v] not equal to expected [%v]", actual, test.expectedFieldErrorsCount)
		}

	}
}

func TestContactFormRequestTypeValidation(t *testing.T) {
	type otherRequest struct {
		other int
	}

	validator := validation.ContactFormValidator{}
	validator.Check(&otherRequest{
		other: 10,
	})
	if validator.Valid() {
		t.Error("Valid() unexpectedly returned true")
	}
	if ge := validator.GlobalError(); ge != "invalid request type" {
		t.Errorf("validator.GlobalError(): expected[%s], got [%s]", "invalid request type", ge)
	}
}

func TestContactFormNameValidation(t *testing.T) {
	type testSpec struct {
		request          api.EmailFormRequest
		expectedErrorMsg string
	}

	testSpecs := []testSpec{
		{
			request: api.EmailFormRequest{
				Name:    "",
				Email:   "test@example.com",
				Message: "This is a test message.",
			},
			expectedErrorMsg: "name must be between 1 and 100 characters",
		},
		{
			request: api.EmailFormRequest{
				Name:    " ",
				Email:   "test@example.com",
				Message: "This is a test message.",
			},
			expectedErrorMsg: "name must be between 1 and 100 characters",
		},
		{
			request: api.EmailFormRequest{
				Name:    generateStringWithLength(101),
				Email:   "test@example.com",
				Message: "This is a test message.",
			},
			expectedErrorMsg: "name must be between 1 and 100 characters",
		},
		{
			request: api.EmailFormRequest{
				Name:    "Valid Name",
				Email:   "test@example.com",
				Message: "This is a test message.",
			},
			expectedErrorMsg: "",
		},
	}

	for _, test := range testSpecs {
		validator := validation.ContactFormValidator{}
		validator.Check(&test.request)
		actual, ok := validator.FieldErrors()["name"]
		if ok && actual != test.expectedErrorMsg {
			t.Errorf("error output [%v] not equal to expected [%v]", actual, test.expectedErrorMsg)
		}
		if !ok && test.expectedErrorMsg != "" {
			t.Errorf("error output missing error")
		}
	}
}

func TestContactFormMessageValidation(t *testing.T) {
	type testSpec struct {
		request          api.EmailFormRequest
		expectedErrorMsg string
	}

	testSpecs := []testSpec{
		{
			request: api.EmailFormRequest{
				Name:    "Gavin Thomas",
				Email:   "test@example.com",
				Message: "",
			},
			expectedErrorMsg: "message must be between 1 and 100 characters",
		},
		{
			request: api.EmailFormRequest{
				Name:    "Gavin Thomas",
				Email:   "test@example.com",
				Message: " ",
			},
			expectedErrorMsg: "message must be between 1 and 100 characters",
		},
		{
			request: api.EmailFormRequest{
				Name:    "Gavin Thomas",
				Email:   "test@example.com",
				Message: generateStringWithLength(1001),
			},
			expectedErrorMsg: "message must be between 1 and 100 characters",
		},
		{
			request: api.EmailFormRequest{
				Name:    "Gavin Thomas",
				Email:   "test@example.com",
				Message: "Valid Message",
			},
			expectedErrorMsg: "",
		},
	}

	for _, test := range testSpecs {
		validator := validation.ContactFormValidator{}
		validator.Check(&test.request)
		actual, ok := validator.FieldErrors()["message"]
		if ok && actual != test.expectedErrorMsg {
			t.Errorf("error output [%v] not equal to expected [%v]", actual, test.expectedErrorMsg)
		}
		if !ok && test.expectedErrorMsg != "" {
			t.Errorf("error output missing error")
		}
	}
}

func TestContactFormEmailValidation(t *testing.T) {
	type testSpec struct {
		request          api.EmailFormRequest
		expectedErrorMsg string
	}

	testSpecs := []testSpec{
		{
			request: api.EmailFormRequest{
				Name:    "Gavin Thomas",
				Email:   "",
				Message: "Valid Message",
			},
			expectedErrorMsg: "email must be a valid email address",
		},
		{
			request: api.EmailFormRequest{
				Name:    "Gavin Thomas",
				Email:   " ",
				Message: "Valid Message",
			},
			expectedErrorMsg: "email must be a valid email address",
		},
		{
			request: api.EmailFormRequest{
				Name:    "Gavin Thomas",
				Email:   "bad-example",
				Message: "Valid Message",
			},
			expectedErrorMsg: "email must be a valid email address",
		},
		{
			request: api.EmailFormRequest{
				Name:    "Gavin Thomas",
				Email:   "bad-example.com",
				Message: "Valid Message",
			},
			expectedErrorMsg: "email must be a valid email address",
		},
		{
			request: api.EmailFormRequest{
				Name:    "Gavin Thomas",
				Email:   "test@example.com",
				Message: "Valid Message",
			},
			expectedErrorMsg: "",
		},
	}

	for _, test := range testSpecs {
		validator := validation.ContactFormValidator{}
		validator.Check(&test.request)
		actual, ok := validator.FieldErrors()["email"]
		if ok && actual != test.expectedErrorMsg {
			t.Errorf("error output [%v] not equal to expected [%v]", actual, test.expectedErrorMsg)
		}
		if !ok && test.expectedErrorMsg != "" {
			t.Errorf("error output missing error")
		}
	}
}

func generateStringWithLength(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
