package secret

import (
	"errors"
	"os"
)

var (
	sysPrivate = `
-----BEGIN PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqBHM9VAYItBHkwdwIBAQQgTb8gsPZI/r6skcKg
SoHDwfCfLkU2XysZCqHdWbf97kKgCgYIKoEcz1UBgi2hRANCAARqcTpGGBKLbmYe
E+wOQFbrV5gGbNd8G+iMeKsws1pOUw80V1CbUFruF0e/MG+weveRcQsv+NUqT6X7
1HK0HTAo
-----END PRIVATE KEY-----
`
	sysPublic = `
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoEcz1UBgi0DQgAEanE6RhgSi25mHhPsDkBW61eYBmzX
fBvojHirMLNaTlMPNFdQm1Ba7hdHvzBvsHr3kXELL/jVKk+l+9RytB0wKA==
-----END PUBLIC KEY-----
`
)

func NewCrypt(private, public string) (*Crypt, error) {
	sc := new(Crypt)
	if private != "" {
		sc.mode = "sm2"
		sm2, err := NewSM2(private, public)
		if err != nil {
			return sc, err
		}
		sc.sm2 = sm2
	} else {
		switch os.Getenv("CRYPTMODE") {
		case "sm2":
			sc.mode = "sm2"
			sm2, err := NewSM2(sysPrivate, sysPublic)
			if err != nil {
				return sc, err
			}
			sc.sm2 = sm2
		default:
			sc.mode = "aes"
			aes := NewAES()
			sc.aes = aes
		}
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
