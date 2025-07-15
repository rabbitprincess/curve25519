package ed25519

import (
	"crypto/sha512"
	"fmt"
	"math/big"
	"slices"
)

var (
	curveL, _ = new(big.Int).SetString("1000000000000000000000000000000014def9dea2f79cd65812631a5cf5d3ed", 16)
)

func GetScalar(seed []byte) ([]byte, error) {
	if len(seed) != 32 {
		return nil, fmt.Errorf("seed must be 32 bytes long, got %d bytes", len(seed))
	}
	hash := sha512.Sum512(seed)
	privateKeyRaw := make([]byte, 32)
	copy(privateKeyRaw, hash[:32])

	privateKeyRaw[0] &= 248  // clear the last 3 bits
	privateKeyRaw[31] &= 127 // clear the highest bit
	privateKeyRaw[31] |= 64  // set the second highest bit

	slices.Reverse(privateKeyRaw) // reverse the byte order
	privateKey := new(big.Int).SetBytes(privateKeyRaw)
	privateKey.Mod(privateKey, curveL) // ensure the scalar is within the field modulus
	if privateKey.Sign() <= 0 {
		privateKey.Add(privateKey, curveL) // ensure positive scalar
	}

	// pad the scalar to 32 bytes
	scalarBytes := make([]byte, 32)
	privateKey.FillBytes(scalarBytes)

	return scalarBytes, nil
}
