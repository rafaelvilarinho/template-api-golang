package libraries

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"api.template.com.br/helpers"
)

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decode(s string) (*[]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func Encrypt(text string) (string, error) {
	env, _ := helpers.GetEnvironment()
	block, err := aes.NewCipher([]byte(env.CRYPT_SECRET))
	if err != nil {
		return "", err
	}

	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)

	return encode(cipherText), nil
}

func Decrypt(text string) (string, error) {
	env, _ := helpers.GetEnvironment()
	block, err := aes.NewCipher([]byte(env.CRYPT_SECRET))
	if err != nil {
		return "", err
	}

	cipherText, err := decode(text)
	if err != nil {
		return "", err
	}

	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(*cipherText))
	cfb.XORKeyStream(plainText, *cipherText)

	return string(plainText), nil
}
