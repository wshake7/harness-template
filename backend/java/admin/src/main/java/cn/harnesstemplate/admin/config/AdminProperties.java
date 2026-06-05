package cn.harnesstemplate.admin.config;

import lombok.Data;
import org.springframework.boot.context.properties.ConfigurationProperties;

@Data
@ConfigurationProperties(prefix = "admin")
public class AdminProperties {

    /** API prefix, default /api */
    private String restPrefix = "/api";

    /** Default language for LanguageInterceptor */
    private String defaultLanguage = "zh-CN";

    /** Auth configuration */
    private AuthConfig auth = new AuthConfig();

    /** Temporal configuration */
    private TemporalConfig temporal = new TemporalConfig();

    /** Storage configuration */
    private StorageConfig storage = new StorageConfig();

    /** Milvus configuration */
    private MilvusConfig milvus = new MilvusConfig();

    /** AI configuration */
    private AiConfig ai = new AiConfig();

    @Data
    public static class AuthConfig {
        /** Token cookie name, must match frontend */
        private String tokenName = "token";
    }

    @Data
    public static class TemporalConfig {
        private boolean enabled = false;
        private String hostPort = "127.0.0.1:7233";
        private String namespace = "default";
        private String identity;
        private String taskQueue = "admin-java";
        private boolean workerEnabled = false;
    }

    @Data
    public static class StorageConfig {
        private boolean enabled = false;
        private String engine = "minio";
        private long maxUploadBytes = 10485760;
        private String objectKeyPrefix = "uploads";
        private int presignedExpiresSeconds = 3600;
        private MinioConfig minio = new MinioConfig();
    }

    @Data
    public static class MinioConfig {
        private String endpoint;
        private String accessKeyId;
        private String secretAccessKey;
        private String bucket = "admin-files";
        private String region = "us-east-1";
        private boolean useSsl = false;
        private boolean autoCreateBucket = true;
    }

    @Data
    public static class MilvusConfig {
        private boolean enabled = true;
        private String address = "127.0.0.1:19530";
        private String apiKey;
        private String dbName;
    }

    @Data
    public static class AiConfig {
        private EmbeddingConfig embedding = new EmbeddingConfig();
        private KnowledgeConfig knowledge = new KnowledgeConfig();
        private MemoryConfig memory = new MemoryConfig();
        private String mcpUrl;
    }

    @Data
    public static class EmbeddingConfig {
        private String provider = "ark";
        private String apiKey;
        private String model = "text-embedding-v3";
        private int dimensions = 2048;
    }

    @Data
    public static class KnowledgeConfig {
        private String collection = "biz";
        private String collectionDesc = "Knowledge documents for admin AI workflows";
        private int idMaxLength = 255;
        private int contentMaxLength = 8192;
        private String indexType = "auto";
        private String indexMetricType = "COSINE";
        private int loadTimeoutSeconds = 60;
        private int flushTimeoutSeconds = 30;
    }

    @Data
    public static class MemoryConfig {
        private String type = "memory";
        private int maxWindowSize = 6;
        private int ttlSeconds = 3600;
    }
}
