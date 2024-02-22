package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"
)

// modExp performs modular exponentiation
func modExp(base *big.Int, exp *big.Int, mod *big.Int) *big.Int {
	result := new(big.Int).Exp(base, exp, mod)
	return result
}

// witness determines the validity of a witness
// for the Miller-Rabin test
func witness(a *big.Int, n *big.Int) bool {
	one := big.NewInt(1)
	d := new(big.Int).Sub(n, one)
	s := 0

	for d.Bit(s) == 0 {
		s++
	}

	r := new(big.Int).Rsh(d, uint(s))
	x := modExp(a, r, n)

	if x.Cmp(one) == 0 || x.Cmp(new(big.Int).Sub(n, one)) == 0 {
		return false
	}

	for i := 0; i < s-1; i++ {
		x.Mod(x.Mul(x, x), n)
		if x.Cmp(one) == 0 {
			return true
		}
		if x.Cmp(new(big.Int).Sub(n, one)) == 0 {
			return false
		}
	}
	return true
}

// isPrime performs the Miller-Rabin primality test
func isPrime(n *big.Int, k int) bool {
	if n.Cmp(big.NewInt(1)) <= 0 || n.Cmp(big.NewInt(4)) == 0 {
		return false
	}
	if n.Cmp(big.NewInt(3)) <= 0 {
		return true
	}

	one := big.NewInt(1)
	two := big.NewInt(2)
	d := new(big.Int).Sub(n, one)

	for d.Bit(0) == 0 {
		d.Rsh(d, 1)
	}

	for i := 0; i < k; i++ {
		a, err := rand.Int(rand.Reader, new(big.Int).Sub(n, two))
		if err != nil {
			fmt.Println("Error generating random number:", err)
			return false
		}
		a.Add(a, two)

		if witness(a, n) {
			return false
		}
	}
	return true
}

// initEnvironment returns a map with the environment
// variables available in the sandboxed environment.
func initEnvironment() map[string]string {
	env := make(map[string]string)

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			key, value := pair[0], pair[1]
			env[key] = value
		}
	}
	return env
}

// parseQuery parses a http query and returns a map
// of parsed key-value pairs.
func parseQuery(query string) map[string]string {
	pairs := strings.Split(query, "&")
	kvMap := make(map[string]string)
	for _, pair := range pairs {
		parsedPair := strings.SplitN(pair, "=", 2)
		if len(parsedPair) == 2 {
			k, v := parsedPair[0], parsedPair[1]
			kvMap[k] = v
		}
	}
	return kvMap
}

func main() {
	number := new(big.Int)
	k := 5 // Number of iterations for the Miller-Rabin test

	env := initEnvironment()
	queryMap := parseQuery(env["http_query"])

	numberQuery, ok := queryMap["number"]
	if !ok {
		fmt.Println("No number provided  for primality test.")
		return
	}

	_, success := number.SetString(numberQuery, 10)
	if !success {
		fmt.Println("Error: Unable to set the big.Int from the given string.")
		return
	}

	if isPrime(number, k) {
		fmt.Println(number, "is likely prime.")
	} else {
		fmt.Println(number, "is composite.")
	}
}
