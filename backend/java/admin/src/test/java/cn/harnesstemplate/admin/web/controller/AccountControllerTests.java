package cn.harnesstemplate.admin.web.controller;

import cn.harnesstemplate.admin.application.port.UserLookupRepository;
import cn.harnesstemplate.admin.domain.model.SysRole;
import cn.harnesstemplate.admin.domain.model.SysUser;
import cn.harnesstemplate.admin.infrastructure.auth.SessionInfo;
import cn.dev33.satoken.stp.StpUtil;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Primary;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.MvcResult;

import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.ConcurrentHashMap;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

@SpringBootTest
@AutoConfigureMockMvc
@ActiveProfiles("test")
class AccountControllerTests {

    @Autowired
    private MockMvc mockMvc;

    @Autowired
    private UserLookupRepository userLookupRepository;

    private static final String CORRECT_PASSWORD = "correct";
    private static String hashedPassword;

    @BeforeEach
    void setUp() {
        // Pre-hash the password with bcrypt cost=4 matching Go's bcrypt.MinCost
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

    // ---- Login tests ----

    @Test
    void publicKeyEndpointShouldBeAccessibleWithoutAuth() throws Exception {
        MvcResult result = mockMvc.perform(get("/api/encrypt/public/key"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.code").value(1))
                .andExpect(jsonPath("$.data.publicKey").isString())
                .andReturn();

        assertPublicKeyMatchesGoContract(result.getResponse().getContentAsString());
    }

    @Test
    void loginWithInvalidUsernameShouldReturnBusinessFailure() throws Exception {
        mockMvc.perform(post("/api/account/login/pwd")
                        .contentType("application/json")
                        .content("{\"username\":\"nobody\",\"pwd\":\"whatever\"}"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.code").value(2))
                .andExpect(jsonPath("$.msg").value("用户名或密码无效"));
    }

    @Test
    void loginWithWrongPasswordShouldReturnBusinessFailure() throws Exception {
        mockMvc.perform(post("/api/account/login/pwd")
                        .contentType("application/json")
                        .content("{\"username\":\"admin\",\"pwd\":\"wrongpassword\"}"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.code").value(2))
                .andExpect(jsonPath("$.msg").value("用户名或密码无效"));
    }

    @Test
    void loginWithValidCredentialsShouldReturnTokenAndPublicKey() throws Exception {
        MvcResult result = mockMvc.perform(post("/api/account/login/pwd")
                        .contentType("application/json")
                        .content("{\"username\":\"admin\",\"pwd\":\"" + CORRECT_PASSWORD + "\"}"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.code").value(1))
                .andExpect(jsonPath("$.msg").value("success"))
                .andExpect(jsonPath("$.data.token").isNotEmpty())
                .andExpect(jsonPath("$.data.publicKey").isString())
                .andExpect(cookie().exists("token"))
                .andReturn();

        assertPublicKeyMatchesGoContract(result.getResponse().getContentAsString());
    }

    @Test
    void loginShouldStoreSessionInfo() throws Exception {
        String responseBody = mockMvc.perform(post("/api/account/login/pwd")
                        .contentType("application/json")
                        .content("{\"username\":\"admin\",\"pwd\":\"" + CORRECT_PASSWORD + "\"}"))
                .andExpect(status().isOk())
                .andReturn().getResponse().getContentAsString();

        // Extract token from response
        String token = com.jayway.jsonpath.JsonPath.read(responseBody, "$.data.token");

        // Verify session info via Sa-Token API (getLoginIdByToken works without request context)
        Object loginId = StpUtil.getLoginIdByToken(token);
        SessionInfo sessionInfo = (SessionInfo) StpUtil.getSessionByLoginId(loginId)
                .get(SessionInfo.SESSION_KEY);
        assert sessionInfo != null : "Session info should be stored";
        assert sessionInfo.getId().equals(1L) : "User ID should be 1";
        assert sessionInfo.getUsername().equals("admin") : "Username should be admin";
        assert sessionInfo.getRoleCodes().contains("admin") : "Should have admin role";
        assert sessionInfo.getRoleIds().contains(1L) : "Should have role ID 1";
        assert sessionInfo.getPrivateKey() != null && !sessionInfo.getPrivateKey().isEmpty()
                : "Private key should be stored";
        assert sessionInfo.getPrivateKey().contains("BEGIN PRIVATE KEY")
                : "Private key should be PEM format";
    }

    // ---- Logout tests ----

    @Test
    void logoutShouldClearSession() throws Exception {
        // First login
        String responseBody = mockMvc.perform(post("/api/account/login/pwd")
                        .contentType("application/json")
                        .content("{\"username\":\"admin\",\"pwd\":\"" + CORRECT_PASSWORD + "\"}"))
                .andExpect(status().isOk())
                .andReturn().getResponse().getContentAsString();

        String token = com.jayway.jsonpath.JsonPath.read(responseBody, "$.data.token");

        // Then logout with cookie
        mockMvc.perform(get("/api/account/logout")
                        .cookie(new jakarta.servlet.http.Cookie("token", token)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.code").value(1));

        // Verify the token is no longer valid
        Object loginId = StpUtil.getLoginIdByToken(token);
        assert loginId == null : "Login session should be cleared after logout";
    }

    @Test
    void logoutWithoutAuthShouldReturnFailure() throws Exception {
        mockMvc.perform(get("/api/account/logout"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.code").value(2));
    }

    // ---- Stub repository ----

    static class StubUserRepository implements UserLookupRepository {
        private final Map<String, SysUser> users = new ConcurrentHashMap<>();
        private final Map<Long, List<SysRole>> userRoles = new ConcurrentHashMap<>();

        void addUser(SysUser user) {
            users.put(user.getUsername(), user);
        }

        void addRole(Long userId, SysRole role) {
            userRoles.computeIfAbsent(userId, k -> new java.util.ArrayList<>()).add(role);
        }

        void clear() {
            users.clear();
            userRoles.clear();
        }

        @Override
        public Optional<SysUser> findByUsername(String username) {
            return Optional.ofNullable(users.get(username));
        }

        @Override
        public List<SysRole> findRolesByUserId(Long userId) {
            return userRoles.getOrDefault(userId, List.of());
        }
    }

    private static void assertPublicKeyMatchesGoContract(String responseBody) {
        String publicKey = com.jayway.jsonpath.JsonPath.read(responseBody, "$.data.publicKey");

        assert !publicKey.contains("BEGIN PUBLIC KEY") : "publicKey should not be PEM formatted";
        assert !publicKey.contains("END PUBLIC KEY") : "publicKey should not be PEM formatted";
        assert publicKey.chars().noneMatch(Character::isWhitespace) : "publicKey should not contain whitespace";
        assert publicKey.matches("^[A-Za-z0-9+/]+={0,2}$") : "publicKey should be base64 encoded";
    }
}
