package cn.harnesstemplate.admin.application.service.system;

import cn.harnesstemplate.admin.domain.entity.SysUser;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.jdbc.AutoConfigureTestDatabase;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import static org.assertj.core.api.Assertions.assertThat;

@SpringBootTest
@ActiveProfiles("test-db")
@Testcontainers
@AutoConfigureTestDatabase(replace = AutoConfigureTestDatabase.Replace.NONE)
class SysUserQueryServiceTests {

    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:16-alpine")
            .withDatabaseName("testdb")
            .withUsername("test")
            .withPassword("test");

    @DynamicPropertySource
    static void configureProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", postgres::getJdbcUrl);
        registry.add("spring.datasource.username", postgres::getUsername);
        registry.add("spring.datasource.password", postgres::getPassword);
        registry.add("spring.datasource.driver-class-name", () -> "org.postgresql.Driver");
    }

    @Autowired
    private SysUserQueryService service;

    @BeforeEach
    void setUp() {
        // Clean up
        var users = service.listAll();
        for (var u : users) {
            service.delete(u.getId());
        }
    }

    @Test
    void shouldSaveAndFindUser() {
        SysUser user = new SysUser();
        user.setUsername("testuser");
        user.setNickname("Test User");
        user.setPassword("hashed");
        user.setLanguageCode("zh-CN");
        user.setIsEnabled(true);
        service.save(user);

        assertThat(user.getId()).isNotNull();

        SysUser found = service.findById(user.getId());
        assertThat(found).isNotNull();
        assertThat(found.getUsername()).isEqualTo("testuser");
    }

    @Test
    void shouldFindByUsername() {
        SysUser user = new SysUser();
        user.setUsername("admin");
        user.setNickname("Admin");
        user.setPassword("secret");
        user.setLanguageCode("zh-CN");
        user.setIsEnabled(true);
        service.save(user);

        SysUser found = service.findByUsername("admin");
        assertThat(found).isNotNull();
        assertThat(found.getNickname()).isEqualTo("Admin");
    }

    @Test
    void shouldPaginateWithFilter() {
        for (int i = 1; i <= 5; i++) {
            SysUser user = new SysUser();
            user.setUsername("user" + i);
            user.setNickname("User " + i);
            user.setPassword("pw");
            user.setLanguageCode("zh-CN");
            user.setIsEnabled(true);
            service.save(user);
        }

        var list = service.page(1, 3, null);
        assertThat(list).hasSizeLessThanOrEqualTo(3);
        long total = service.count(null);
        assertThat(total).isGreaterThanOrEqualTo(5);
    }

    @Test
    void shouldFilterByUsername() {
        SysUser user = new SysUser();
        user.setUsername("special-user");
        user.setNickname("Special");
        user.setPassword("pw");
        user.setLanguageCode("zh-CN");
        user.setIsEnabled(true);
        service.save(user);

        var list = service.page(1, 10, "special");
        assertThat(list).hasSizeGreaterThanOrEqualTo(1);
    }

    @Test
    void shouldUpdateUser() {
        SysUser user = new SysUser();
        user.setUsername("update-me");
        user.setNickname("Original");
        user.setPassword("pw");
        user.setLanguageCode("zh-CN");
        user.setIsEnabled(true);
        service.save(user);

        user.setNickname("Updated");
        service.update(user);

        SysUser reloaded = service.findById(user.getId());
        assertThat(reloaded.getNickname()).isEqualTo("Updated");
    }

    @Test
    void shouldDeleteUser() {
        SysUser user = new SysUser();
        user.setUsername("delete-me");
        user.setNickname("Delete");
        user.setPassword("pw");
        user.setLanguageCode("zh-CN");
        user.setIsEnabled(true);
        service.save(user);

        assertThat(service.count(null)).isEqualTo(1);
        service.delete(user.getId());
        assertThat(service.count(null)).isEqualTo(0);
    }
}
