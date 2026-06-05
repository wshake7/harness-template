package cn.harnesstemplate.admin.infrastructure.auth;

import cn.harnesstemplate.admin.domain.model.EncryptKeyPair;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.stereotype.Service;

import java.security.KeyPair;

/**
 * Cache-aside service for the global RSA encrypt key pair, backed by Redis.
 * Mirrors Go's appsvc.RedisCache + GenerateAndCacheKeyPair.
 */
@Service
public class EncryptKeyPairService {

    private static final Logger log = LoggerFactory.getLogger(EncryptKeyPairService.class);
    private static final String REDIS_KEY = "global:encrypt:public:key";
    private static final ObjectMapper objectMapper = new ObjectMapper();

    private final StringRedisTemplate redis;

    public EncryptKeyPairService(StringRedisTemplate redis) {
        this.redis = redis;
    }

    /**
     * Get the key pair from Redis. Returns null on cache miss.
     */
    public EncryptKeyPair getEncryptKeyPair() {
        String json = redis.opsForValue().get(REDIS_KEY);
        if (json == null || json.isEmpty()) {
            return null;
        }
        try {
            return objectMapper.readValue(json, EncryptKeyPair.class);
        } catch (Exception e) {
            log.error("Failed to deserialize encrypt key pair from Redis", e);
            return null;
        }
    }

    /**
     * Store the key pair in Redis.
     */
    public void setEncryptKeyPair(String publicKey, String privateKey) {
        try {
            EncryptKeyPair pair = new EncryptKeyPair(publicKey, privateKey);
            String json = objectMapper.writeValueAsString(pair);
            redis.opsForValue().set(REDIS_KEY, json);
            log.info("Global encrypt key pair stored in Redis");
        } catch (Exception e) {
            log.error("Failed to serialize encrypt key pair to Redis", e);
        }
    }

    /**
     * Generate a new RSA key pair and cache it in Redis.
     * Returns the base64-encoded public key and PEM-encoded private key.
     */
    public KeyPairResult generateAndCacheKeyPair() {
        KeyPair keyPair = CryptoService.generateRsaKeyPair();
        String publicKey = CryptoService.toBase64(keyPair.getPublic());
        String privateKeyPem = CryptoService.toPem(keyPair.getPrivate());
        setEncryptKeyPair(publicKey, privateKeyPem);
        return new KeyPairResult(publicKey, privateKeyPem);
    }

    public record KeyPairResult(String publicKey, String privateKeyPem) {}
}
