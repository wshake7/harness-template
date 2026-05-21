package aes_util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"go-common/utils/encrypt"
	"io"
)

// GenerateAESKey 生成 Base64 编码的 256-bit key
func GenerateAESKey() (string, error) {
	key := make([]byte, encrypt.AES_KEY_SIZE)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// EncryptResult 加密结果
type EncryptResult struct {
	CiphertextRaw []byte // 密文（GET时为空）
	Ciphertext    string // base64(密文)
	TagIvRaw      []byte // 认证标签
	TagIv         string
	CombinedRaw   []byte
	Combined      string // base64(IV + ciphertext + tag)
}

// Encrypt AES-GCM 加密（返回 Base64(iv + ciphertext + tag)）
func Encrypt(data string, base64Key string, aad string) (*EncryptResult, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, err
	}
	if len(keyBytes) != encrypt.AES_KEY_SIZE {
		return nil, errors.New("AES密钥长度必须为256位")
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCMWithTagSize(block, encrypt.GCM_TAG_LENGTH)
	if err != nil {
		return nil, err
	}

	// 生成 IV
	iv := make([]byte, encrypt.GCM_IV_LENGTH)
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// Seal: 输出 ciphertext + tag (最后16字节)
	sealed := gcm.Seal(nil, iv, []byte(data), []byte(aad))

	// 分离密文和标签
	if len(sealed) < encrypt.GCM_TAG_LENGTH {
		return nil, errors.New("sealed data too short")
	}

	ciphertext := sealed[:len(sealed)-encrypt.GCM_TAG_LENGTH]
	tag := sealed[len(sealed)-encrypt.GCM_TAG_LENGTH:]

	tagIv := append(tag, iv...)
	combined := append(sealed, iv...)
	return &EncryptResult{
		CiphertextRaw: ciphertext,
		Ciphertext:    base64.StdEncoding.EncodeToString(ciphertext),
		TagIvRaw:      tagIv,
		TagIv:         base64.StdEncoding.EncodeToString(tagIv),
		CombinedRaw:   combined,
		Combined:      base64.StdEncoding.EncodeToString(combined),
	}, nil
}

// Decrypt AES-GCM 解密
func Decrypt(combined string, key string, aad string) ([]byte, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	if len(keyBytes) != encrypt.AES_KEY_SIZE {
		return nil, errors.New("AES密钥长度必须为256位")
	}

	data, err := base64.StdEncoding.DecodeString(combined)
	if err != nil {
		return nil, err
	}
	if len(data) < encrypt.GCM_IV_LENGTH+encrypt.GCM_TAG_LENGTH {
		return nil, errors.New("密文格式错误")
	}

	// combined = sealed(ciphertext + tag) + iv
	sealed := data[:len(data)-encrypt.GCM_IV_LENGTH]
	iv := data[len(data)-encrypt.GCM_IV_LENGTH:]

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCMWithTagSize(block, encrypt.GCM_TAG_LENGTH)
	if err != nil {
		return nil, err
	}

	plain, err := gcm.Open(nil, iv, sealed, []byte(aad))
	if err != nil {
		return nil, errors.New("AES解密失败，密钥不匹配、数据被篡改或AAD不一致")
	}

	return plain, nil
}

func DecryptCiphertextAndTag(ciphertext string, tagIv string, key string, aad string) ([]byte, error) {
	// 解码 key
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("key decode error: %w", err)
	}
	if len(keyBytes) != encrypt.AES_KEY_SIZE {
		return nil, errors.New("AES密钥长度必须为256位")
	}

	// 解码 ciphertext
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("ciphertext decode error: %w", err)
	}

	// 解码 tagIv（tag + iv 拼接）
	tagIvBytes, err := base64.StdEncoding.DecodeString(tagIv)
	if err != nil {
		return nil, fmt.Errorf("tagIv decode error: %w", err)
	}
	if len(tagIvBytes) != encrypt.GCM_TAG_LENGTH+encrypt.GCM_IV_LENGTH {
		return nil, errors.New("tagIv length invalid")
	}

	// 分离 tag 和 iv
	tag := tagIvBytes[:encrypt.GCM_TAG_LENGTH]
	iv := tagIvBytes[encrypt.GCM_TAG_LENGTH:]

	// 重新拼接为 GCM 需要的格式：ciphertext + tag
	sealed := append(ciphertextBytes, tag...)

	// 初始化 AES-GCM
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCMWithTagSize(block, encrypt.GCM_TAG_LENGTH)
	if err != nil {
		return nil, err
	}

	// 解密（同时验证 tag）
	plaintext, err := gcm.Open(nil, iv, sealed, []byte(aad))
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plaintext, nil
}
