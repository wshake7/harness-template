package cn.harnesstemplate.admin.infrastructure.temporal;

import cn.harnesstemplate.admin.config.AdminProperties;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

class TemporalClientFactoryTests {

    @Test
    void shouldNotBeEnabledByDefault() {
        var config = new AdminProperties.TemporalConfig();
        var factory = new TemporalClientFactory(config);
        assertThat(factory.isEnabled()).isFalse();
        assertThat(factory.getScheduleClient()).isNull();
        assertThat(factory.getWorkflowClient()).isNull();
    }

    @Test
    void shouldSetCorrectDefaults() {
        var config = new AdminProperties.TemporalConfig();
        assertThat(config.getHostPort()).isEqualTo("127.0.0.1:7233");
        assertThat(config.getNamespace()).isEqualTo("default");
        assertThat(config.getTaskQueue()).isEqualTo("admin-java");
    }

    @Test
    void shouldNotStartWhenDisabled() {
        var config = new AdminProperties.TemporalConfig();
        config.setEnabled(false);
        var factory = new TemporalClientFactory(config);
        factory.start(); // should not throw
        assertThat(factory.getWorkflowClient()).isNull();
    }
}
