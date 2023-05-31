package secret

import (
	"github.com/feelingsray/ray-utils-go/v2/serialize"
	"github.com/spf13/cast"
	"log"
	"testing"
	"time"

	"github.com/ivanlebron/gmsm/sm2"
)

func TestNewSecretCrypt(t *testing.T) {
	crypt, err := NewCrypt()
	if err != nil {
		t.Fatal(err)
	}
	test := map[string]int{}
	for i := 0; i < 10000; i++ {
		test[cast.ToString(i)] = i
	}
	json, _ := serialize.DumpJson(test)
	start := time.Now()
	for x := 0; x < 1000; x++ {
		encrypt, err := crypt.Encrypt(string(json))
		if err != nil {
			t.Fatal(err)
		}
		//t.Log(encrypt)

		decrypt, err := crypt.Decrypt(encrypt)
		if err != nil {
			t.Fatal(err)
		}
		_, _ = encrypt, decrypt
	}
	//t.Log(decrypt)
	time.Since(start)
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
//

//=== RUN   TestNewSecretCrypt
//--- PASS: TestNewSecretCrypt (7.97s)
//PASS
//=== RUN   TestNewSecretCrypt
//--- PASS: TestNewSecretCrypt (187.6s)
