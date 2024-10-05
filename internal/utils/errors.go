package utils

import "fmt"

type RequiredFieldError struct {
	Field string
}

func (e RequiredFieldError) Error() string {
	return fmt.Sprintf("%s field is missing", e.Field)
}

type WrongFormatError struct {
	Field string
}

func (e WrongFormatError) Error() string {
	return fmt.Sprintf("Wrong format for field %s", e.Field)
}