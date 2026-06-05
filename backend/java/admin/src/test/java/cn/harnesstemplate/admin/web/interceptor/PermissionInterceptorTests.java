package cn.harnesstemplate.admin.web.interceptor;

import cn.harnesstemplate.admin.application.port.UserLookupRepository;
import cn.harnesstemplate.admin.domain.model.SysRole;
import cn.harnesstemplate.admin.domain.model.SysUser;
import cn.harnesstemplate.admin.infrastructure.permission.CasbinPermissionService;
import cn.harnesstemplate.admin.web.dto.R;
import com.jayway.jsonpath.JsonPath;
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
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.*;
import java.util.concurrent.ConcurrentHashMap;

import static org.assertj.core.api.Assertions.assertThat;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

@SpringBootTest
@AutoConfigureMockMvc
@ActiveProfiles("test")
class PermissionInterceptorTests {

    @Autowired
    private MockMvc mockMvc;

    @Autowired
    private CasbinPermissionService casbinPermissionService;

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
        repo.addUser(new SysUser(2L, "viewer", hashedPassword));
        repo.addRole(1L, new SysRole(1L, "admin"));
        repo.addRole(2L, new SysRole(2L, "viewer"));

        // Set up policies
        casbinPermissionService.removePoliciesForSubject("role:admin");
        casbinPermissionService.removePoliciesForSubject("role:viewer");
        casbinPermissionService.addPolicy("role:admin", "/api/test/protected", "GET");
        casbinPermissionService.addPolicy("role:admin", "/api/test/admin-only", "POST");
        casbinPermissionService.addPolicy("role:viewer", "/api/test/protected", "GET");
    }

    @TestConfiguration
    static class TestConfig {
        @Bean
        @Primary
        public UserLookupRepository userLookupRepository() {
            return new StubUserRepository();
        }

        @Bean
        public TestProtectedController testProtectedController() {
            return new TestProtectedController();
        }
    }

    @RestController
    @RequestMapping("/api/test")
    static class TestProtectedController {
        @GetMapping("/protected")
        public R<String> protectedEndpoint() {
            return R.ok("protected data");
        }

        @GetMapping("/admin-only")
        public R<String> adminOnly() {
            return R.ok("admin data");
        }

        @GetMapping("/public")
        public R<String> publicEndpoint() {
            return R.ok("public data");
        }
    }

    // ---- Tests ----

    @Test
    void roleWithPermissionShouldAccessProtectedRoute() throws Exception {
        String token = loginAs("admin");

        mockMvc.perform(get("/api/test/protected")
                        .cookie(new jakarta.servlet.http.Cookie("token", token)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.code").value(1));
    }

    @Test
    void roleWithoutPermissionShouldReceiveUnauthorized() throws Exception {
        String token = loginAs("viewer");

        mockMvc.perform(get("/api/test/admin-only")
                        .cookie(new jakarta.servlet.http.Cookie("token", token)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.code").value(200))
                .andExpect(jsonPath("$.msg").value("未授权"));
    }

    @Test
    void unauthenticatedRequestShouldReceiveUnauthorized() throws Exception {
        mockMvc.perform(get("/api/test/protected"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.code").value(200))
                .andExpect(jsonPath("$.msg").value("未授权"));
    }

    @Test
    void adminRoleWithWildcardPolicyShouldMatchViaKeyMatch2() throws Exception {
        // Add wildcard policy for admin role
        casbinPermissionService.addPolicy("role:admin", "/api/sys/*", "GET");

        List<String> subjects = CasbinPermissionService.buildSubjects(1L, List.of("admin"));

        // keyMatch2: /api/sys/* should match /api/sys/user/page
        assertThat(casbinPermissionService.enforce(subjects, "/api/sys/user/page", "GET")).isTrue();
        assertThat(casbinPermissionService.enforce(subjects, "/api/sys/role/list", "GET")).isTrue();
        assertThat(casbinPermissionService.enforce(subjects, "/api/other/endpoint", "GET")).isFalse();
    }

    // ---- Helpers ----

    private String loginAs(String username) throws Exception {
        String responseBody = mockMvc.perform(post("/api/account/login/pwd")
                        .contentType("application/json")
                        .content("{\"username\":\"" + username + "\",\"pwd\":\"" + CORRECT_PASSWORD + "\"}"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.code").value(1))
                .andReturn().getResponse().getContentAsString();

        return JsonPath.read(responseBody, "$.data.token");
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
