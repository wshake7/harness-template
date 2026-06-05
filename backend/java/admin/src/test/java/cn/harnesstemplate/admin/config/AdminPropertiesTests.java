package cn.harnesstemplate.admin.config;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.ActiveProfiles;

import static org.assertj.core.api.Assertions.assertThat;

@SpringBootTest
@ActiveProfiles("test")
class AdminPropertiesTests {

    @Autowired
    private AdminProperties properties;

    @Test
    void shouldBindDefaultRestPrefix() {
        assertThat(properties.getRestPrefix()).isEqualTo("/api");
    }

    @Test
    void shouldBindServerPort() {
        assertThat(properties.getDefaultLanguage()).isEqualTo("zh-CN");
    }

    @Test
    void shouldBindTokenName() {
        assertThat(properties.getAuth().getTokenName()).isEqualTo("token");
    }

    @Test
    void shouldBindTemporalDefaults() {
        assertThat(properties.getTemporal().getTaskQueue()).isEqualTo("admin-java");
    }
}
