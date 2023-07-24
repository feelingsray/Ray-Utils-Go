package secret

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type AESCrypt struct {
	key []byte
	iv  []byte
}

func NewAES() *AESCrypt {
	sc := new(AESCrypt)
	sc.key = []byte("Dv3KMrhJyxZiQU6p")
	sc.iv = []byte("0000000000000000")
	return sc
}

func (a *AESCrypt) encrypt(data string) (string, error) {
	//var block cipher.Block
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", err
	}
	encrypt := cipher.NewCBCEncrypter(block, a.iv)
	source := PKCS5Padding([]byte(data), block.BlockSize())
	dst := make([]byte, len(source))
	encrypt.CryptBlocks(dst, source)
	return base64.RawStdEncoding.EncodeToString(dst), nil
}

func (a *AESCrypt) decrypt(data string) (string, error) {
	//var block cipher.Block
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", err
	}
	encrypt := cipher.NewCBCDecrypter(block, a.iv)
	var source []byte
	if source, err = base64.RawStdEncoding.DecodeString(data); err != nil {
		return "", err
	}
	dst := make([]byte, len(source))
	encrypt.CryptBlocks(dst, source)
	return string(PKCS5UnPadding(dst)), nil
}

func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func PKCS5UnPadding(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
