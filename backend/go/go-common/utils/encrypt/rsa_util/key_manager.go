package rsa_util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"go-common/utils/encrypt"
)

type KeyPairManager struct {
	PublicKey       *rsa.PublicKey
	PrivateKey      *rsa.PrivateKey
	PublicKeyBase64 string
}

// NewKeyPairManager 生成 RSA 密钥对
func NewKeyPairManager() (*KeyPairManager, error) {

	privateKey, err := rsa.GenerateKey(rand.Reader, encrypt.RSA_KEY_SIZE)
	if err != nil {
		return nil, err
	}

	publicKey := &privateKey.PublicKey

	// Go: X.509 PKIX（等价 Java publicKey.getEncoded()）
	pubBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	return &KeyPairManager{
		PublicKey:       publicKey,
		PrivateKey:      privateKey,
		PublicKeyBase64: base64.StdEncoding.EncodeToString(pubBytes),
	}, nil
}

func (k *KeyPairManager) PrivateKeyPEM() string {
	der := x509.MarshalPKCS1PrivateKey(k.PrivateKey)
	return base64.StdEncoding.EncodeToString(der)
}
