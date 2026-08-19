package domain

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
)

var ErrInvalidUUID = errors.New("invalid UUID")

// NewUUID creates a random RFC 4122 version 4 UUID for offline-safe IDs.
func NewUUID() (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("generate UUID: %w", err)
	}

	raw[6] = (raw[6] & 0x0f) | 0x40
	raw[8] = (raw[8] & 0x3f) | 0x80

	return formatUUID(raw), nil
}

// ValidateUUID accepts canonical lower-case UUID strings with the RFC 4122
// variant and a defined UUID version.
func ValidateUUID(value string) error {
	if len(value) != 36 || value[8] != '-' || value[13] != '-' || value[18] != '-' || value[23] != '-' {
		return ErrInvalidUUID
	}

	var compact [32]byte
	compactIndex := 0
	for index := 0; index < len(value); index++ {
		if index == 8 || index == 13 || index == 18 || index == 23 {
			continue
		}
		character := value[index]
		if !isLowerHex(character) {
			return ErrInvalidUUID
		}
		compact[compactIndex] = character
		compactIndex++
	}

	var raw [16]byte
	if _, err := hex.Decode(raw[:], compact[:]); err != nil {
		return ErrInvalidUUID
	}
	version := raw[6] >> 4
	if raw[8]&0xc0 != 0x80 || version == 0 || version > 8 {
		return ErrInvalidUUID
	}

	return nil
}

func formatUUID(raw [16]byte) string {
	encoded := make([]byte, 36)
	hex.Encode(encoded[0:8], raw[0:4])
	encoded[8] = '-'
	hex.Encode(encoded[9:13], raw[4:6])
	encoded[13] = '-'
	hex.Encode(encoded[14:18], raw[6:8])
	encoded[18] = '-'
	hex.Encode(encoded[19:23], raw[8:10])
	encoded[23] = '-'
	hex.Encode(encoded[24:36], raw[10:16])
	return string(encoded)
}

func isLowerHex(character byte) bool {
	return character >= '0' && character <= '9' || character >= 'a' && character <= 'f'
}
