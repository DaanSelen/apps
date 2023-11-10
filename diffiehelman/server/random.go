package main

import (
	"crypto/rand"
	"math/big"
)

func isPrime(n *big.Int) bool {
	// 0 and 1 are not prime
	if n.Cmp(big.NewInt(2)) < 0 {
		return false
	}
	// Check for divisibility by numbers from 2 to the square root of n
	sqrtN := new(big.Int).Sqrt(n)
	for i := big.NewInt(2); i.Cmp(sqrtN) <= 0; i.Add(i, big.NewInt(1)) {
		if new(big.Int).Mod(n, i).Cmp(big.NewInt(0)) == 0 {
			return false
		}
	}
	return true
}

func generateRandomPrime(min, max *big.Int) (*big.Int, error) {
	for {
		// Generate a random number between min and max
		randomNum, err := rand.Int(rand.Reader, new(big.Int).Sub(max, min))
		if err != nil {
			return nil, err
		}
		// Add min to the random number to ensure it's in the desired range
		randomNum.Add(randomNum, min)

		// Check if the random number is prime
		if isPrime(randomNum) {
			return randomNum, nil
		}
	}
}
