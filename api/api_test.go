package api_test

import (
	"reflect"
	"testing"

	"github.com/ippoippo/ippoippophotography-com-functions-contact/api"
)

func TestSuccessResponse(t *testing.T) {
	actual := api.SuccessResponse()
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
		t.Errorf("SuccessResponse() actual[%v], does not match expected[%v]", actual, expected)
	}
}

func TestValidationFailureResponse(t *testing.T) {
	fieldErrors := map[string]string{
		"field1": "error message 1",
		"field2": "error message 2",
	}
	actual := api.ValidationFailureResponse("Global error message", fieldErrors)
	expected := api.EmailFormResponse{
		StatusCode: 400,
		Headers: api.ResponseHeaders{
			ContentType: "application/json",
		},
		Body: api.ResponseBody{
			GlobalErrorMessage: "Global error message",
			FieldErrors: []api.FieldError{
				{
					Field:        "field1",
					ErrorMessage: "error message 1",
				},
				{
					Field:        "field2",
					ErrorMessage: "error message 2",
				},
			},
			Message: "error",
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("SuccessResponse() actual[%v], does not match expected[%v]", actual, expected)
	}
}

func TestInternalFailureResponse(t *testing.T) {
	actual := api.InternalFailureResponse("Internal error message")
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
		t.Errorf("SuccessResponse() actual[%v], does not match expected[%v]", actual, expected)
	}
}
