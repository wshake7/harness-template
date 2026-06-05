package cn.harnesstemplate.admin.infrastructure.auth;

import javax.crypto.Cipher;
import javax.crypto.spec.OAEPParameterSpec;
import javax.crypto.spec.PSource;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

import java.security.KeyPair;
import java.security.spec.MGF1ParameterSpec;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.util.LinkedHashMap;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

class CryptoServiceTests {

    private static CryptoService service;
    private static String aesKeyBase64;
    private static PrivateKey privateKey;
    private static PublicKey publicKey;
    private static String privateKeyPem;
    private static String publicKeyPem;

    // Fixed golden parameters from frontend-api-samples.md
    private static final long TIMESTAMP = 1711411200000L;
    private static final String REQUEST_ID = "test-request-id-0001";
    private static final String PLAIN_BODY = "{\"username\":\"admin\"}";

    @BeforeAll
    static void setUp() {
        service = new CryptoService();
        KeyPair rsaKeyPair = CryptoService.generateRsaKeyPair();
        privateKey = rsaKeyPair.getPrivate();
        publicKey = rsaKeyPair.getPublic();
        privateKeyPem = CryptoService.toPem(privateKey);
        publicKeyPem = CryptoService.toPem(publicKey);
        aesKeyBase64 = CryptoService.generateAesKey();
    }

    // ---- RSA ----

    @Test
    void rsaShouldDecryptAesKeyEncryptedWithPublicKey() {
        String encryptedAesKey = service.rsaEncrypt(aesKeyBase64, publicKey);
        String decrypted = service.rsaDecrypt(encryptedAesKey, privateKey);
        assertEquals(aesKeyBase64, decrypted, "RSA round-trip should preserve AES key");
    }

    @Test
    void rsaShouldDecryptViaPemParsing() {
        PrivateKey parsedKey = CryptoService.parsePrivateKeyPem(privateKeyPem);
        PublicKey parsedPub = CryptoService.parsePublicKeyPem(publicKeyPem);

        String encrypted = service.rsaEncrypt(aesKeyBase64, parsedPub);
        String decrypted = service.rsaDecrypt(encrypted, parsedKey);
        assertEquals(aesKeyBase64, decrypted);
    }

    @Test
    void rsaShouldDecryptPayloadEncryptedWithSha256Mgf1Sha256() throws Exception {
        Cipher cipher = Cipher.getInstance("RSA/ECB/OAEPWithSHA-256AndMGF1Padding");
        cipher.init(
                Cipher.ENCRYPT_MODE,
                publicKey,
                new OAEPParameterSpec(
                        "SHA-256",
                        "MGF1",
                        MGF1ParameterSpec.SHA256,
                        PSource.PSpecified.DEFAULT
                )
        );

        String encrypted = java.util.Base64.getEncoder()
                .encodeToString(cipher.doFinal(aesKeyBase64.getBytes()));

        String decrypted = service.rsaDecrypt(encrypted, privateKey);
        assertEquals(aesKeyBase64, decrypted);
    }

    @Test
    void rsaDecryptWithWrongKeyShouldFail() {
        KeyPair otherKey = CryptoService.generateRsaKeyPair();
        String encryptedWithOther = service.rsaEncrypt(aesKeyBase64, otherKey.getPublic());
        assertThrows(CryptoService.CryptoException.class, () ->
                service.rsaDecrypt(encryptedWithOther, privateKey));
    }

    // ---- AES-GCM Encrypt/Decrypt (Combined format) ----

    @Test
    void aesShouldEncryptAndDecryptCombinedWithAad() {
        String aad = buildStandardAad();

        CryptoService.EncryptResult result = service.aesEncrypt(PLAIN_BODY, aesKeyBase64, aad);
        assertNotNull(result.ciphertext);
        assertNotNull(result.tagIv);
        assertNotNull(result.combined);

        // Combined format: base64(ciphertext + tag + iv)
        byte[] decrypted = service.aesDecryptCombined(result.combined, aesKeyBase64, aad);
        assertEquals(PLAIN_BODY, new String(decrypted));
    }

    @Test
    void aesShouldEncryptAndDecryptCiphertextAndTagWithAad() {
        String aad = buildStandardAad();

        CryptoService.EncryptResult result = service.aesEncrypt(PLAIN_BODY, aesKeyBase64, aad);

        // Decrypt using the separate ciphertext + tagIv method (as Go middleware does)
        byte[] decrypted = service.aesDecryptCiphertextAndTag(
                result.ciphertext, result.tagIv, aesKeyBase64, aad);
        assertEquals(PLAIN_BODY, new String(decrypted));
    }

    @Test
    void aesDecryptWithWrongAadShouldFail() {
        String aad = buildStandardAad();
        CryptoService.EncryptResult result = service.aesEncrypt(PLAIN_BODY, aesKeyBase64, aad);

        String wrongAad = "X-Request-ID=wrong&X-Timestamp=0";
        assertThrows(CryptoService.CryptoException.class, () ->
                service.aesDecryptCombined(result.combined, aesKeyBase64, wrongAad));
    }

    @Test
    void aesDecryptWithWrongKeyShouldFail() {
        String aad = buildStandardAad();
        CryptoService.EncryptResult result = service.aesEncrypt(PLAIN_BODY, aesKeyBase64, aad);
        String wrongKey = CryptoService.generateAesKey();

        assertThrows(CryptoService.CryptoException.class, () ->
                service.aesDecryptCombined(result.combined, wrongKey, aad));
    }

    // ---- Response encryption (empty AAD) ----

    @Test
    void responseEncryptionShouldUseEmptyAad() {
        String responseBody = "{\"code\":1,\"msg\":\"success\",\"data\":null}";

        CryptoService.EncryptResult result = service.aesEncrypt(responseBody, aesKeyBase64, "");
        assertNotNull(result.combined);

        byte[] decrypted = service.aesDecryptCombined(result.combined, aesKeyBase64, "");
        assertEquals(responseBody, new String(decrypted));
    }

    @Test
    void responseEncryptedWithEmptyAadShouldFailWithNonEmptyAad() {
        String responseBody = "{\"code\":1,\"msg\":\"success\"}";

        CryptoService.EncryptResult result = service.aesEncrypt(responseBody, aesKeyBase64, "");

        // Decrypting with non-empty AAD should fail
        assertThrows(CryptoService.CryptoException.class, () ->
                service.aesDecryptCombined(result.combined, aesKeyBase64, "some=aad"));
    }

    // ---- AAD construction ----

    @Test
    void aadShouldSortKeysAlphabetically() {
        // Keys intentionally in non-alphabetical order
        Map<String, String> params = new LinkedHashMap<>();
        params.put("X-Timestamp", String.valueOf(TIMESTAMP));
        params.put("pageSize", "10");
        params.put("X-Request-ID", REQUEST_ID);
        params.put("page", "1");

        String aad = service.buildAad(params);

        // AAD must have keys sorted alphabetically
        assertEquals("X-Request-ID=test-request-id-0001&X-Timestamp=1711411200000&page=1&pageSize=10", aad);
    }

    @Test
    void aadShouldMatchGoExpectedFormat() {
        // This is the exact AAD format Go's RequestAAD produces:
        // params: {X-Request-ID, X-Timestamp, ...query keys}, UriSort with filter always true
        Map<String, String> params = new LinkedHashMap<>();
        params.put("X-Request-ID", REQUEST_ID);
        params.put("X-Timestamp", String.valueOf(TIMESTAMP));
        params.put("page", "1");
        params.put("pageSize", "10");

        String aad = service.buildAad(params);
        assertEquals("X-Request-ID=test-request-id-0001&X-Timestamp=1711411200000&page=1&pageSize=10", aad);
    }

    @Test
    void aadShouldSkipEmptyValues() {
        Map<String, String> params = new LinkedHashMap<>();
        params.put("X-Request-ID", REQUEST_ID);
        params.put("empty", "");
        params.put("X-Timestamp", String.valueOf(TIMESTAMP));
        params.put("nullVal", null);

        String aad = service.buildAad(params);
        assertEquals("X-Request-ID=test-request-id-0001&X-Timestamp=1711411200000", aad);
    }

    @Test
    void aadWithSingleParamShouldNotHaveLeadingAmpersand() {
        Map<String, String> params = Map.of("X-Request-ID", REQUEST_ID);
        String aad = service.buildAad(params);
        assertEquals("X-Request-ID=test-request-id-0001", aad);
    }

    // ---- Sign verification ----

    @Test
    void signVerificationShouldPassWithCorrectAad() {
        String aad = buildStandardAad();

        // Encrypt empty body to get the tagIv as "signature"
        CryptoService.EncryptResult result = service.aesEncrypt("", aesKeyBase64, aad);

        // sign = tagIv (as used by Go's SignMiddleware)
        assertTrue(service.verifySign(result.tagIv, aesKeyBase64, aad));
    }

    @Test
    void signVerificationShouldFailWithWrongAad() {
        String aad = buildStandardAad();
        CryptoService.EncryptResult result = service.aesEncrypt("", aesKeyBase64, aad);

        assertFalse(service.verifySign(result.tagIv, aesKeyBase64, "wrong=aad"));
    }

    // ---- Combined format structural verification ----

    @Test
    void combinedFormatShouldContainCiphertextPlusTagPlusIv() {
        String aad = buildStandardAad();
        CryptoService.EncryptResult result = service.aesEncrypt(PLAIN_BODY, aesKeyBase64, aad);

        byte[] combined = java.util.Base64.getDecoder().decode(result.combined);
        // combined = ciphertext(plaintext length) + tag(16) + iv(12)
        // GCM ciphertext has same length as plaintext
        int expectedLen = PLAIN_BODY.getBytes().length + 16 + 12;
        assertEquals(expectedLen, combined.length);

        // Last 12 bytes = IV
        byte[] iv = new byte[12];
        System.arraycopy(combined, combined.length - 12, iv, 0, 12);
        assertEquals(12, iv.length);
    }

    @Test
    void tagIvFormatShouldContainTagPlusIv() {
        String aad = buildStandardAad();
        CryptoService.EncryptResult result = service.aesEncrypt(PLAIN_BODY, aesKeyBase64, aad);

        byte[] tagIv = java.util.Base64.getDecoder().decode(result.tagIv);
        assertEquals(28, tagIv.length); // 16 (tag) + 12 (iv)

        byte[] tag = new byte[16];
        byte[] iv = new byte[12];
        System.arraycopy(tagIv, 0, tag, 0, 16);
        System.arraycopy(tagIv, 16, iv, 0, 12);

        // Reconstruct combined: ciphertext + tag + iv
        byte[] ciphertext = java.util.Base64.getDecoder().decode(result.ciphertext);
        byte[] reconstructed = new byte[ciphertext.length + 16 + 12];
        System.arraycopy(ciphertext, 0, reconstructed, 0, ciphertext.length);
        System.arraycopy(tag, 0, reconstructed, ciphertext.length, 16);
        System.arraycopy(iv, 0, reconstructed, ciphertext.length + 16, 12);

        String reconstructedCombined = java.util.Base64.getEncoder().encodeToString(reconstructed);
        assertEquals(result.combined, reconstructedCombined);
    }

    // ---- Key generation ----

    @Test
    void generatedAesKeyShouldBe32BytesBase64() {
        String key = CryptoService.generateAesKey();
        byte[] decoded = java.util.Base64.getDecoder().decode(key);
        assertEquals(32, decoded.length);
    }

    @Test
    void generatedRsaKeyPairShouldBe2048Bits() {
        KeyPair kp = CryptoService.generateRsaKeyPair();
        assertEquals("RSA", kp.getPrivate().getAlgorithm());
        assertEquals(2048, ((java.security.interfaces.RSAPrivateKey) kp.getPrivate()).getModulus().bitLength());
    }

    // ---- PEM formatting ----

    @Test
    void pemPrivateKeyShouldBeParseable() {
        String pem = CryptoService.toPem(privateKey);
        assertTrue(pem.contains("-----BEGIN PRIVATE KEY-----"));
        assertTrue(pem.contains("-----END PRIVATE KEY-----"));

        PrivateKey parsed = CryptoService.parsePrivateKeyPem(pem);
        assertEquals(privateKey, parsed);
    }

    @Test
    void pemPublicKeyShouldBeParseable() {
        String pem = CryptoService.toPem(publicKey);
        assertTrue(pem.contains("-----BEGIN PUBLIC KEY-----"));
        assertTrue(pem.contains("-----END PUBLIC KEY-----"));

        PublicKey parsed = CryptoService.parsePublicKeyPem(pem);
        assertEquals(publicKey, parsed);
    }

    // ---- helpers ----

    private String buildStandardAad() {
        Map<String, String> params = new LinkedHashMap<>();
        params.put("X-Request-ID", REQUEST_ID);
        params.put("X-Timestamp", String.valueOf(TIMESTAMP));
        params.put("page", "1");
        params.put("pageSize", "10");
        return service.buildAad(params);
    }
}
