package validation

import (
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"

	"github.com/ippoippo/ippoippophotography-com-functions-contact/api"
)

const (
	MinNameLength    = 1
	MaxNameLength    = 100
	MinMessageLength = 1
	MaxMessageLength = 1000
	NameField        = "name"
	EmailField       = "email"
	MessageField     = "message"
)

// Define interface for validation
type Validator interface {
	Check(request any)
	Valid() bool
	FieldErrors() map[string]string
	GlobalError() string
}

type ContactFormValidator struct {
	fieldErrors map[string]string
	globalError string
}

func NewContactFormValidator() *ContactFormValidator {
	return &ContactFormValidator{}
}

func (v *ContactFormValidator) Valid() bool {
	return len(v.fieldErrors) == 0 && v.globalError == ""
}

func (v *ContactFormValidator) FieldErrors() map[string]string {
	return v.fieldErrors
}

func (v *ContactFormValidator) GlobalError() string {
	return v.globalError
}

func (v *ContactFormValidator) Check(request any) {
	efr, ok := request.(*api.EmailFormRequest)
	if !ok {
		fmt.Printf("request not expected type of *api.EmailFormRequest: [%v]", request)
		v.globalError = "invalid request type"
		return
	}

	v.checkField(validMinMaxChars(efr.Name, MinNameLength, MaxNameLength), NameField,
		validMinMaxCharsErrorMsg(NameField, MinNameLength, MaxNameLength))
	v.checkField(validMinMaxChars(efr.Message, MinMessageLength, MaxMessageLength), MessageField,
		validMinMaxCharsErrorMsg(MessageField, MinNameLength, MaxNameLength))

	v.checkField(validEmail(efr.Email), EmailField, validMailErrorMsg(EmailField))
}

func (v *ContactFormValidator) addFieldError(field, message string) {
	if v.fieldErrors == nil {
		v.fieldErrors = make(map[string]string)
	}
	if _, exists := v.fieldErrors[field]; !exists {
		v.fieldErrors[field] = message
	}
}

func (v *ContactFormValidator) checkField(ok bool, field, message string) {
	if !ok {
		v.addFieldError(field, message)
	}
}

func validMaxChars(value string, n int) bool {
	return utf8.RuneCountInString(strings.TrimSpace(value)) <= n
}

func validMinChars(value string, n int) bool {
	return utf8.RuneCountInString(strings.TrimSpace(value)) >= n
}

func validMinMaxChars(value string, min, max int) bool {
	trimmed := strings.TrimSpace(value)
	return validMinChars(trimmed, min) && validMaxChars(trimmed, max)
}

func validMinMaxCharsErrorMsg(field string, min, max int) string {
	return fmt.Sprintf("%s must be between %d and %d characters", field, min, max)
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func validMailErrorMsg(field string) string {
	return fmt.Sprintf("%s must be a valid email address", field)
}
