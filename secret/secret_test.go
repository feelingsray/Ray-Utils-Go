package secret

import (
	"log"
	"testing"
	
	"github.com/ivanlebron/gmsm/sm2"
)

func TestNewSecretCrypt(t *testing.T) {
	
	crypt, err := NewSecretCrypt()
	if err != nil {
		t.Fatal(err)
	}
	encrypt, err := crypt.Encrypt("hello")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(encrypt)
	
	decrypt, err := crypt.Decrypt(encrypt)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(decrypt)
}
func TestGeneralFile(t *testing.T) {
	
	priv, err := sm2.GenerateKey() // 生成密钥对
	if err != nil {
		log.Fatal(err)
	}
	privByte, err := sm2.WritePrivateKeytoMem(priv, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(privByte))
	
	pubKey, _ := priv.Public().(*sm2.PublicKey)
	
	pubByte, err := sm2.WritePublicKeytoMem(pubKey, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(pubByte))
}

//-----BEGIN PRIVATE KEY-----
//MIGTAgEAMBMGByqGSM49AgEGCCqBHM9VAYItBHkwdwIBAQQgfUKyp1fzOhixB2OP
//kmyynzBbM1byYMJJQL76N/BVg+CgCgYIKoEcz1UBgi2hRANCAAQeounVibe0P2iE
//0wUmbUmjy+uMMavtrm3cOKO+SDgli+PXUuEpakwvsKv5VcRlRPPHm7GtLuGNDmXH
//K5r3lIZz
//-----END PRIVATE KEY-----
//
//-----BEGIN PUBLIC KEY-----
//MFkwEwYHKoZIzj0CAQYIKoEcz1UBgi0DQgAEHqLp1Ym3tD9ohNMFJm1Jo8vrjDGr
//7a5t3Dijvkg4JYvj11LhKWpML7Cr+VXEZUTzx5uxrS7hjQ5lxyua95SGcw==
//-----END PUBLIC KEY-----
