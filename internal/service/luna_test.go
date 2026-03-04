package service

import "testing"

func TestValidLuna_ValidNumbers(t *testing.T) {
	validNumbers := []string{
		"4532015112830366", // Visa
		"6011000990139424", // Discover
		"378282246310005",  // Amex
	}

	for _, number := range validNumbers {
		if !ValidLuhn(number) {
			t.Errorf("ValidLuna(%s) = false, want true", number)
		}
	}
}

func TestValidLuna_InvalidNumbers(t *testing.T) {
	invalidNumbers := []string{
		"4532015112830367", // изменена последняя цифра
		"1234567890123456", // случайный набор цифр
		"abcd1234",         // содержит буквы
		"",                 // пустая строка
	}

	for _, number := range invalidNumbers {
		if ValidLuhn(number) {
			t.Errorf("ValidLuna(%s) = true, want false", number)
		}
	}
}
