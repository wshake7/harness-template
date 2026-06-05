package cn.harnesstemplate.admin.ai;

import cn.harnesstemplate.admin.infrastructure.milvus.MilvusVectorStore;
import dev.langchain4j.model.chat.ChatModel;
import dev.langchain4j.service.AiServices;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.ObjectProvider;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

@Service
public class RagChatService {

    private static final Logger log = LoggerFactory.getLogger(RagChatService.class);

    private final EmbeddingService embeddingService;
    private final MilvusVectorStore vectorStore;
    private final ObjectProvider<ChatModel> chatModelProvider;

    public RagChatService(EmbeddingService embeddingService,
                          MilvusVectorStore vectorStore,
                          ObjectProvider<ChatModel> chatModelProvider) {
        this.embeddingService = embeddingService;
        this.vectorStore = vectorStore;
        this.chatModelProvider = chatModelProvider;
    }

    public RagResponse chat(RagRequest request) {
        log.info("RAG chat: collection={}, prompt length={}", request.getCollectionName(), request.getPrompt().length());

        // 1. Embed the query
        List<Float> queryVector = embeddingService.embed(request.getPrompt());

        // 2. Search for relevant documents
        List<Map<String, Object>> retrieved = List.of();
        if (!queryVector.isEmpty()) {
            retrieved = vectorStore.search(request.getCollectionName(), queryVector, request.getTopK());
        }

        // 3. Build context from retrieved documents
        StringBuilder context = new StringBuilder();
        for (Map<String, Object> doc : retrieved) {
            context.append(doc.getOrDefault("content", "")).append("\n");
        }

        // 4. Build augmented prompt
        String augmentedPrompt = "基于以下参考文档回答问题:\n\n" +
                "参考文档:\n" + context + "\n\n" +
                "用户问题: " + request.getPrompt();

        RagResponse response = new RagResponse();
        response.setPrompt(request.getPrompt());
        response.setRetrievedCount(retrieved.size());
        response.setAugmentedPrompt(augmentedPrompt);
        response.setAnswer(generateAnswer(augmentedPrompt));
        response.setSources(toSourceList(retrieved));
        return response;
    }

    private String generateAnswer(String augmentedPrompt) {
        ChatModel chatModel = chatModelProvider.getIfAvailable();
        if (chatModel == null) {
            log.warn("LangChain4j ChatModel bean not configured — returning augmented prompt fallback");
            return "LangChain4j ChatModel is not configured";
        }
        Assistant assistant = AiServices.builder(Assistant.class)
                .chatModel(chatModel)
                .build();
        return assistant.chat(augmentedPrompt);
    }

    private List<RagSource> toSourceList(List<Map<String, Object>> docs) {
        List<RagSource> sources = new ArrayList<>();
        for (Map<String, Object> doc : docs) {
            RagSource src = new RagSource();
            src.setTitle((String) doc.getOrDefault("title", ""));
            src.setContent((String) doc.getOrDefault("content", ""));
            sources.add(src);
        }
        return sources;
    }

    private interface Assistant {
        String chat(String userMessage);
    }
}
