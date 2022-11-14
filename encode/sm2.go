package encode

import (
	"encoding/hex"
	"errors"

	"github.com/tjfoc/gmsm/sm2"
)

type SM2Crypt struct {
	pwd []byte
}

func NewSM2Crypt(pwd string) *SM2Crypt {
	return &SM2Crypt{
		pwd: []byte(pwd),
	}
}

/*
 * 生产密钥
 */
func (s *SM2Crypt) GenerateKeyToMem() (string, string, error) {
	privateKey, err := sm2.GenerateKey()
	if err != nil {
		return "", "", err
	}
	priv, err := sm2.WritePrivateKeytoMem(privateKey, s.pwd)
	if err != nil {
		return "", "", err
	}
	pub, err := sm2.WritePublicKeytoMem(&privateKey.PublicKey, s.pwd)
	if err != nil {
		return "", "", err
	}
	return string(priv), string(pub), nil
}

/*
 * 生产密钥Pem
 */
func (s *SM2Crypt) GenerateKeyToPem(privKeyPath, pubKeyPath string) (bool, error) {
	privateKey, err := sm2.GenerateKey()
	if err != nil {
		return false, err
	}
	priv, err := sm2.WritePrivateKeytoPem(privKeyPath, privateKey, s.pwd)
	if err != nil {
		return false, err
	}
	pub, err := sm2.WritePublicKeytoPem(pubKeyPath, &privateKey.PublicKey, s.pwd)
	if err != nil {
		return false, err
	}
	if priv && pub {
		return true, nil
	} else {
		return false, errors.New("未生成密钥文件")
	}
}

/*
 * 读取私钥
 */
func (s *SM2Crypt) ReadPrivateKeyFromMem(key string) (*sm2.PrivateKey, error) {
	privateKey, err := sm2.ReadPrivateKeyFromMem([]byte(key), s.pwd)
	if err != nil {
		return nil, err
	} else {
		return privateKey, nil
	}
}

/*
 * 读取公钥
 */
func (s *SM2Crypt) ReadPublicKeyFromMem(key string) (*sm2.PublicKey, error) {
	publicKey, err := sm2.ReadPublicKeyFromMem([]byte(key), s.pwd)
	if err != nil {
		return nil, err
	} else {
		return publicKey, nil
	}
}

/*
 * 读取私钥
 */
func (s *SM2Crypt) ReadPrivateKeyFromPem(path string) (*sm2.PrivateKey, error) {
	privateKey, err := sm2.ReadPrivateKeyFromPem(path, s.pwd)
	if err != nil {
		return nil, err
	} else {
		return privateKey, nil
	}
}

/*
 * 读取公钥
 */
func (s *SM2Crypt) ReadPublicKeyFromPem(path string) (*sm2.PublicKey, error) {
	publicKey, err := sm2.ReadPublicKeyFromPem(path, s.pwd)
	if err != nil {
		return nil, err
	} else {
		return publicKey, nil
	}
}

/*
 * 加密
 */
func (s *SM2Crypt) Encrypt(data string, publicKey *sm2.PublicKey) (string, error) {
	encryptData, err := publicKey.Encrypt([]byte(data))
	if err != nil {
		return "", err
	} else {
		return hex.EncodeToString(encryptData), nil
	}
}

/*
 * 解密
 */
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
