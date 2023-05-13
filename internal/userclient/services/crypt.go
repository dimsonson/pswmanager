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

func (c *Crypt) EncryptAES(key, plaintext string) (string, error) {

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
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return hex.EncodeToString(ciphertext), err
	// log.Print(key, "key")
	// keyHex, err := hex.DecodeString(key)
	// if err != nil {
	// 	log.Print("encrypt error: ", err)
	// 	return "", err
	// }

	// log.Print(keyHex)

	// cipher, err := aes.NewCipher(keyHex)
	// if err != nil {
	// 	log.Print("encrypt error: ", err)
	// 	return "", err
	// }

	// out := make([]byte, len(plaintext))
	// cipher.Encrypt(out, []byte(plaintext))
}

func (c *Crypt) DecryptAES(key, ciphertxt string) (string, error) {
	ciphertext, _ := hex.DecodeString(ciphertxt)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
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

	// keyHex, _ := hex.DecodeString(key)
	// ciphertext, _ := hex.DecodeString(ciphertxt)
	// cipher, err := aes.NewCipher(keyHex)
	// if err != nil {
	// 	log.Print("decrypt error: ", err)
	// 	return "", err
	// }
	// pt := make([]byte, len(ciphertext))
	// cipher.Decrypt(pt, ciphertext)
	// s := string(pt[:])

}

// // RandProvider интерфейс для вызова метода генератора псевдо случайной последовательности знаков.
// type RandProvider interface {
// 	RandSeq(n int) (random string, ok error)
// }

// // Rand структура для вызова метода генератора псевдо случайной последовательности знаков.
// type Rand struct{}

// // RandSeq функция генерации псевдо случайной последовательности знаков.
// func (r *Rand) RandSeq(n int) (random string, ok error) {
// 	if n < 1 {
// 		err := fmt.Errorf("wromg argument: number %v less than 1\n ", n)
// 		return "", err
// 	}
// 	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
// 	rand.Seed(time.Now().UnixNano())
// 	b := make([]rune, n)
// 	for i := range b {
// 		b[i] = letters[rand.Intn(len(letters))]
// 	}
// 	random = string(b)
// 	return random, nil
// }
