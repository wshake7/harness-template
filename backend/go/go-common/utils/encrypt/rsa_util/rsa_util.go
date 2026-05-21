package rsa_util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"go-common/utils/encrypt"
)

func EncryptRaw(data string, publicKey *rsa.PublicKey) (string, error) {
	hash := sha256.New()
	encryptedBytes, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, []byte(data), nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

func Encrypt(data string, publicKeyStr string) (string, error) {
	publicKey, err := ParsePublicKeyAuto(publicKeyStr)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	encryptedBytes, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, []byte(data), nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

func DecryptRaw(encryptedData string, privateKey *rsa.PrivateKey) (string, error) {
	hash := sha256.New()

	decodedData, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	decryptedBytes, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, decodedData, nil)
	if err != nil {
		return "", errors.New("RSA解密失败，密钥不匹配或数据被篡改")
	}

	return string(decryptedBytes), nil
}

func Decrypt(encryptedData string, privateKeyStr string) (string, error) {
	privateKey, err := ParsePrivateKeyAuto(privateKeyStr)
	if err != nil {
		return "", err
	}

	hash := sha256.New()

	decodedData, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	decryptedBytes, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, decodedData, nil)
	if err != nil {
		return "", errors.New("RSA解密失败，密钥不匹配或数据被篡改")
	}

	return string(decryptedBytes), nil
}

func GenerateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, encrypt.RSA_KEY_SIZE)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func GeneratePEMKeyPair() (privateKeyStr string, publicKeyStr string, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, encrypt.RSA_KEY_SIZE)
	if err != nil {
		return "", "", err
	}

	privateBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateBytes,
	}
	privateKeyStr = string(pem.EncodeToMemory(privateBlock))

	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}
	publicKeyStr = string(pem.EncodeToMemory(pubBlock))

	return
}

func GenerateKeyPair() (privateKeyStr string, publicKeyStr string, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, encrypt.RSA_KEY_SIZE)
	if err != nil {
		return "", "", err
	}

	// 私钥 PKCS1 DER → Base64
	privateBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyStr = base64.StdEncoding.EncodeToString(privateBytes)

	// 公钥 PKIX DER → Base64
	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	publicKeyStr = base64.StdEncoding.EncodeToString(pubBytes)

	return
}

func ParsePrivateKeyAuto(keyStr string) (*rsa.PrivateKey, error) {
	// 尝试 PEM
	if block, _ := pem.Decode([]byte(keyStr)); block != nil {
		if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
			return key, nil
		}
		if keyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
			if key, ok := keyInterface.(*rsa.PrivateKey); ok {
				return key, nil
			}
		}
		return nil, errors.New("解析 PEM 私钥失败")
	}

	// 尝试 Base64
	bytes, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return nil, errors.New("私钥格式错误")
	}

	if key, err := x509.ParsePKCS1PrivateKey(bytes); err == nil {
		return key, nil
	}

	keyInterface, err := x509.ParsePKCS8PrivateKey(bytes)
	if err != nil {
		return nil, errors.New("解析私钥失败")
	}

	key, ok := keyInterface.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("不是 RSA 私钥")
	}

	return key, nil
}

func ParsePublicKeyAuto(keyStr string) (*rsa.PublicKey, error) {
	// 尝试 PEM
	if block, _ := pem.Decode([]byte(keyStr)); block != nil {
		pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, errors.New("解析 PEM 公钥失败")
		}
		if pubKey, ok := pubInterface.(*rsa.PublicKey); ok {
			return pubKey, nil
		}
		return nil, errors.New("不是 RSA 公钥")
	}

	// 尝试 Base64
	bytes, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return nil, errors.New("公钥格式错误")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(bytes)
	if err != nil {
		return nil, errors.New("解析公钥失败")
	}

	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("不是 RSA 公钥")
	}

	return pubKey, nil
}
