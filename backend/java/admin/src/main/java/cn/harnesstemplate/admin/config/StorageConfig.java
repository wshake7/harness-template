package cn.harnesstemplate.admin.config;

import cn.harnesstemplate.admin.infrastructure.objectstore.ObjectStorageService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class StorageConfig {

    private static final Logger log = LoggerFactory.getLogger(StorageConfig.class);

    private final AdminProperties adminProperties;

    public StorageConfig(AdminProperties adminProperties) {
        this.adminProperties = adminProperties;
    }

    @Bean
    public ObjectStorageService objectStorageService() {
        var svc = new ObjectStorageService(adminProperties.getStorage());
        log.info("Object storage: {}", svc.isEnabled() ? "enabled" : "disabled");
        return svc;
    }
}
