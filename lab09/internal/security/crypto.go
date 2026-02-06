package security

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/nacl/box"
)

// LoadKey parses a 64-character hex string into a [32]byte key.
func LoadKey(hexKey string) (*[32]byte, error) {
	b, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, err
	}
	if len(b) != 32 {
		return nil, fmt.Errorf("key must be 32 bytes, got %d", len(b))
	}
	var key [32]byte
	copy(key[:], b)
	return &key, nil
}

// Encrypt seals a message using the sender's private key and receiver's public key.
// Format: [Nonce (24 bytes)][Ciphertext]
func Encrypt(msg []byte, mySecret, peerPublic *[32]byte) ([]byte, error) {
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}

	// box.Seal appends the encrypted message to the first argument.
	// We start with the nonce so the receiver can extract it.
	encrypted := box.Seal(nonce[:], msg, &nonce, peerPublic, mySecret)
	return encrypted, nil
}

// Decrypt opens a sealed box.
// Expects: [Nonce (24 bytes)][Ciphertext]
func Decrypt(encrypted []byte, mySecret, peerPublic *[32]byte) ([]byte, error) {
	if len(encrypted) < 24 {
		return nil, fmt.Errorf("message too short")
	}

	var nonce [24]byte
	copy(nonce[:], encrypted[:24])
	ciphertext := encrypted[24:]

	decrypted, ok := box.Open(nil, ciphertext, &nonce, peerPublic, mySecret)
	if !ok {
		return nil, fmt.Errorf("decryption failed - invalid key or corrupted message")
	}
	return decrypted, nil
}
