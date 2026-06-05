package cn.harnesstemplate.admin.ai;

import cn.harnesstemplate.admin.config.AdminProperties;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class EmbeddingService {

    private static final Logger log = LoggerFactory.getLogger(EmbeddingService.class);

    private final AdminProperties.AiConfig aiConfig;

    public EmbeddingService(AdminProperties adminProperties) {
        this.aiConfig = adminProperties.getAi();
    }

    public boolean isAvailable() {
        var emb = aiConfig.getEmbedding();
        return emb.getApiKey() != null && !emb.getApiKey().isEmpty();
    }

    public List<Float> embed(String text) {
        if (!isAvailable()) {
            log.warn("Embedding API key not configured — returning empty vector");
            return List.of();
        }
        // TODO: Implement actual embedding API call (Ark/Doubao)
        log.debug("Would embed text of length {} via provider {} model {}",
                text.length(), aiConfig.getEmbedding().getProvider(), aiConfig.getEmbedding().getModel());
        return List.of();
    }

    public int getDimensions() {
        return aiConfig.getEmbedding().getDimensions();
    }
}
