package cn.harnesstemplate.admin.web.controller;

import cn.harnesstemplate.admin.application.service.knowledge.KnowledgeCollectionQueryService;
import cn.harnesstemplate.admin.domain.entity.KnowledgeCollection;
import cn.harnesstemplate.admin.web.dto.PageResult;
import cn.harnesstemplate.admin.web.dto.R;
import org.springframework.web.bind.annotation.*;

import java.util.*;

@RestController
@RequestMapping("/api/knowledge/collection")
public class KnowledgeCollectionController {

    private final KnowledgeCollectionQueryService service;

    public KnowledgeCollectionController(KnowledgeCollectionQueryService service) {
        this.service = service;
    }

    @PostMapping("/list")
    public R<PageResult<Map<String, Object>>> list() {
        List<KnowledgeCollection> all = service.listAll();
        List<Map<String, Object>> items = new ArrayList<>();
        for (KnowledgeCollection c : all) items.add(toMap(c));
        return R.ok(PageResult.of(items, items.size()));
    }

    @PostMapping("/detail")
    public R<Map<String, Object>> detail(@RequestBody Map<String, Long> req) {
        KnowledgeCollection c = service.findById(req.get("id"));
        if (c == null) return R.fail("知识库不存在");
        return R.ok(toMap(c));
    }

    @PostMapping("/create")
    public R<Void> create(@RequestBody Map<String, Object> req) {
        KnowledgeCollection c = new KnowledgeCollection();
        c.setCollectionName((String) req.get("collectionName"));
        c.setDisplayName((String) req.get("displayName"));
        c.setMetricType((String) req.getOrDefault("metricType", "COSINE"));
        c.setIndexType((String) req.getOrDefault("indexType", "auto"));
        c.setIsEnabled(true);
        c.setRemark((String) req.getOrDefault("remark", ""));
        service.save(c);
        return R.ok();
    }

    @PostMapping("/update")
    public R<Void> update(@RequestBody Map<String, Object> req) {
        Long id = ((Number) req.get("id")).longValue();
        KnowledgeCollection c = service.findById(id);
        if (c == null) return R.fail("知识库不存在");
        if (req.containsKey("collectionName")) c.setCollectionName((String) req.get("collectionName"));
        if (req.containsKey("displayName")) c.setDisplayName((String) req.get("displayName"));
        if (req.containsKey("metricType")) c.setMetricType((String) req.get("metricType"));
        if (req.containsKey("indexType")) c.setIndexType((String) req.get("indexType"));
        if (req.containsKey("isEnabled")) c.setIsEnabled((Boolean) req.get("isEnabled"));
        if (req.containsKey("remark")) c.setRemark((String) req.get("remark"));
        service.update(c);
        return R.ok();
    }

    @PostMapping("/del")
    public R<Void> del(@RequestBody Map<String, Long> req) {
        service.delete(req.get("id"));
        return R.ok();
    }

    private Map<String, Object> toMap(KnowledgeCollection c) {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("id", c.getId());
        m.put("collectionName", c.getCollectionName());
        m.put("displayName", c.getDisplayName());
        m.put("metricType", c.getMetricType());
        m.put("indexType", c.getIndexType());
        m.put("isEnabled", c.getIsEnabled());
        m.put("remark", c.getRemark());
        m.put("createdAt", c.getCreatedAt());
        m.put("updatedAt", c.getUpdatedAt());
        return m;
    }
}
