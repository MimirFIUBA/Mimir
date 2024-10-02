package validators

import "fmt"

func CheckIdNotEmpty(id string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	return nil
}
