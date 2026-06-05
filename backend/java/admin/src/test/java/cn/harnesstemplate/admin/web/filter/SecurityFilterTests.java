package cn.harnesstemplate.admin.web.filter;

import cn.harnesstemplate.admin.application.port.UserLookupRepository;
import cn.harnesstemplate.admin.domain.model.SysRole;
import cn.harnesstemplate.admin.domain.model.SysUser;
import cn.harnesstemplate.admin.infrastructure.auth.CryptoService;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Primary;
import org.springframework.mock.web.MockMultipartFile;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.web.servlet.MockMvc;

import java.security.KeyPair;
import java.util.*;
import java.util.concurrent.ConcurrentHashMap;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

@SpringBootTest
@AutoConfigureMockMvc
@ActiveProfiles("test")
class SecurityFilterTests {

    @Autowired
    private MockMvc mockMvc;

    @Autowired
    private ServerKeyPairProvider serverKeyPairProvider;

    @Autowired
    private CryptoService cryptoService;

    @Autowired
    private UserLookupRepository userLookupRepository;

    private static final String CORRECT_PASSWORD = "correct";
    private static String hashedPassword;

    @BeforeEach
    void setUp() {
        if (hashedPassword == null) {
            hashedPassword = new BCryptPasswordEncoder(4).encode(CORRECT_PASSWORD);
        }
        StubUserRepository repo = (StubUserRepository) userLookupRepository;
        repo.clear();
        repo.addUser(new SysUser(1L, "admin", hashedPassword));
        repo.addRole(1L, new SysRole(1L, "admin"));
    }

    @TestConfiguration
    static class TestConfig {
        @Bean
        @Primary
        public UserLookupRepository userLookupRepository() {
            return new StubUserRepository();
        }
    }

    // ---- Timestamp filter tests ----

    @Test
    void expiredTimestampShouldReturnRequestExpiredCode() throws Exception {
        long expiredTime = System.currentTimeMillis() - (6 * 60 * 1000); // 6 minutes ago

        mockMvc.perform(post("/api/account/login/pwd")
                        .contentType("application/json")
                        .header("X-Timestamp", String.valueOf(expiredTime))
                        .content("{\"username\":\"admin\",\"pwd\":\"correct\"}"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.code").value(3))
                .andExpect(jsonPath("$.msg").value("请求已过期"));
    }

    @Test
    void validTimestampShouldPass() throws Exception {
        long now = System.currentTimeMillis();

        mockMvc.perform(post("/api/account/login/pwd")
                        .contentType("application/json")
                        .header("X-Timestamp", String.valueOf(now))
                        .content("{\"username\":\"admin\",\"pwd\":\"correct\"}"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.code").value(1));
    }

    @Test
    void missingTimestampShouldPass() throws Exception {
        // No X-Timestamp header at all — should pass through
        mockMvc.perform(post("/api/account/login/pwd")
                        .contentType("application/json")
                        .content("{\"username\":\"admin\",\"pwd\":\"correct\"}"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.code").value(1));
    }

    // ---- Language tests ----

    @Test
    void languageHeaderShouldBeStoredInRequest() throws Exception {
        mockMvc.perform(post("/api/account/login/pwd")
                        .contentType("application/json")
                        .header("X-Language", "en-US")
                        .content("{\"username\":\"admin\",\"pwd\":\"correct\"}"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.code").value(1))
                .andExpect(request().attribute("language", "en-US"));
    }

    @Test
    void missingLanguageHeaderShouldUseDefault() throws Exception {
        mockMvc.perform(post("/api/account/login/pwd")
                        .contentType("application/json")
                        .content("{\"username\":\"admin\",\"pwd\":\"correct\"}"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.code").value(1))
                .andExpect(request().attribute("language", "zh-CN"));
    }

    // ---- Encrypt filter tests ----

    @Test
    void requestWithoutEncryptionHeadersShouldPassThrough() throws Exception {
        mockMvc.perform(post("/api/account/login/pwd")
                        .contentType("application/json")
                        .content("{\"username\":\"admin\",\"pwd\":\"correct\"}"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.code").value(1))
                .andExpect(header().doesNotExist("X-Response-Is-Encrypt"));
    }

    @Test
    void encryptedRequestShouldReturnEncryptedResponse() throws Exception {
        // Encrypt a request body using the server's public key
        String aesKey = CryptoService.generateAesKey();
        String serverPublicKeyPem = serverKeyPairProvider.getPublicKey();

        // RSA-encrypt the AES key with server's public key
        String encryptedAesKey = cryptoService.rsaEncrypt(aesKey,
                CryptoService.parsePublicKeyPem(serverPublicKeyPem));

        // Build AAD
        Map<String, String> aadParams = new LinkedHashMap<>();
        aadParams.put("X-Request-ID", "test-req-id");
        aadParams.put("X-Timestamp", String.valueOf(System.currentTimeMillis()));
        String aad = cryptoService.buildAad(aadParams);

        // AES-encrypt the body
        String plainBody = "{\"username\":\"admin\",\"pwd\":\"correct\"}";
        CryptoService.EncryptResult result = cryptoService.aesEncrypt(plainBody, aesKey, aad);

        mockMvc.perform(post("/api/account/login/pwd")
                        .contentType("application/json")
                        .header("X-Request-Encrypted-Key", encryptedAesKey)
                        .header("X-Sign", result.tagIv)
                        .header("X-Request-ID", "test-req-id")
                        .header("X-Timestamp", aadParams.get("X-Timestamp"))
                        .content(result.ciphertext))
                .andExpect(status().isOk())
                .andExpect(header().string("X-Response-Is-Encrypt", "true"));
    }

    @Test
    void encryptedRequestWithFrontendTimestampHeaderShouldReturnEncryptedResponse() throws Exception {
        String aesKey = CryptoService.generateAesKey();
        String serverPublicKeyPem = serverKeyPairProvider.getPublicKey();
        String encryptedAesKey = cryptoService.rsaEncrypt(aesKey,
                CryptoService.parsePublicKeyPem(serverPublicKeyPem));

        Map<String, String> aadParams = new LinkedHashMap<>();
        aadParams.put("X-Request-ID", "frontend-req-id");
        aadParams.put("X-Request-Timestamp", String.valueOf(System.currentTimeMillis()));
        String aad = cryptoService.buildAad(aadParams);

        String plainBody = "{\"username\":\"admin\",\"pwd\":\"correct\"}";
        CryptoService.EncryptResult result = cryptoService.aesEncrypt(plainBody, aesKey, aad);

        mockMvc.perform(post("/api/account/login/pwd")
                        .contentType("application/json")
                        .header("X-Request-Encrypted-Key", encryptedAesKey)
                        .header("X-Request-Signature", result.tagIv)
                        .header("X-Request-ID", aadParams.get("X-Request-ID"))
                        .header("X-Request-Timestamp", aadParams.get("X-Request-Timestamp"))
                        .content(result.ciphertext))
                .andExpect(status().isOk())
                .andExpect(header().string("X-Response-Is-Encrypt", "true"));
    }

    @Test
    void encryptedBusinessFailureShouldStillReturnEncryptedResponse() throws Exception {
        String aesKey = CryptoService.generateAesKey();
        String serverPublicKeyPem = serverKeyPairProvider.getPublicKey();
        String encryptedAesKey = cryptoService.rsaEncrypt(aesKey,
                CryptoService.parsePublicKeyPem(serverPublicKeyPem));

        Map<String, String> aadParams = new LinkedHashMap<>();
        aadParams.put("X-Request-ID", "missing-user-req-id");
        aadParams.put("X-Request-Timestamp", String.valueOf(System.currentTimeMillis()));
        String aad = cryptoService.buildAad(aadParams);

        String plainBody = "{\"username\":\"root\",\"pwd\":\"correct\"}";
        CryptoService.EncryptResult result = cryptoService.aesEncrypt(plainBody, aesKey, aad);

        String encryptedResponse = mockMvc.perform(post("/api/account/login/pwd")
                        .contentType("application/json")
                        .header("X-Request-Encrypted-Key", encryptedAesKey)
                        .header("X-Request-Signature", result.tagIv)
                        .header("X-Request-ID", aadParams.get("X-Request-ID"))
                        .header("X-Request-Timestamp", aadParams.get("X-Request-Timestamp"))
                        .content(result.ciphertext))
                .andExpect(status().isOk())
                .andExpect(header().string("X-Response-Is-Encrypt", "true"))
                .andReturn()
                .getResponse()
                .getContentAsString();

        String decryptedResponse = new String(
                cryptoService.aesDecryptCombined(encryptedResponse, aesKey, ""),
                java.nio.charset.StandardCharsets.UTF_8
        );
        Integer code = com.jayway.jsonpath.JsonPath.read(decryptedResponse, "$.code");
        String msg = com.jayway.jsonpath.JsonPath.read(decryptedResponse, "$.msg");
        org.assertj.core.api.Assertions.assertThat(code).isEqualTo(2);
        org.assertj.core.api.Assertions.assertThat(msg).isEqualTo("用户名或密码无效");
    }

    @Test
    void encryptedRequestWithStalePublicKeyShouldReturnRequestKeyFailure() throws Exception {
        String aesKey = CryptoService.generateAesKey();
        KeyPair staleKeyPair = CryptoService.generateRsaKeyPair();
        String encryptedAesKey = cryptoService.rsaEncrypt(aesKey, staleKeyPair.getPublic());

        Map<String, String> aadParams = new LinkedHashMap<>();
        aadParams.put("X-Request-ID", "stale-key-req-id");
        aadParams.put("X-Request-Timestamp", String.valueOf(System.currentTimeMillis()));
        String aad = cryptoService.buildAad(aadParams);

        String plainBody = "{\"username\":\"admin\",\"pwd\":\"correct\"}";
        CryptoService.EncryptResult result = cryptoService.aesEncrypt(plainBody, aesKey, aad);

        mockMvc.perform(post("/api/account/login/pwd")
                        .contentType("application/json")
                        .header("X-Request-Encrypted-Key", encryptedAesKey)
                        .header("X-Request-Signature", result.tagIv)
                        .header("X-Request-ID", aadParams.get("X-Request-ID"))
                        .header("X-Request-Timestamp", aadParams.get("X-Request-Timestamp"))
                        .content(result.ciphertext))
                .andExpect(status().isOk())
                .andExpect(header().doesNotExist("X-Response-Is-Encrypt"))
                .andExpect(jsonPath("$.code").value(5))
                .andExpect(jsonPath("$.msg").value("请求错误"));
    }

    @Test
    void requestBodyShouldRemainReadableWhenEncryptedKeyExistsButSignatureIsMissing() throws Exception {
        String aesKey = CryptoService.generateAesKey();
        String serverPublicKeyPem = serverKeyPairProvider.getPublicKey();
        String encryptedAesKey = cryptoService.rsaEncrypt(aesKey,
                CryptoService.parsePublicKeyPem(serverPublicKeyPem));

        String encryptedResponse = mockMvc.perform(post("/api/account/login/pwd")
                        .contentType("application/json")
                        .header("X-Request-Encrypted-Key", encryptedAesKey)
                        .header("X-Request-ID", "missing-sign-req-id")
                        .header("X-Request-Timestamp", String.valueOf(System.currentTimeMillis()))
                        .content("{\"username\":\"admin\",\"pwd\":\"correct\"}"))
                .andExpect(status().isOk())
                .andExpect(header().string("X-Response-Is-Encrypt", "true"))
                .andReturn()
                .getResponse()
                .getContentAsString();

        String decryptedResponse = new String(
                cryptoService.aesDecryptCombined(encryptedResponse, aesKey, ""),
                java.nio.charset.StandardCharsets.UTF_8
        );
        Integer code = com.jayway.jsonpath.JsonPath.read(decryptedResponse, "$.code");
        org.assertj.core.api.Assertions.assertThat(code).isEqualTo(1);
    }

    // ---- Stub repository ----

    static class StubUserRepository implements UserLookupRepository {
        private final Map<String, SysUser> users = new ConcurrentHashMap<>();
        private final Map<Long, List<SysRole>> userRoles = new ConcurrentHashMap<>();

        void addUser(SysUser user) { users.put(user.getUsername(), user); }
        void addRole(Long userId, SysRole role) {
            userRoles.computeIfAbsent(userId, k -> new java.util.ArrayList<>()).add(role);
        }
        void clear() { users.clear(); userRoles.clear(); }

        @Override
        public Optional<SysUser> findByUsername(String username) {
            return Optional.ofNullable(users.get(username));
        }

        @Override
        public List<SysRole> findRolesByUserId(Long userId) {
            return userRoles.getOrDefault(userId, List.of());
        }
    }
}
