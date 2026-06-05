package cn.harnesstemplate.admin.config;

import cn.harnesstemplate.admin.infrastructure.temporal.TemporalClientFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class TemporalConfig {

    private static final Logger log = LoggerFactory.getLogger(TemporalConfig.class);

    private final AdminProperties adminProperties;

    public TemporalConfig(AdminProperties adminProperties) {
        this.adminProperties = adminProperties;
    }

    @Bean
    public TemporalClientFactory temporalClientFactory() {
        var factory = new TemporalClientFactory(adminProperties.getTemporal());
        if (adminProperties.getTemporal().isEnabled()) {
            log.info("Temporal enabled — starting client");
            factory.start();
        } else {
            log.info("Temporal disabled");
        }
        return factory;
    }
}
