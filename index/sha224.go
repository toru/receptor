package index

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

const (
	SHA224DexIDLen    = sha256.Size224
	SHA224DexIDHexLen = sha256.Size224 * 2
)

// DexID is a 224-bit (or 28-byte) long primary key that every
// object inside Dexter has been alloted with.
type DexID = [sha256.Size224]byte

// NewDexIDFromString returns a new DexID based on the given string.
func NewDexIDFromString(src string) DexID {
	return sha256.Sum224([]byte(src))
}

// NewDexIDFromHexDigest returns a new DexID based on the given hex digest.
func NewDexIDFromHexDigest(src string) (DexID, error) {
	if !ValidateHexID(src) {
		return DexID{}, errors.New("invalid dexter id")
	}
	rv := DexID{}
	raw, err := hex.DecodeString(src)
	if err != nil {
		return rv, err
	}
	copy(rv[:], raw)
	return rv, nil
}

// DexIDToHexDigest returns the hexadecimal representation of the given
// DexID as a string. Sadly, DexID can't be used as a method receiver.
func DexIDToHexDigest(id []byte) string {
	return hex.EncodeToString(id)
}

// ValidateHexID returns a boolean indicating the validity of the given
// hexadecimal string. Mostly syntax sugar at this point.
func ValidateHexID(digest string) bool {
	return len(digest) == DexHexIDLen
}

type SHA224DexID struct {
	value [sha256.Size224]byte
}

func NewSHA224DexIDFromString(src string) ID {
	return &SHA224DexID{sha256.Sum224([]byte(src))}
}

// Algo implements the ID interface.
func (id SHA224DexID) Algo() uint8 {
	return DexIDTypeSHA224
}

// Value implements the ID interface.
func (id SHA224DexID) Value() []byte {
	return id.value[:]
}

// String implements the ID interface.
func (id SHA224DexID) String() string {
	return hex.EncodeToString(id.value[:])
}

// SetValue implements the ID interface.
func (id *SHA224DexID) SetValue(val []byte) {
	copy(id.value[:], val)
}

// SetValueFromHex implements the ID interface.
func (id *SHA224DexID) SetValueFromHexString(val string) error {
	if len(val) != SHA224DexIDHexLen {
		return errors.New("invalid hex string")
	}
	raw, err := hex.DecodeString(val)
	if err != nil {
		return err
	}
	copy(id.value[:], raw)
	return nil
}