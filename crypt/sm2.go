package crypt

import (
	"encoding/hex"

	"github.com/ivanlebron/gmsm/sm2"
)

type SM2Crypt struct {
	pwd []byte
}

func NewSM2Crypt(pwd string) *SM2Crypt {
	if pwd == "" {
		return &SM2Crypt{}
	}
	return &SM2Crypt{
		pwd: []byte(pwd),
	}
}

func (s *SM2Crypt) ReadPrivateKeyFromMem(key string) (*sm2.PrivateKey, error) {
	if len(s.pwd) == 0 {
		s.pwd = nil
	}
	privateKey, err := sm2.ReadPrivateKeyFromMem([]byte(key), s.pwd)
	if err != nil {
		return nil, err
	} else {
		return privateKey, nil
	}
}

func (s *SM2Crypt) ReadPublicKeyFromMem(key string) (*sm2.PublicKey, error) {
	if len(s.pwd) == 0 {
		s.pwd = nil
	}
	publicKey, err := sm2.ReadPublicKeyFromMem([]byte(key), s.pwd)
	if err != nil {
		return nil, err
	} else {
		return publicKey, nil
	}
}

func (s *SM2Crypt) Encrypt(data string, publicKey *sm2.PublicKey) (string, error) {
	encryptData, err := publicKey.Encrypt([]byte(data))
	if err != nil {
		return "", err
	} else {
		return hex.EncodeToString(encryptData), nil
	}
}

func (s *SM2Crypt) Decrypt(data string, privateKey *sm2.PrivateKey) (string, error) {
	tmp, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}
	decryptData, err := privateKey.Decrypt(tmp)
	if err != nil {
		return "", err
	} else {
		return string(decryptData), nil
	}
}
