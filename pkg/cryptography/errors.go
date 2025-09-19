package cryptography

import (
	"fmt"
)

func EncryptMessageError(err error) error {
	return fmt.Errorf("error encrypting message %w", err)
}

func DecryptMessageError(err error) error {
	return fmt.Errorf("error decrypting message %w", err)
}
