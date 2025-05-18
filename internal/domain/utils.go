package domain

import (
	"crypto/rand"
	"math/big"
)

// cryptoRandShuffler provides a Shuffle method using crypto/rand.
type cryptoRandShuffler struct{}

// Shuffle shuffles a collection of n elements using the Fisher-Yates algorithm.
// The swap function is called with indices i and j to perform the swap.
// Random numbers are generated using crypto/rand for cryptographic security.
// It returns an error if random number generation fails.
func (s *cryptoRandShuffler) Shuffle(n int, swap func(i, j int)) error {
	if n <= 1 {
		return nil
	}
	// Needed for crypto/rand.Int, ensure "math/big" is imported.
	// Needed for rand.Reader and rand.Int, ensure "crypto/rand" is imported.
	for i := n - 1; i > 0; i-- {
		// Generate a random integer j such that 0 <= j <= i.
		// rand.Int returns a uniform random value in [0, max-1).
		// So, for the range [0, i], max is i+1.
		bigJ, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			// A failure in crypto/rand is critical for secure shuffling.
			return err
		}
		j := bigJ.Int64()
		swap(i, int(j))
	}
	return nil
}

// IntN returns a cryptographically secure random integer in [0, n).
// It returns an error if n <= 0 or if random number generation fails.
func (s *cryptoRandShuffler) IntN(n int) (int, error) {
	if n <= 0 {
		return 0, &invalidArgumentError{message: "argument to IntN must be positive"}
	}
	bigN := big.NewInt(int64(n))
	bigR, err := rand.Int(rand.Reader, bigN)
	if err != nil {
		return 0, err
	}
	return int(bigR.Int64()), nil
}

// invalidArgumentError is a custom error type for invalid arguments.
type invalidArgumentError struct {
	message string
}

// Error implements the error interface for invalidArgumentError.
func (e *invalidArgumentError) Error() string {
	return e.message
}

// Rng is an instance of cryptoRandShuffler.
// Calls to Rng.Shuffle will now use the crypto/rand-based implementation.
var Rng = &cryptoRandShuffler{}
