package cn.harnesstemplate.admin.application.service.system;

import cn.harnesstemplate.admin.domain.entity.SysRole;
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
class SysRoleQueryServiceTests {

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
    private SysRoleQueryService service;

    @BeforeEach
    void setUp() {
        var roles = service.listAll();
        for (var r : roles) {
            service.delete(r.getId());
        }
    }

    @Test
    void shouldSaveAndFindRole() {
        SysRole role = new SysRole();
        role.setName("管理员");
        role.setCode("admin");
        role.setIsEnabled(true);
        service.save(role);

        assertThat(role.getId()).isNotNull();

        SysRole found = service.findById(role.getId());
        assertThat(found).isNotNull();
        assertThat(found.getCode()).isEqualTo("admin");
    }

    @Test
    void shouldFindByCode() {
        SysRole role = new SysRole();
        role.setName("User");
        role.setCode("user");
        role.setIsEnabled(true);
        service.save(role);

        SysRole found = service.findByCode("user");
        assertThat(found).isNotNull();
        assertThat(found.getName()).isEqualTo("User");
    }

    @Test
    void shouldReturnNullWhenNotFound() {
        assertThat(service.findByCode("nonexistent")).isNull();
        assertThat(service.findById(99999L)).isNull();
    }

    @Test
    void shouldUpdateRole() {
        SysRole role = new SysRole();
        role.setName("Original");
        role.setCode("original");
        role.setIsEnabled(true);
        service.save(role);

        role.setName("Updated");
        service.update(role);

        SysRole reloaded = service.findById(role.getId());
        assertThat(reloaded.getName()).isEqualTo("Updated");
    }

    @Test
    void shouldCreateHierarchy() {
        SysRole parent = new SysRole();
        parent.setName("Parent");
        parent.setCode("parent");
        parent.setIsEnabled(true);
        service.save(parent);

        SysRole child = new SysRole();
        child.setName("Child");
        child.setCode("child");
        child.setParentId(parent.getId());
        child.setIsEnabled(true);
        service.save(child);

        assertThat(service.hasChildren(parent.getId())).isTrue();
        assertThat(service.hasChildren(child.getId())).isFalse();
    }

    @Test
    void shouldSoftDelete() {
        SysRole role = new SysRole();
        role.setName("DeleteMe");
        role.setCode("delete-me");
        role.setIsEnabled(true);
        service.save(role);

        assertThat(service.count(null)).isEqualTo(1);
        service.delete(role.getId());
        assertThat(service.count(null)).isEqualTo(0);
    }
}
