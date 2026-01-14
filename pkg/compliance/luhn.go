// Package compliance provides the Luhn algorithm for credit card validation.
package compliance

import "strings"

// ValidateLuhn validates a credit card number using the Luhn algorithm.
// Returns true if the number passes the Luhn check.
func ValidateLuhn(number string) bool {
	// Remove spaces and dashes
	cleaned := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, number)

	// Check length (credit cards are typically 13-19 digits)
	if len(cleaned) < 13 || len(cleaned) > 19 {
		return false
	}

	sum := 0
	alternate := false

	// Process from right to left
	for i := len(cleaned) - 1; i >= 0; i-- {
		digit := int(cleaned[i] - '0')

		if alternate {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
}

// GetCardType returns the card type based on the number prefix.
func GetCardType(number string) string {
	// Remove non-digits
	cleaned := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, number)

	if len(cleaned) < 4 {
		return "unknown"
	}

	// Check prefixes
	switch {
	case strings.HasPrefix(cleaned, "4"):
		return "visa"
	case strings.HasPrefix(cleaned, "51"), strings.HasPrefix(cleaned, "52"),
		strings.HasPrefix(cleaned, "53"), strings.HasPrefix(cleaned, "54"),
		strings.HasPrefix(cleaned, "55"):
		return "mastercard"
	case strings.HasPrefix(cleaned, "34"), strings.HasPrefix(cleaned, "37"):
		return "amex"
	case strings.HasPrefix(cleaned, "6011"), strings.HasPrefix(cleaned, "65"):
		return "discover"
	default:
		return "unknown"
	}
}

