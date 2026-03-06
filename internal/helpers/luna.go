package helpers

// ValidLuhn проверяет корректность номера карты/числа по алгоритму Луна (Luhn).
// Принимает строку number, содержащую только цифры.
// Возвращает true, если число проходит проверку, иначе false.
func ValidLuhn(number string) bool {
	if len(number) == 0 {
		return false
	}

	var sum int
	alt := false

	for i := len(number) - 1; i >= 0; i-- {
		d := int(number[i] - '0')
		if d < 0 || d > 9 {
			return false
		}

		if alt {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		alt = !alt
	}

	return sum%10 == 0
}
