package api

import (
	"fmt"
	"net/http"
)

type EmailFormRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

type ResponseHeaders struct {
	ContentType string `json:"Content-Type"`
}

type FieldError struct {
	Field        string `json:"field"`
	ErrorMessage string `json:"errorMessage"`
}

type ResponseBody struct {
	Message            string       `json:"message"`
	GlobalErrorMessage string       `json:"globalErrorMessage"`
	FieldErrors        []FieldError `json:"fieldErrors"`
}

type EmailFormResponse struct {
	Body       ResponseBody    `json:"body"`
	StatusCode int             `json:"statusCode"`
	Headers    ResponseHeaders `json:"headers"`
}

func InternalFailureResponse(errorMessage string) EmailFormResponse {
	fmt.Printf("Internal Error: [%v]", errorMessage)
	res := baseResponse(http.StatusInternalServerError)
	res.Body = ResponseBody{
		GlobalErrorMessage: "Unexpected error occurred. Please try again later.",
		Message:            "error",
	}
	return res
}

func ValidationFailureResponse(globalError string, fieldErrors map[string]string) EmailFormResponse {
	res := baseResponse(http.StatusBadRequest)
	res.Body = ResponseBody{
		GlobalErrorMessage: globalError,
		FieldErrors:        toFieldErrors(fieldErrors),
		Message:            "error",
	}
	return res
}

func SuccessResponse() EmailFormResponse {
	res := baseResponse(http.StatusOK)
	res.Body = ResponseBody{
		Message: "success",
	}
	return res
}

func baseResponse(statusCode int) EmailFormResponse {
	return EmailFormResponse{
		StatusCode: statusCode,
		Headers: ResponseHeaders{
			ContentType: "application/json",
		},
	}
}

func toFieldErrors(errors map[string]string) []FieldError {
	var fieldErrors []FieldError
	for field, errorMessage := range errors {
		fieldErrors = append(fieldErrors, FieldError{
			Field:        field,
			ErrorMessage: errorMessage,
		})
	}
	return fieldErrors
}
