package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"

	"github.com/dimsonson/pswmanager/pkg/log"
)

type CryptProvider interface {
	EncryptAES(key, plaintext string) (string, error)
	DecryptAES(key, ciphertxt string) (string, error)
}

type Crypt struct{}

// EncryptAES метод AES шифрования
func (c *Crypt) EncryptAES(key, plaintxt string) (string, error) {
	plaintext := []byte(plaintxt)
	keyHex, _ := hex.DecodeString(key)
	block, err := aes.NewCipher(keyHex)
	if err != nil {
		log.Print("encrypt error: ", err)
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Print("encrypt error: ", err)
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return hex.EncodeToString(ciphertext), err
}

// DecryptAES метод AES дешифрования
func (c *Crypt) DecryptAES(key, ciphertxt string) (string, error) {
	keyHex, _ := hex.DecodeString(key)
	ciphertext, _ := hex.DecodeString(ciphertxt)
	block, err := aes.NewCipher(keyHex)
	if err != nil {
		log.Print("decrypt error: ", err)
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		err = errors.New("wrong ciphertext, it is too small")
		log.Print("decrypt error: ", err)
		return "", err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	//return hex.EncodeToString(ciphertext), err
	return string(ciphertext), err
}
