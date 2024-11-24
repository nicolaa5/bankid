package bankid

import (
	"fmt"
	"strconv"
)

// Swedish personal numbers are created with the "mod 10" algorithm, aka Luhn algorithm.
// The steps are as follows:
// 1. Start from the Right: Begin with the rightmost digit (check digit) and move leftwards.
// 2. Double Every Second Digit: Double the value of every second digit. If doubling results in a number greater than 9, subtract 9 from the result.
// 3. Sum All Digits: Add together all the digits, including those not doubled.
// 4. Calculate Modulo 10: The total sum modulo 10 should be 0 for a valid number.
func validateChecksum(number string) error {
	mod := len(number) % 2
	sum := 0
	for i, ch := range number {
		num, err := strconv.Atoi(string(ch))
		if err != nil {
			return fmt.Errorf("invalid number: %w", err)
		}
		if i%2 == mod {
			if num < 5 {
				sum += num * 2
			} else {
				sum += num*2 - 9
			}
		} else {
			sum += num
		}
	}
	if sum%10 != 0 {
		return fmt.Errorf("invalid checksum in personnummer")
	}
	return nil
}
