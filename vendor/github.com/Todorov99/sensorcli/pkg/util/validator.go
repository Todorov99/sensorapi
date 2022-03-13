package util

import "fmt"

// ValidateFormat validate the chosen format.
func ValidateFormat(format string) error {
	if format != jsonType && format != yamlType {
		return fmt.Errorf("invalid format")
	}

	return nil
}
