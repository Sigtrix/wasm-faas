package main

import (
	"math/big"
	"testing"
)

func TestModExp(t *testing.T) {
	base := big.NewInt(5)
	exp := big.NewInt(3)
	mod := big.NewInt(13)

	result := modExp(base, exp, mod)
	expected := big.NewInt(8)

	if result.Cmp(expected) != 0 {
		t.Errorf("ModExp(%s, %s, %s) = %s, expected %s", base, exp, mod, result, expected)
	}
}

func TestIsPrime(t *testing.T) {
	tests := []struct {
		number     *big.Int
		iterations int
		expected   bool
	}{
		{big.NewInt(1), 5, false},
		{big.NewInt(2), 5, true},
		{big.NewInt(17), 5, true},
		{big.NewInt(25), 5, false},
	}

	for _, test := range tests {
		result := isPrime(test.number, test.iterations)
		if result != test.expected {
			t.Errorf("isPrime(%s, %d) = %t, expected %t", test.number, test.iterations, result, test.expected)
		}
	}
}
