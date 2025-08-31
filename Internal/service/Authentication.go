package service

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var emailRX = regexp.MustCompile(
	"(?i)^[a-z0-9!#$%&'*+/=?^_{|}~.-]+@" +
		"[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?" +
		"(?:\\.[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)+$",
)

func (service *Service) isValidPassword(pw string) bool {
	if len(pw) < 8 {
		return false
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, ch := range pw {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case regexp.MustCompile(`[!@#\$%\^&\*]`).MatchString(string(ch)):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}

func (service *Service) IsValidEmail(email string) bool {
	if len(email) == 0 || len(email) > 254 {
		return false
	}
	if !emailRX.MatchString(email) {
		return false
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	local := parts[0]
	domain := parts[1]

	
	if len(local) == 0 || len(local) > 64 {
		return false
	}
	if local[0] == '.' || local[len(local)-1] == '.' {
		return false 
	}
	if strings.Contains(local, "..") {
		return false 
	}

	if len(domain) == 0 || len(domain) > 253 {
		return false
	}
	if strings.Contains(domain, "..") {
		return false
	}

	labels := strings.Split(domain, ".")
	for _, lab := range labels {
		if len(lab) == 0 || len(lab) > 63 {
			return false
		}
		if lab[0] == '-' || lab[len(lab)-1] == '-' {
			return false
		}
	}

	tld := strings.ToLower(labels[len(labels)-1])
	if len(tld) < 2 || len(tld) > 63 {
		return false
	}
	for _, r := range tld {
		if r < 'a' || r > 'z' {
			return false
		}
	}
	return true
}

func hashPassword(password string) string {
	// use bcrypt to hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash)
}
