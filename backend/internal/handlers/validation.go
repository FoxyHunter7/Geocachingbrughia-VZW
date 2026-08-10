package handlers

import (
	"net/mail"
	"strings"
)

const (
	maxStringLength = 10000
	maxTitleLength  = 200
	maxPriceCents   = 100000000
	maxQuantity     = 999
	maxStock        = 999999
)

func validateEmail(email string) bool {
	email = strings.TrimSpace(email)
	if len(email) > 254 || email == "" {
		return false
	}
	_, err := mail.ParseAddress(email)
	return err == nil
}

func truncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen]
	}
	return s
}

func sanitizeString(s string) string {
	s = strings.TrimSpace(s)
	return truncateString(s, maxStringLength)
}

func validatePriceCents(cents int) bool {
	return cents > 0 && cents <= maxPriceCents
}

func validateQuantity(qty int) bool {
	return qty > 0 && qty <= maxQuantity
}

func validateStock(stock *int) bool {
	if stock == nil {
		return true
	}
	return *stock >= 0 && *stock <= maxStock
}

var validCountryCodes = map[string]bool{
	"BE": true, "NL": true, "FR": true, "DE": true, "LU": true,
	"GB": true, "ES": true, "IT": true, "PT": true, "AT": true,
	"CH": true, "IE": true, "DK": true, "SE": true, "NO": true,
	"FI": true, "PL": true, "CZ": true, "SK": true, "HU": true,
	"RO": true, "BG": true, "HR": true, "SI": true, "EE": true,
	"LV": true, "LT": true, "GR": true,
}

func validateCountryCode(code string) bool {
	return validCountryCodes[strings.ToUpper(code)]
}
