package cn.harnesstemplate.admin.web.controller;

import cn.harnesstemplate.admin.application.service.storage.FileAssetQueryService;
import cn.harnesstemplate.admin.domain.entity.FileAsset;
import cn.harnesstemplate.admin.infrastructure.objectstore.ObjectStorageService;
import cn.harnesstemplate.admin.web.dto.PageResult;
import cn.harnesstemplate.admin.web.dto.R;
import org.springframework.web.bind.annotation.*;

import java.util.*;

@RestController
@RequestMapping("/api/storage/file")
public class StorageFileController {

    private final FileAssetQueryService service;
    private final ObjectStorageService storage;

    public StorageFileController(FileAssetQueryService service, ObjectStorageService storage) {
        this.service = service;
        this.storage = storage;
    }

    @PostMapping("/list")
    public R<PageResult<Map<String, Object>>> list(@RequestBody Map<String, Object> req) {
        int page = req.containsKey("page") ? ((Number) req.get("page")).intValue() : 1;
        int pageSize = req.containsKey("pageSize") ? ((Number) req.get("pageSize")).intValue() : 10;
        List<FileAsset> files = service.page(page, pageSize);
        long total = service.count();
        List<Map<String, Object>> items = new ArrayList<>();
        for (FileAsset f : files) items.add(toMap(f));
        return R.ok(PageResult.of(items, total));
    }

    @PostMapping("/detail")
    public R<Map<String, Object>> detail(@RequestBody Map<String, Long> req) {
        FileAsset f = service.findById(req.get("id"));
        if (f == null) return R.fail("文件不存在");
        return R.ok(toMap(f));
    }

    @PostMapping("/prepareUpload")
    public R<Map<String, Object>> prepareUpload(@RequestBody Map<String, Object> req) {
        String originalName = (String) req.get("originalName");
        String contentType = (String) req.getOrDefault("contentType", "application/octet-stream");
        long size = req.containsKey("size") ? ((Number) req.get("size")).longValue() : 0;
        String bizType = (String) req.getOrDefault("bizType", "");
        String bizId = (String) req.getOrDefault("bizId", "");

        if (size > storage.getMaxUploadBytes()) {
            return R.fail("文件大小超过限制");
        }

        String objectKey = storage.getBucket() + "/" + UUID.randomUUID() + "-" + originalName;
        var presigned = storage.prepareUpload(objectKey, contentType, size);

        FileAsset asset = new FileAsset();
        asset.setEngine("minio");
        asset.setBucket(storage.getBucket());
        asset.setObjectKey(objectKey);
        asset.setOriginalName(originalName);
        asset.setContentType(contentType);
        asset.setExtension(getExtension(originalName));
        asset.setSize(size);
        asset.setBizType(bizType);
        asset.setBizId(bizId);
        asset.setStatus("pending_upload");
        service.save(asset);

        Map<String, Object> result = new LinkedHashMap<>();
        result.put("id", asset.getId());
        result.put("presignedUrl", presigned.get("presignedUrl"));
        result.put("objectKey", objectKey);
        return R.ok(result);
    }

    @PostMapping("/completeUpload")
    public R<Void> completeUpload(@RequestBody Map<String, Long> req) {
        Long id = req.get("id");
        FileAsset f = service.findById(id);
        if (f == null) return R.fail("文件不存在");
        if (!"pending_upload".equals(f.getStatus())) return R.fail("文件状态不正确");

        if (!storage.verifyObjectExists(f.getObjectKey())) {
            return R.fail("对象存储验证失败");
        }
        service.updateStatus(id, "active");
        return R.ok();
    }

    @PostMapping("/del")
    public R<Void> del(@RequestBody Map<String, Long> req) {
        Long id = req.get("id");
        FileAsset f = service.findById(id);
        if (f == null) return R.fail("文件不存在");

        storage.deleteObject(f.getObjectKey());
        service.delete(id);
        return R.ok();
    }

    private Map<String, Object> toMap(FileAsset f) {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("id", f.getId());
        m.put("originalName", f.getOriginalName());
        m.put("contentType", f.getContentType());
        m.put("extension", f.getExtension());
        m.put("size", f.getSize());
        m.put("sha256", f.getSha256());
        m.put("bizType", f.getBizType());
        m.put("bizId", f.getBizId());
        m.put("status", f.getStatus());
        m.put("engine", f.getEngine());
        m.put("bucket", f.getBucket());
        m.put("objectKey", f.getObjectKey());
        m.put("createdAt", f.getCreatedAt());
        m.put("updatedAt", f.getUpdatedAt());
        return m;
    }

    private String getExtension(String name) {
        if (name == null || !name.contains(".")) return "";
        return name.substring(name.lastIndexOf('.'));
    }
}
