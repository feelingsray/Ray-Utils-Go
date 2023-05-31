package secret

import (
	"errors"
	"os"
)

func NewCrypt() (*Crypt, error) {
	sc := new(Crypt)
	switch os.Getenv("CRYPTMODE") {
	case "sm2":
		sc.mode = "sm2"
		sm2, err := NewSM2()
		if err != nil {
			return sc, err
		}
		sc.sm2 = sm2
	default:
		sc.mode = "aes"
		aes := NewAES()
		sc.aes = aes
	}
	return sc, nil
}

type Crypt struct {
	mode string
	sm2  *SM2Crypt
	aes  *AESCrypt
}

func (s *Crypt) Decrypt(raw string) (string, error) {
	switch s.mode {
	case "aes":
		return s.aes.decrypt(raw)
	case "sm2":
		return s.sm2.decrypt(raw)
	}
	return "", errors.New("crypt mode error")
}

func (s *Crypt) Encrypt(raw string) (string, error) {
	switch s.mode {
	case "aes":
		return s.aes.encrypt(raw)
	case "sm2":
		return s.sm2.encrypt(raw)
	}
	return "", nil
}
