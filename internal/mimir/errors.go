package mimir

import "fmt"

type ValueNotFoundError struct {
	Field string
}

func (e ValueNotFoundError) Error() string {
	return fmt.Sprintf("%s field is missing", e.Field)
}

type WrongFormatError struct {
	Field string
}

func (e WrongFormatError) Error() string {
	return fmt.Sprintf("Wrong format for field %s", e.Field)
}
