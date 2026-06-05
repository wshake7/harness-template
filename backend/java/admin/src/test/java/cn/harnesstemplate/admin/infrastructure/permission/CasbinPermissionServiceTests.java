package cn.harnesstemplate.admin.infrastructure.permission;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.List;

import static org.assertj.core.api.Assertions.assertThat;

class CasbinPermissionServiceTests {

    private CasbinPermissionService service;

    @BeforeEach
    void setUp() {
        service = new CasbinPermissionService();
    }

    @Test
    void roleWithApiPermissionShouldAllowAccess() {
        service.addPolicy("role:admin", "/api/sys/user/page", "GET");
        service.addPolicy("role:admin", "/api/sys/user/save", "POST");

        List<String> subjects = CasbinPermissionService.buildSubjects(1L, List.of("admin"));

        assertThat(service.enforce(subjects, "/api/sys/user/page", "GET")).isTrue();
    }

    @Test
    void roleWithoutApiPermissionShouldBeDenied() {
        service.addPolicy("role:viewer", "/api/sys/user/page", "GET");

        List<String> subjects = CasbinPermissionService.buildSubjects(1L, List.of("guest"));

        assertThat(service.enforce(subjects, "/api/sys/user/save", "POST")).isFalse();
    }

    @Test
    void userSubjectShouldBeChecked() {
        service.addPolicy("user:42", "/api/sys/user/profile", "GET");

        List<String> subjects = CasbinPermissionService.buildSubjects(42L, List.of());

        assertThat(service.enforce(subjects, "/api/sys/user/profile", "GET")).isTrue();
    }

    @Test
    void orLogicShouldAllowIfAnySubjectMatches() {
        service.addPolicy("role:viewer", "/api/sys/user/page", "GET");

        List<String> subjects = CasbinPermissionService.buildSubjects(1L, List.of("admin", "viewer"));

        // Fails for "user:1", "role:admin" but passes for "role:viewer"
        assertThat(service.enforce(subjects, "/api/sys/user/page", "GET")).isTrue();
    }

    @Test
    void removePolicyShouldWork() {
        service.addPolicy("role:admin", "/api/sys/user/delete", "DELETE");
        assertThat(service.enforce(
                CasbinPermissionService.buildSubjects(1L, List.of("admin")),
                "/api/sys/user/delete", "DELETE")).isTrue();

        service.removePolicy("role:admin", "/api/sys/user/delete", "DELETE");
        assertThat(service.enforce(
                CasbinPermissionService.buildSubjects(1L, List.of("admin")),
                "/api/sys/user/delete", "DELETE")).isFalse();
    }

    @Test
    void removePoliciesForSubjectShouldRemoveAll() {
        service.addPolicy("role:temp", "/api/a", "GET");
        service.addPolicy("role:temp", "/api/b", "POST");

        service.removePoliciesForSubject("role:temp");

        List<String> subjects = CasbinPermissionService.buildSubjects(1L, List.of("temp"));
        assertThat(service.enforce(subjects, "/api/a", "GET")).isFalse();
        assertThat(service.enforce(subjects, "/api/b", "POST")).isFalse();
    }

    @Test
    void caseInsensitiveMethodMatching() {
        service.addPolicy("role:admin", "/api/sys/user/page", "GET");

        List<String> subjects = CasbinPermissionService.buildSubjects(1L, List.of("admin"));

        assertThat(service.enforce(subjects, "/api/sys/user/page", "get")).isTrue();
        assertThat(service.enforce(subjects, "/api/sys/user/page", "Get")).isTrue();
    }

    @Test
    void keyMatch2WildcardShouldWork() {
        service.addPolicy("role:admin", "/api/sys/*", "GET");

        List<String> subjects = CasbinPermissionService.buildSubjects(1L, List.of("admin"));

        assertThat(service.enforce(subjects, "/api/sys/user/page", "GET")).isTrue();
        assertThat(service.enforce(subjects, "/api/sys/role/list", "GET")).isTrue();
    }

    @Test
    void keyMatch2ResourcePatternShouldWork() {
        // keyMatch2 supports :resource pattern matching
        service.addPolicy("role:admin", "/api/sys/:resource", "GET");

        List<String> subjects = CasbinPermissionService.buildSubjects(1L, List.of("admin"));

        assertThat(service.enforce(subjects, "/api/sys/user", "GET")).isTrue();
        assertThat(service.enforce(subjects, "/api/other/thing", "GET")).isFalse();
    }

    @Test
    void buildSubjectsWithEmptyRoleCodes() {
        List<String> subjects = CasbinPermissionService.buildSubjects(1L, List.of());

        assertThat(subjects).containsExactly("user:1");
    }

    @Test
    void buildSubjectsWithBlankRoles() {
        List<String> subjects = CasbinPermissionService.buildSubjects(1L,
                List.of("admin", "  ", "", "viewer"));

        assertThat(subjects).containsExactly("user:1", "role:admin", "role:viewer");
    }

    @Test
    void noPoliciesShouldDenyAll() {
        List<String> subjects = CasbinPermissionService.buildSubjects(1L, List.of("admin"));

        assertThat(service.enforce(subjects, "/api/anything", "GET")).isFalse();
    }
}
