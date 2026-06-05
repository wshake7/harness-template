package cn.harnesstemplate.admin.web.controller;

import cn.harnesstemplate.admin.application.service.knowledge.KnowledgeDocumentQueryService;
import cn.harnesstemplate.admin.application.service.storage.FileAssetQueryService;
import cn.harnesstemplate.admin.domain.entity.FileAsset;
import cn.harnesstemplate.admin.domain.entity.KnowledgeDocument;
import cn.harnesstemplate.admin.web.dto.PageResult;
import cn.harnesstemplate.admin.web.dto.R;
import org.springframework.web.bind.annotation.*;

import java.util.*;

@RestController
@RequestMapping("/api/knowledge/document")
public class KnowledgeDocumentController {

    private final KnowledgeDocumentQueryService service;
    private final FileAssetQueryService fileService;

    public KnowledgeDocumentController(KnowledgeDocumentQueryService service, FileAssetQueryService fileService) {
        this.service = service;
        this.fileService = fileService;
    }

    @PostMapping("/list")
    public R<PageResult<Map<String, Object>>> list(@RequestBody Map<String, Object> req) {
        Long collectionId = ((Number) req.get("collectionId")).longValue();
        List<KnowledgeDocument> docs = service.findByCollectionId(collectionId);
        List<Map<String, Object>> items = new ArrayList<>();
        for (KnowledgeDocument d : docs) items.add(toMap(d));
        return R.ok(PageResult.of(items, items.size()));
    }

    @PostMapping("/detail")
    public R<Map<String, Object>> detail(@RequestBody Map<String, Long> req) {
        KnowledgeDocument d = service.findById(req.get("id"));
        if (d == null) return R.fail("文档不存在");
        return R.ok(toMap(d));
    }

    @PostMapping("/create")
    public R<Void> create(@RequestBody Map<String, Object> req) {
        KnowledgeDocument d = new KnowledgeDocument();
        d.setCollectionId(((Number) req.get("collectionId")).longValue());
        d.setTitle((String) req.get("title"));
        d.setContent((String) req.getOrDefault("content", ""));
        d.setContentType((String) req.getOrDefault("contentType", "text"));
        d.setSource((String) req.getOrDefault("source", "manual"));
        d.setMetadata(req.containsKey("metadata") ? req.get("metadata").toString() : "{}");
        d.setDocumentId("doc-" + System.currentTimeMillis());
        d.setIsEnabled(true);
        d.setVectorStatus("pending");
        d.setRemark((String) req.getOrDefault("remark", ""));
        service.save(d);
        return R.ok();
    }

    @PostMapping("/update")
    public R<Void> update(@RequestBody Map<String, Object> req) {
        Long id = ((Number) req.get("id")).longValue();
        KnowledgeDocument d = service.findById(id);
        if (d == null) return R.fail("文档不存在");
        if (req.containsKey("title")) d.setTitle((String) req.get("title"));
        if (req.containsKey("content")) d.setContent((String) req.get("content"));
        if (req.containsKey("contentType")) d.setContentType((String) req.get("contentType"));
        if (req.containsKey("isEnabled")) d.setIsEnabled((Boolean) req.get("isEnabled"));
        if (req.containsKey("metadata")) d.setMetadata(req.get("metadata").toString());
        if (req.containsKey("remark")) d.setRemark((String) req.get("remark"));
        service.update(d);
        return R.ok();
    }

    @PostMapping("/del")
    public R<Void> del(@RequestBody Map<String, Long> req) {
        service.delete(req.get("id"));
        return R.ok();
    }

    @PostMapping("/importFile")
    public R<Map<String, Object>> importFile(@RequestBody Map<String, Object> req) {
        Long collectionId = ((Number) req.get("collectionId")).longValue();
        Long fileId = ((Number) req.get("fileId")).longValue();

        FileAsset file = fileService.findById(fileId);
        if (file == null || !"active".equals(file.getStatus())) {
            return R.fail("文件不可用");
        }

        KnowledgeDocument d = new KnowledgeDocument();
        d.setCollectionId(collectionId);
        d.setDocumentId("doc-" + file.getId() + "-" + System.currentTimeMillis());
        d.setTitle(file.getOriginalName());
        d.setContent("");
        d.setContentType(file.getContentType());
        d.setSource("file-import");
        d.setMetadata("{\"fileAssetId\": " + fileId + "}");
        d.setIsEnabled(true);
        d.setVectorStatus("pending");
        d.setRemark("Imported from file: " + file.getOriginalName());
        service.save(d);

        Map<String, Object> result = new LinkedHashMap<>();
        result.put("id", d.getId());
        result.put("vectorStatus", d.getVectorStatus());
        return R.ok(result);
    }

    private Map<String, Object> toMap(KnowledgeDocument d) {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("id", d.getId());
        m.put("collectionId", d.getCollectionId());
        m.put("documentId", d.getDocumentId());
        m.put("title", d.getTitle());
        m.put("contentType", d.getContentType());
        m.put("source", d.getSource());
        m.put("chunkIndex", d.getChunkIndex());
        m.put("totalChunks", d.getTotalChunks());
        m.put("vectorStatus", d.getVectorStatus());
        m.put("vectorId", d.getVectorId());
        m.put("indexingError", d.getIndexingError());
        m.put("lastIndexedAt", d.getLastIndexedAt());
        m.put("isEnabled", d.getIsEnabled());
        m.put("remark", d.getRemark());
        m.put("metadata", d.getMetadata());
        m.put("createdAt", d.getCreatedAt());
        m.put("updatedAt", d.getUpdatedAt());
        return m;
    }
}
