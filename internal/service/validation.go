package service

import (
	"net/mail"
	"unicode"
)

func validateRegisterInput(email, password string) error {
	if !isValidEmail(email) {
		return badRequest("Укажите корректную электронную почту.")
	}
	return validatePassword(password)
}

func validateLoginInput(email, password string) error {
	if !isValidEmail(email) {
		return badRequest("Укажите корректную электронную почту.")
	}
	if password == "" {
		return badRequest("Введите пароль.")
	}
	return nil
}

func validatePassword(password string) error {
	if len([]rune(password)) < 8 {
		return badRequest("Пароль должен быть не короче 8 символов.")
	}

	hasLetter := false
	hasDigit := false
	for _, r := range password {
		if unicode.IsLetter(r) {
			hasLetter = true
		}
		if unicode.IsDigit(r) {
			hasDigit = true
		}
	}

	if !hasLetter || !hasDigit {
		return badRequest("Пароль должен содержать буквы и цифры.")
	}

	return nil
}

func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	addr, err := mail.ParseAddress(email)
	return err == nil && addr.Address == email
}
