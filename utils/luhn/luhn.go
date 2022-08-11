package luhn

import (
	"fmt"
	"strconv"
)

// CalculateToLuhn return Lunh true number by adding 1 digit
func CalculateToLuhn(number int) int {
	if Valid(number) {
		return number
	}
	checkNumber := checksum(number)

	if checkNumber == 0 {
		strNumber := fmt.Sprintf("%d", number) + "0"
		finNumber, _ := strconv.Atoi(strNumber)
		return finNumber
	}

	strNumber := fmt.Sprintf("%d", number) + fmt.Sprintf("%d", 10-checkNumber)
	finNumber, _ := strconv.Atoi(strNumber)

	return finNumber
}

// Valid check number is valid or not based on Luhn algorithm
func Valid(number int) bool {
	return (number%10+checksum(number/10))%10 == 0
}

// checksum check control sum of digits
func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur *= 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number /= 10
	}
	return luhn % 10
}
