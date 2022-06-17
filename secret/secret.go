package secret

import (
	"github.com/feelingsray/Ray-Utils-Go/encode"
	"github.com/tjfoc/gmsm/sm2"
)

func NewSecretCrypt() (*SecretCrypt,error) {
	sc := new(SecretCrypt)
	sc.privatK = `
-----BEGIN ENCRYPTED PRIVATE KEY-----
MIH8MFcGCSqGSIb3DQEFDTBKMCkGCSqGSIb3DQEFDDAcBAiDaVfqFsXo3wICCAAw
DAYIKoZIhvcNAgcFADAdBglghkgBZQMEASoEEJ4c0sQeuMSsPuoUqLuQkkEEgaBm
kM09xMwjTXNxTgHdxKP+9y5jknqinmAYfQoqMBq3E4VGVK07+2NbKKfnjpdSLsTV
tPol9qY1XkI8EI4S4EqzExpawC4NG1uhXs7LxhcFXBtOk0bPadLH5/BLyRHc06tG
oGephca/lxlPzU4LsjB5kNeuo2fhflxCSEByxYXOqkLYRERip5MqM0/NxN09zFTe
K2SFtkCC1Tb2I8nkTU9w
-----END ENCRYPTED PRIVATE KEY-----
`
	sc.pubK = `
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoEcz1UBgi0DQgAESp2XedoldFaeedu8b6OMZDyMQmJR
Kuto6F8KIcLd58OoYfsfaC78E+T8kI2CS9U8jLR9X6fzGDSEGpYxEzN4Ng==
-----END PUBLIC KEY-----
`
	sc.pwd = "r19a224y"
	crypt := encode.NewSM2Crypt(sc.pwd)
	privateKey, err := crypt.ReadPrivateKeyFromMem(sc.privatK)
	if err != nil {
		return nil, err
	}
	publicKey, err := crypt.ReadPublicKeyFromMem(sc.pubK)
	if err != nil {
		return nil, err
	}
	sc.Crypt = crypt
	sc.PrivateKey = privateKey
	sc.PublicKey = publicKey
	return sc,nil
}

type SecretCrypt struct {
	Crypt                   *encode.SM2Crypt    // SM2加密
	PrivateKey              *sm2.PrivateKey     // 私钥
	PublicKey               *sm2.PublicKey      // 公钥
	pwd                     string              // 密码
	privatK                 string              // 私钥
	pubK                    string              // 公钥
}

func (s *SecretCrypt)Decrypt(raw string) (string,error) {
	data,err := s.Crypt.Decrypt(raw,s.PrivateKey)
	if err != nil {
		return "",err
	} else {
		return data,nil
	}
}

func (s *SecretCrypt)Encrypt(raw string) (string,error) {
	data,err := s.Crypt.Encrypt(raw,s.PublicKey)
	if err != nil {
		return "",err
	} else {
		return data,nil
	}
}
