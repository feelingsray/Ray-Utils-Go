package secret

import (
	"github.com/ivanlebron/gmsm/sm2"

	"github.com/feelingsray/ray-utils-go/v2/crypt"
)

func NewSecretCrypt() (*Sm2Crypt, error) {
	sc := new(Sm2Crypt)
	sc.privateK = `
-----BEGIN PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqBHM9VAYItBHkwdwIBAQQgTb8gsPZI/r6skcKg
SoHDwfCfLkU2XysZCqHdWbf97kKgCgYIKoEcz1UBgi2hRANCAARqcTpGGBKLbmYe
E+wOQFbrV5gGbNd8G+iMeKsws1pOUw80V1CbUFruF0e/MG+weveRcQsv+NUqT6X7
1HK0HTAo
-----END PRIVATE KEY-----
`
	sc.publicK = `
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoEcz1UBgi0DQgAEanE6RhgSi25mHhPsDkBW61eYBmzX
fBvojHirMLNaTlMPNFdQm1Ba7hdHvzBvsHr3kXELL/jVKk+l+9RytB0wKA==
-----END PUBLIC KEY-----
`
	//sc.pwd = "r19a224y"
	sm2Crypt := crypt.NewSM2Crypt(sc.pwd)
	privateKey, err := sm2Crypt.ReadPrivateKeyFromMem(sc.privateK)
	if err != nil {
		return nil, err
	}
	publicKey, err := sm2Crypt.ReadPublicKeyFromMem(sc.publicK)
	if err != nil {
		return nil, err
	}
	sc.Crypt = sm2Crypt
	sc.PrivateKey = privateKey
	sc.PublicKey = publicKey
	return sc, nil
}

type Sm2Crypt struct {
	Crypt      *crypt.SM2Crypt // SM2加密
	PrivateKey *sm2.PrivateKey // 私钥
	PublicKey  *sm2.PublicKey  // 公钥
	pwd        string          // 密码
	privateK   string          // 私钥
	publicK    string          // 公钥
}

func (s *Sm2Crypt) Decrypt(raw string) (string, error) {
	data, err := s.Crypt.Decrypt(raw, s.PrivateKey)
	if err != nil {
		return "", err
	} else {
		return data, nil
	}
}

func (s *Sm2Crypt) Encrypt(raw string) (string, error) {
	data, err := s.Crypt.Encrypt(raw, s.PublicKey)
	if err != nil {
		return "", err
	} else {
		return data, nil
	}
}
