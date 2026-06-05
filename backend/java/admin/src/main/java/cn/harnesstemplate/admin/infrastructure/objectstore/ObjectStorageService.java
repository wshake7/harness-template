package cn.harnesstemplate.admin.infrastructure.objectstore;

import cn.harnesstemplate.admin.config.AdminProperties;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;

public class ObjectStorageService {

    private static final Logger log = LoggerFactory.getLogger(ObjectStorageService.class);

    private final AdminProperties.StorageConfig config;

    public ObjectStorageService(AdminProperties.StorageConfig config) {
        this.config = config;
    }

    public boolean isEnabled() {
        return config.isEnabled();
    }

    public Map<String, Object> prepareUpload(String objectKey, String contentType, long size) {
        if (!isEnabled()) {
            log.debug("Storage disabled — skipping prepareUpload for {}", objectKey);
            return Map.of("presignedUrl", "", "objectKey", objectKey);
        }
        // TODO: Implement MinIO presigned URL generation
        log.info("Storage enabled — would generate presigned URL for {}", objectKey);
        return Map.of("presignedUrl", "http://minio:9000/" + config.getMinio().getBucket() + "/" + objectKey,
                "objectKey", objectKey);
    }

    public boolean verifyObjectExists(String objectKey) {
        if (!isEnabled()) return true; // assume exists when storage is disabled
        // TODO: Check MinIO for object
        return true;
    }

    public void deleteObject(String objectKey) {
        if (!isEnabled()) return;
        // TODO: Delete from MinIO
        log.info("Would delete object: {}", objectKey);
    }

    public String getBucket() {
        return config.getMinio().getBucket();
    }

    public long getMaxUploadBytes() {
        return config.getMaxUploadBytes();
    }
}
