package cn.harnesstemplate.admin.infrastructure.persistence;

import cn.harnesstemplate.admin.application.service.system.SysUserQueryService;
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

import java.time.LocalDateTime;

import static org.assertj.core.api.Assertions.assertThat;

@SpringBootTest
@ActiveProfiles("test-db")
@Testcontainers
@AutoConfigureTestDatabase(replace = AutoConfigureTestDatabase.Replace.NONE)
class EasyQuerySupportTests {

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
    private EasyQuerySupport easyQuerySupport;

    @Autowired
    private SysUserQueryService sysUserQueryService;

    @BeforeEach
    void setUp() {
        var users = sysUserQueryService.listAll();
        for (var user : users) {
            sysUserQueryService.delete(user.getId());
        }
    }

    @Test
    void shouldExposeQueryableAndUpdatableOrmApi() {
        SysUser user = new SysUser();
        user.setUsername("orm-user");
        user.setNickname("Before");
        user.setPassword("pw");
        user.setLanguageCode("zh-CN");
        user.setIsEnabled(true);
        sysUserQueryService.save(user);

        SysUser found = easyQuerySupport.queryable(SysUser.class)
                .where(w -> {
                    w.eq("deletedAt", 0L);
                    w.like("username", "orm");
                })
                .orderByAsc(o -> o.column("id"))
                .firstOrNull();

        assertThat(found).isNotNull();
        assertThat(found.getId()).isEqualTo(user.getId());

        long updatedRows = easyQuerySupport.updatable(SysUser.class)
                .set("nickname", "After")
                .set("updatedAt", LocalDateTime.now())
                .where(w -> {
                    w.eq("id", user.getId());
                    w.eq("deletedAt", 0L);
                })
                .executeRows();

        assertThat(updatedRows).isEqualTo(1);
        assertThat(sysUserQueryService.findById(user.getId()).getNickname()).isEqualTo("After");
    }
}
