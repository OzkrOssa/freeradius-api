package handler

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

type errorResponse struct {
	Success  bool     `json:"success" example:"false"`
	Messages []string `json:"messages" example:"Error message 1, Error message 2"`
}

func newErrorResponse(messages []string) errorResponse {
	return errorResponse{
		Success:  false,
		Messages: messages,
	}
}

func parseError(err error) []string {
	var errorMessages []string

	if errors.As(err, &validator.ValidationErrors{}) {
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, err.Error())
		}
	} else {
		errorMessages = append(errorMessages, err.Error())
	}

	return errorMessages
}
