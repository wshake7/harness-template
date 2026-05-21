package aes_util

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"
)

// ─── GenerateAESKey ───────────────────────────────────────────────────────────

func TestGenerateAESKey_Length(t *testing.T) {
	key, err := GenerateAESKey()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key == "" {
		t.Fatal("key should not be empty")
	}
}

func TestGenerateAESKey_Unique(t *testing.T) {
	k1, _ := GenerateAESKey()
	k2, _ := GenerateAESKey()
	if k1 == k2 {
		t.Fatal("two generated keys should not be identical")
	}
}

// ─── Encrypt ─────────────────────────────────────────────────────────────────

func TestEncrypt_Success(t *testing.T) {
	key, _ := GenerateAESKey()
	result, err := Encrypt("hello world", key, "aad-data")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Ciphertext == "" || result.TagIv == "" || result.Combined == "" {
		t.Fatal("result fields should not be empty")
	}
}

func TestEncrypt_DifferentIvEachTime(t *testing.T) {
	key, _ := GenerateAESKey()
	r1, _ := Encrypt("same plaintext", key, "aad")
	r2, _ := Encrypt("same plaintext", key, "aad")
	if r1.Combined == r2.Combined {
		t.Fatal("two encryptions of the same data should produce different Combined (random IV)")
	}
}

func TestEncrypt_InvalidKeyBase64(t *testing.T) {
	_, err := Encrypt("data", "not-valid-base64!!!", "aad")
	if err == nil {
		t.Fatal("expected error for invalid base64 key")
	}
}

func TestEncrypt_InvalidKeyLength(t *testing.T) {
	// 只有 8 字节，不足 256-bit
	import_b64 := "c2hvcnQ=" // base64("short")
	_, err := Encrypt("data", import_b64, "aad")
	if err == nil {
		t.Fatal("expected error for wrong key length")
	}
}

func TestEncrypt_EmptyPlaintext(t *testing.T) {
	key, _ := GenerateAESKey()
	result, err := Encrypt("", key, "aad")
	if err != nil {
		t.Fatalf("unexpected error encrypting empty string: %v", err)
	}
	if result.Combined == "" {
		t.Fatal("Combined should not be empty even for empty plaintext")
	}
}

func TestEncrypt_EmptyAAD(t *testing.T) {
	key, _ := GenerateAESKey()
	_, err := Encrypt("hello", key, "")
	if err != nil {
		t.Fatalf("unexpected error with empty AAD: %v", err)
	}
}

// ─── Decrypt ─────────────────────────────────────────────────────────────────

func TestDecrypt_RoundTrip(t *testing.T) {
	key, _ := GenerateAESKey()
	plain := "the quick brown fox"
	aad := "some-context"

	result, err := Encrypt(plain, key, aad)
	if err != nil {
		t.Fatalf("encrypt error: %v", err)
	}

	got, err := Decrypt(result.Combined, key, aad)
	if err != nil {
		t.Fatalf("decrypt error: %v", err)
	}
	if string(got) != plain {
		t.Fatalf("expected %q, got %q", plain, string(got))
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	key, _ := GenerateAESKey()
	wrongKey, _ := GenerateAESKey()

	result, _ := Encrypt("secret", key, "aad")
	_, err := Decrypt(result.Combined, wrongKey, "aad")
	if err == nil {
		t.Fatal("expected error when decrypting with wrong key")
	}
}

func TestDecrypt_WrongAAD(t *testing.T) {
	key, _ := GenerateAESKey()
	result, _ := Encrypt("secret", key, "correct-aad")
	_, err := Decrypt(result.Combined, key, "wrong-aad")
	if err == nil {
		t.Fatal("expected error when decrypting with wrong AAD")
	}
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	key, _ := GenerateAESKey()
	result, _ := Encrypt("secret", key, "aad")

	// 翻转第一个字节
	raw := result.CombinedRaw
	raw[0] ^= 0xFF
	tampered := base64.StdEncoding.EncodeToString(raw)

	_, err := Decrypt(tampered, key, "aad")
	if err == nil {
		t.Fatal("expected error for tampered ciphertext")
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	key, _ := GenerateAESKey()
	_, err := Decrypt("!!!not-base64", key, "aad")
	if err == nil {
		t.Fatal("expected error for invalid base64 combined")
	}
}

func TestDecrypt_TooShortData(t *testing.T) {
	key, _ := GenerateAESKey()
	// 只有 3 字节，远小于 IV+Tag 最小长度
	short := base64.StdEncoding.EncodeToString([]byte{0x01, 0x02, 0x03})
	_, err := Decrypt(short, key, "aad")
	if err == nil {
		t.Fatal("expected error for too-short data")
	}
}

func TestDecrypt_EmptyPlaintext_RoundTrip(t *testing.T) {
	key, _ := GenerateAESKey()
	result, _ := Encrypt("", key, "aad")
	got, err := Decrypt(result.Combined, key, "aad")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty plaintext, got %q", got)
	}
}

// ─── DecryptCiphertextAndTag ──────────────────────────────────────────────────

func TestDecryptCiphertextAndTag_RoundTrip(t *testing.T) {
	key, _ := GenerateAESKey()
	plain := "hello from tag+iv path"
	aad := "ctx"

	result, err := Encrypt(plain, key, aad)
	if err != nil {
		t.Fatalf("encrypt error: %v", err)
	}

	got, err := DecryptCiphertextAndTag(result.Ciphertext, result.TagIv, key, aad)
	if err != nil {
		t.Fatalf("decrypt error: %v", err)
	}
	if string(got) != plain {
		t.Fatalf("expected %q, got %q", plain, string(got))
	}
}

func TestDecryptCiphertextAndTag_WrongKey(t *testing.T) {
	key, _ := GenerateAESKey()
	wrongKey, _ := GenerateAESKey()
	result, _ := Encrypt("secret", key, "aad")
	_, err := DecryptCiphertextAndTag(result.Ciphertext, result.TagIv, wrongKey, "aad")
	if err == nil {
		t.Fatal("expected error for wrong key")
	}
}

func TestDecryptCiphertextAndTag_WrongAAD(t *testing.T) {
	key, _ := GenerateAESKey()
	result, _ := Encrypt("secret", key, "correct")
	_, err := DecryptCiphertextAndTag(result.Ciphertext, result.TagIv, key, "wrong")
	if err == nil {
		t.Fatal("expected error for wrong AAD")
	}
}

func TestDecryptCiphertextAndTag_InvalidTagIvLength(t *testing.T) {
	key, _ := GenerateAESKey()
	result, _ := Encrypt("data", key, "aad")

	// 截断 tagIv，使其长度不对
	raw := result.TagIvRaw[:5]
	badTagIv := base64.StdEncoding.EncodeToString(raw)

	_, err := DecryptCiphertextAndTag(result.Ciphertext, badTagIv, key, "aad")
	if err == nil {
		t.Fatal("expected error for invalid tagIv length")
	}
}

func TestDecryptCiphertextAndTag_InvalidCiphertextBase64(t *testing.T) {
	key, _ := GenerateAESKey()
	result, _ := Encrypt("data", key, "aad")
	_, err := DecryptCiphertextAndTag("!!!bad-base64", result.TagIv, key, "aad")
	if err == nil {
		t.Fatal("expected error for invalid ciphertext base64")
	}
}

func TestDecryptCiphertextAndTag_InvalidTagIvBase64(t *testing.T) {
	key, _ := GenerateAESKey()
	result, _ := Encrypt("data", key, "aad")
	_, err := DecryptCiphertextAndTag(result.Ciphertext, "!!!bad-base64", key, "aad")
	if err == nil {
		t.Fatal("expected error for invalid tagIv base64")
	}
}

// ─── 两种解密路径结果一致性 ───────────────────────────────────────────────────

func TestDecrypt_BothPathsAgree(t *testing.T) {
	key, _ := GenerateAESKey()
	plain := ""
	aad := "aad-val"

	result, _ := Encrypt(plain, key, aad)

	via1, err1 := Decrypt(result.Combined, key, aad)
	via2, err2 := DecryptCiphertextAndTag(result.Ciphertext, result.TagIv, key, aad)

	if err1 != nil || err2 != nil {
		t.Fatalf("errors: Decrypt=%v, DecryptCiphertextAndTag=%v", err1, err2)
	}
	if !bytes.Equal(via1, via2) {
		t.Fatalf("two decrypt paths returned different results: %q vs %q", via1, via2)
	}
}

// ─── 长文本 & Unicode ─────────────────────────────────────────────────────────

func TestEncryptDecrypt_LongText(t *testing.T) {
	key, _ := GenerateAESKey()
	plain := strings.Repeat("数据加密测试 ", 500)
	result, err := Encrypt(plain, key, "aad")
	if err != nil {
		t.Fatalf("encrypt error: %v", err)
	}
	got, err := Decrypt(result.Combined, key, "aad")
	if err != nil {
		t.Fatalf("decrypt error: %v", err)
	}
	if string(got) != plain {
		t.Fatal("long text round-trip mismatch")
	}
}

func TestEncryptDecrypt_UnicodeAAD(t *testing.T) {
	key, _ := GenerateAESKey()
	_, err := Encrypt("data", key, "用户ID=42&场景=登录")
	if err != nil {
		t.Fatalf("unexpected error with unicode AAD: %v", err)
	}
}
