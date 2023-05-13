package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/dimsonson/pswmanager/pkg/log"
)

type CryptProvider interface {
	EncryptAES(key, plaintext string) (string, error)
	DecryptAES(key, ciphertxt string) (string, error)
}

type Crypt struct{}

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

func (c *Crypt) DecryptAES(key, ciphertxt string) (string, error) {
	keyHex, _ := hex.DecodeString(key)
	ciphertext, _ := hex.DecodeString(ciphertxt)

	block, err := aes.NewCipher(keyHex)
	if err != nil {
		log.Print("decrypt error: ", err)
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		log.Print("decrypt error: ", err)
		return "", err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), err
	//return hex.EncodeToString(ciphertext), err

}

