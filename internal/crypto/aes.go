package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

func NewAESStream(key []byte) (cipher.Stream, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewCTR(block, key[:block.BlockSize()]), nil
}
