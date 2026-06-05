package cn.harnesstemplate.admin.web.filter;

import cn.harnesstemplate.admin.infrastructure.auth.EncryptKeyPairService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;

/**
 * Holds the global server RSA key pair for encrypting public routes (login).
 * Uses Redis cache-aside: tries Redis first, generates and caches on miss.
 * This matches Go's EncryptHandler.PublicKey pattern.
 *
 * Results are cached in volatile fields so the EncryptFilter can call
 * {@link #getPrivateKeyPem()} on every encrypted request without hitting Redis.
 */
@Component
public class ServerKeyPairProvider {

    private static final Logger log = LoggerFactory.getLogger(ServerKeyPairProvider.class);

    private final EncryptKeyPairService service;

    private volatile String privateKeyPem;
    private volatile String publicKey;

    public ServerKeyPairProvider(EncryptKeyPairService service) {
        this.service = service;
    }

    public String getPrivateKeyPem() {
        if (privateKeyPem == null) {
            synchronized (this) {
                if (privateKeyPem == null) {
                    loadKeyPair();
                }
            }
        }
        return privateKeyPem;
    }

    public String getPublicKey() {
        if (publicKey == null) {
            synchronized (this) {
                if (publicKey == null) {
                    loadKeyPair();
                }
            }
        }
        return publicKey;
    }

    private void loadKeyPair() {
        var pair = service.getEncryptKeyPair();
        if (pair != null) {
            publicKey = pair.publicKey();
            privateKeyPem = pair.privateKey();
            log.debug("Loaded encrypt key pair from Redis");
            return;
        }

        log.info("Encrypt key pair not found in Redis, generating new pair");
        var result = service.generateAndCacheKeyPair();
        publicKey = result.publicKey();
        privateKeyPem = result.privateKeyPem();
    }
}
