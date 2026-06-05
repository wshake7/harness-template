package cn.harnesstemplate.admin.infrastructure.milvus;

import cn.harnesstemplate.admin.config.AdminProperties;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.Map;

public class MilvusVectorStore {

    private static final Logger log = LoggerFactory.getLogger(MilvusVectorStore.class);

    private final AdminProperties.MilvusConfig config;

    public MilvusVectorStore(AdminProperties.MilvusConfig config) {
        this.config = config;
    }

    public boolean isEnabled() {
        return config.isEnabled();
    }

    public void createCollection(String collectionName, int dimension, String metricType) {
        if (!isEnabled()) {
            log.debug("Milvus disabled — skipping collection creation: {}", collectionName);
            return;
        }
        log.info("Would create Milvus collection: {} (dim={}, metric={})", collectionName, dimension, metricType);
    }

    public void insert(String collectionName, List<Float> vector, Map<String, Object> metadata) {
        if (!isEnabled()) return;
        log.debug("Would insert vector into collection {}", collectionName);
    }

    public void deleteByFilter(String collectionName, String filter) {
        if (!isEnabled()) return;
        log.debug("Would delete by filter '{}' from collection {}", filter, collectionName);
    }

    public List<Map<String, Object>> search(String collectionName, List<Float> queryVector, int topK) {
        if (!isEnabled()) return List.of();
        log.debug("Would search collection {} with topK={}", collectionName, topK);
        return List.of();
    }
}
