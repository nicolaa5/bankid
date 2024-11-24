package bankid

import (
	"fmt"
	"testing"
)

func TestLuhn(t *testing.T) {
	personnummers := map[string]bool{
		"9512011294": false,
		"6310234224": false,
		"3810260632": true,
	}
	for num, expected := range personnummers {
		err := validateChecksum(num)
		actual := err == nil

		if expected != actual {
			t.Errorf("Expected %t actual %t for %s", expected, actual, num)

			if err != nil {
				fmt.Printf("checksum error: %v", err.Error())
			}
		}
	}
}
