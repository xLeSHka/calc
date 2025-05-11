package password

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
)

const keyLength = 32

func Encrypt(data, key []byte) ([]byte, error) {
	if len(key) != keyLength {
		return nil, fmt.Errorf("invalid key length has been transmitted, recieve len %d", len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	cipherData := aesGCM.Seal(nonce, nonce, data, nil)
	return cipherData, err
}
func Decrypt(data, key []byte) ([]byte, error) {
	if len(key) != keyLength {
		return nil, fmt.Errorf("invalid key length has been transmitted, recieve len %d", len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("the length of ciphertext is less than the nonce size")
	}
	nonce, data := data[:nonceSize], data[nonceSize:]
	plainText, err := aesGCM.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, err
	}
	return plainText, err
}
func Compare(password, expected, key []byte) error {
	decrypted, err := Decrypt(expected, key)
	if err != nil {
		return err
	}
	if !bytes.Equal(decrypted, password) {
		log.Println("decryted")
		log.Println(string(decrypted))
		log.Println("password")
		log.Println(string(password))
		return errors.New("passwords do not match")
	}
	return nil
}
