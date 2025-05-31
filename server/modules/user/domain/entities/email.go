package entities

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"
)

// Email represents an email value object
type Email struct {
	value string
}

// emailRegex is a simple email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmail creates a new Email value object
func NewEmail(value string) (Email, error) {
	if value == "" {
		return Email{}, fmt.Errorf("email cannot be empty")
	}

	value = strings.TrimSpace(strings.ToLower(value))

	if !emailRegex.MatchString(value) {
		return Email{}, fmt.Errorf("invalid email format: %s", value)
	}

	return Email{value: value}, nil
}

// String returns the string representation of the email
func (e Email) String() string {
	return e.value
}

// GetValue returns the email value
func (e Email) GetValue() string {
	return e.value
}

// Equals checks if two emails are equal
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// Scan implements the sql.Scanner interface for GORM
func (e *Email) Scan(value interface{}) error {
	if value == nil {
		*e = Email{}
		return nil
	}

	switch v := value.(type) {
	case string:
		email, err := NewEmail(v)
		if err != nil {
			return err
		}
		*e = email
		return nil
	case []byte:
		email, err := NewEmail(string(v))
		if err != nil {
			return err
		}
		*e = email
		return nil
	default:
		return fmt.Errorf("cannot scan %T into Email", value)
	}
}

// Value implements the driver.Valuer interface for GORM
func (e Email) Value() (driver.Value, error) {
	return e.value, nil
}
