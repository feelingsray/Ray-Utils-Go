package encode

import (
    "crypto/aes"
    "crypto/cipher"
    "encoding/base64"
)

type AESCryptor struct {
    Key []byte
    IV  []byte
}

func NewAESCryptor(key, iv string) *AESCryptor {
    k, _ := base64.StdEncoding.DecodeString(key)
    v, _ := base64.StdEncoding.DecodeString(iv)
    return &AESCryptor{
        Key: k,
        IV:  v,
    }
}

func (a *AESCryptor) AESBase64Encrypt(data string) (string, error) {
    iv := a.IV
    key := a.Key
    var block cipher.Block
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    encrypt := cipher.NewCBCEncrypter(block, iv)
    source := PKCS5Padding([]byte(data), block.BlockSize())
    dst := make([]byte, len(source))
    encrypt.CryptBlocks(dst, source)
    return base64.RawStdEncoding.EncodeToString(dst), nil
}

func (a *AESCryptor) AESBase64Decrypt(data string) (string, error) {
    iv := a.IV
    key := a.Key
    var block cipher.Block
    block, err := aes.NewCipher([]byte(key))
    if err != nil {
        return "", err
    }
    encrypt := cipher.NewCBCDecrypter(block, iv)
    var source []byte
    if source, err = base64.RawStdEncoding.DecodeString(data); err != nil {
        return "", err
    }
    dst := make([]byte, len(source))
    encrypt.CryptBlocks(dst, source)
    return string(PKCS5UnPadding(dst)), nil
}
