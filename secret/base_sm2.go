package secret

import (
	"github.com/ivanlebron/gmsm/sm2"

	"github.com/feelingsray/ray-utils-go/v2/crypt"
)

func NewSM2(private, public string) (*SM2Crypt, error) {
	sc := new(SM2Crypt)
	sc.privateK = private
	sc.publicK = public
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

type SM2Crypt struct {
	Crypt      *crypt.SM2Crypt // SM2加密
	PrivateKey *sm2.PrivateKey // 私钥
	PublicKey  *sm2.PublicKey  // 公钥
	pwd        string          // 密码
	privateK   string          // 私钥
	publicK    string          // 公钥
}

func (s *SM2Crypt) decrypt(raw string) (string, error) {
	data, err := s.Crypt.Decrypt(raw, s.PrivateKey)
	if err != nil {
		return "", err
	} else {
		return data, nil
	}
}

func (s *SM2Crypt) encrypt(raw string) (string, error) {
	data, err := s.Crypt.Encrypt(raw, s.PublicKey)
	if err != nil {
		return "", err
	} else {
		return data, nil
	}
}
