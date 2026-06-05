package cn.harnesstemplate.admin.application.converter;

import cn.harnesstemplate.admin.domain.entity.SysRole;
import cn.harnesstemplate.admin.domain.entity.SysUser;
import cn.harnesstemplate.admin.web.dto.*;
import io.github.linpeilie.Converter;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.ActiveProfiles;

import java.time.LocalDateTime;
import java.util.List;

import static org.assertj.core.api.Assertions.assertThat;

@SpringBootTest
@ActiveProfiles("test")
class ConverterTests {

    @Autowired
    private Converter converter;

    @Test
    void userEntityToResponseShouldExcludePassword() {
        SysUser user = new SysUser();
        user.setId(1L);
        user.setUsername("admin");
        user.setNickname("Admin");
        user.setPassword("secret-hash");
        user.setLanguageCode("zh-CN");
        user.setIsEnabled(true);
        user.setCreatedAt(LocalDateTime.now());
        user.setUpdatedAt(LocalDateTime.now());

        SysUserResponse response = converter.convert(user, SysUserResponse.class);

        assertThat(response.id()).isEqualTo(1L);
        assertThat(response.username()).isEqualTo("admin");
        assertThat(response.nickname()).isEqualTo("Admin");
        assertThat(response.languageCode()).isEqualTo("zh-CN");
        assertThat(response.isEnabled()).isTrue();
        // Password field should NOT exist on SysUserResponse
        assertThat(response.getClass().getDeclaredFields())
                .noneMatch(f -> f.getName().equals("password"));
    }

    @Test
    void userCreateRequestToEntityShouldNotSetServerFields() {
        SysUserCreateRequest request = new SysUserCreateRequest(
                "newuser",
                "New User",
                "plain-password",
                "en-US",
                null
        );

        SysUser entity = converter.convert(request, SysUser.class);

        assertThat(entity.getUsername()).isEqualTo("newuser");
        assertThat(entity.getNickname()).isEqualTo("New User");
        assertThat(entity.getPassword()).isEqualTo("plain-password");
        assertThat(entity.getLanguageCode()).isEqualTo("en-US");
        // Server-owned fields should be null/default
        assertThat(entity.getId()).isNull();
        assertThat(entity.getCreatedAt()).isNull();
        assertThat(entity.getUpdatedAt()).isNull();
    }

    @Test
    void roleEntityToResponseShouldPreserveFields() {
        SysRole role = new SysRole();
        role.setId(1L);
        role.setName("管理员");
        role.setCode("admin");
        role.setRemark("Administrator role");
        role.setIsEnabled(true);
        role.setCreatedAt(LocalDateTime.now());
        role.setUpdatedAt(LocalDateTime.now());

        SysRoleResponse response = converter.convert(role, SysRoleResponse.class);

        assertThat(response.id()).isEqualTo(1L);
        assertThat(response.name()).isEqualTo("管理员");
        assertThat(response.code()).isEqualTo("admin");
        assertThat(response.remark()).isEqualTo("Administrator role");
        assertThat(response.isEnabled()).isTrue();
    }

    @Test
    void roleCreateRequestToEntityShouldNotSetServerFields() {
        SysRoleCreateRequest request = new SysRoleCreateRequest(
                "Viewer",
                "viewer",
                null,
                "Read-only role",
                null,
                null
        );

        SysRole entity = converter.convert(request, SysRole.class);

        assertThat(entity.getName()).isEqualTo("Viewer");
        assertThat(entity.getCode()).isEqualTo("viewer");
        assertThat(entity.getRemark()).isEqualTo("Read-only role");
        assertThat(entity.getId()).isNull();
        assertThat(entity.getCreatedAt()).isNull();
        assertThat(entity.getUpdatedAt()).isNull();
    }

    @Test
    void userResponseShouldNotContainPasswordField() {
        // Structural test: SysUserResponse must not declare a password field
        java.lang.reflect.Field[] fields = SysUserResponse.class.getDeclaredFields();
        for (java.lang.reflect.Field f : fields) {
            assertThat(f.getName()).isNotIn("password", "privateKey", "token");
        }
    }

    @Test
    void userResponseListConversionShouldWork() {
        SysUser user1 = new SysUser();
        user1.setId(1L);
        user1.setUsername("user1");
        SysUser user2 = new SysUser();
        user2.setId(2L);
        user2.setUsername("user2");

        List<SysUserResponse> responses = converter.convert(
                List.of(user1, user2), SysUserResponse.class);

        assertThat(responses).hasSize(2);
        assertThat(responses.get(0).username()).isEqualTo("user1");
        assertThat(responses.get(1).username()).isEqualTo("user2");
    }
}
