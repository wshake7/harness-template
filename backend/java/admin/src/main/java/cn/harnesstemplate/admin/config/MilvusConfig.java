package cn.harnesstemplate.admin.config;

import cn.harnesstemplate.admin.infrastructure.milvus.MilvusVectorStore;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class MilvusConfig {

    private static final Logger log = LoggerFactory.getLogger(MilvusConfig.class);

    private final AdminProperties adminProperties;

    public MilvusConfig(AdminProperties adminProperties) {
        this.adminProperties = adminProperties;
    }

    @Bean
    public MilvusVectorStore milvusVectorStore() {
        var store = new MilvusVectorStore(adminProperties.getMilvus());
        log.info("Milvus: {}", store.isEnabled() ? "enabled" : "disabled");
        return store;
    }
}
