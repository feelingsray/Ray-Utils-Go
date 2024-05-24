package crypt

import (
	"bytes"
	"crypto/cipher"
	"encoding/hex"

	"github.com/ivanlebron/gmsm/sm4"
)

type SM4Crypt struct {
	key []byte
}

func NewSM4Crypt(key string) *SM4Crypt {
	return &SM4Crypt{
		key: []byte(key),
	}
}

func (s *SM4Crypt) Encrypt(data string) (string, error) {
	block, err := sm4.NewCipher(s.key)
	if err != nil {
		return "", err
	}
	src := []byte(data)
	a := block.BlockSize() - len(src)%block.BlockSize()
	repeat := bytes.Repeat([]byte{byte(a)}, a)
	newsrc := append(src, repeat...)
	dst := make([]byte, len(newsrc))
	blockMode := cipher.NewCBCEncrypter(block, s.key[:block.BlockSize()])
	blockMode.CryptBlocks(dst, newsrc)
	return hex.EncodeToString(dst), nil
}

func (s *SM4Crypt) Decrypt(data string) (string, error) {
	block, err := sm4.NewCipher(s.key)
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, s.key[:block.BlockSize()])
	dst, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}
	src := make([]byte, len(dst))
	blockMode.CryptBlocks(src, dst)
	num := int(src[len(src)-1])
	newsrc := src[:len(src)-num]
	return string(newsrc), nil
}
