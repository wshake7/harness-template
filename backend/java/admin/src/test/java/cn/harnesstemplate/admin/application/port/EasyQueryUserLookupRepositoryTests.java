package cn.harnesstemplate.admin.application.port;

import cn.harnesstemplate.admin.domain.entity.SysRole;
import cn.harnesstemplate.admin.domain.entity.SysUser;
import cn.harnesstemplate.admin.domain.entity.SysUserRole;
import cn.harnesstemplate.admin.infrastructure.persistence.EasyQuerySupport;
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

import java.util.List;

import static org.assertj.core.api.Assertions.assertThat;

@SpringBootTest
@ActiveProfiles("test-db")
@Testcontainers
@AutoConfigureTestDatabase(replace = AutoConfigureTestDatabase.Replace.NONE)
class EasyQueryUserLookupRepositoryTests {

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
    private UserLookupRepository repository;

    @Autowired
    private EasyQuerySupport easyQuerySupport;

    @BeforeEach
    void setUp() {
        easyQuerySupport.execute("DELETE FROM sys_user_role");
        easyQuerySupport.execute("DELETE FROM sys_role");
        easyQuerySupport.execute("DELETE FROM sys_user");
    }

    @Test
    void shouldFindUserAndRolesFromDatabase() {
        SysUser user = new SysUser();
        user.setUsername("root");
        user.setNickname("Root");
        user.setPassword("hashed-password");
        user.setLanguageCode("zh-CN");
        user.setIsEnabled(true);
        user.setDeletedAt(0L);
        easyQuerySupport.insert(user);

        SysRole rootRole = new SysRole();
        rootRole.setName("Root");
        rootRole.setCode("role:root");
        rootRole.setIsEnabled(true);
        rootRole.setDeletedAt(0L);
        easyQuerySupport.insert(rootRole);

        SysRole disabledRole = new SysRole();
        disabledRole.setName("Disabled");
        disabledRole.setCode("role:disabled");
        disabledRole.setIsEnabled(false);
        disabledRole.setDeletedAt(0L);
        easyQuerySupport.insert(disabledRole);

        SysUserRole activeLink = new SysUserRole();
        activeLink.setUserId(user.getId());
        activeLink.setRoleId(rootRole.getId());
        activeLink.setDeletedAt(0L);
        easyQuerySupport.insert(activeLink);

        SysUserRole disabledLink = new SysUserRole();
        disabledLink.setUserId(user.getId());
        disabledLink.setRoleId(disabledRole.getId());
        disabledLink.setDeletedAt(0L);
        easyQuerySupport.insert(disabledLink);

        var foundUser = repository.findByUsername("root");
        assertThat(foundUser).isPresent();
        assertThat(foundUser.get().getUsername()).isEqualTo("root");
        assertThat(foundUser.get().getPassword()).isEqualTo("hashed-password");

        List<cn.harnesstemplate.admin.domain.model.SysRole> roles = repository.findRolesByUserId(user.getId());
        assertThat(roles).extracting(cn.harnesstemplate.admin.domain.model.SysRole::getCode)
                .containsExactly("role:root");
    }
}
