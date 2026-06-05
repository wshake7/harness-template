package cn.harnesstemplate.admin.domain.model;

/**
 * RSA key pair stored in Redis under the global encrypt key.
 * Mirrors Go's domains.EncryptKeyPair.
 */
public record EncryptKeyPair(String publicKey, String privateKey) {}
