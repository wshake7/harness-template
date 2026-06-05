package cn.harnesstemplate.admin.web.controller;

import cn.harnesstemplate.admin.application.service.system.SysResourceApiQueryService;
import cn.harnesstemplate.admin.domain.entity.SysResourceApi;
import cn.harnesstemplate.admin.web.dto.PageResult;
import cn.harnesstemplate.admin.web.dto.R;
import org.springframework.web.bind.annotation.*;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/sys/resource/api")
public class ApiController {

    private final SysResourceApiQueryService apiService;

    public ApiController(SysResourceApiQueryService apiService) {
        this.apiService = apiService;
    }

    @PostMapping("/list")
    public R<PageResult<Map<String, Object>>> list(@RequestBody Map<String, Object> req) {
        List<SysResourceApi> all = apiService.listAll();
        List<Map<String, Object>> items = new ArrayList<>();
        for (SysResourceApi a : all) {
            Map<String, Object> item = new HashMap<>();
            item.put("id", a.getId());
            item.put("module", a.getModule());
            item.put("path", a.getPath());
            item.put("method", a.getMethod());
            item.put("sortOrder", a.getSortOrder());
            item.put("isEnabled", a.getIsEnabled());
            item.put("remark", a.getRemark());
            item.put("createdAt", a.getCreatedAt());
            item.put("updatedAt", a.getUpdatedAt());
            items.add(item);
        }
        return R.ok(PageResult.of(items, items.size()));
    }

    @PostMapping("/create")
    public R<Void> create(@RequestBody Map<String, Object> req) {
        SysResourceApi api = new SysResourceApi();
        api.setModule((String) req.getOrDefault("module", ""));
        api.setPath((String) req.get("path"));
        api.setMethod(((String) req.get("method")).toUpperCase());
        api.setSortOrder(req.containsKey("sortOrder") ? ((Number) req.get("sortOrder")).intValue() : 0);
        api.setIsEnabled(req.containsKey("isEnabled") ? (Boolean) req.get("isEnabled") : true);
        api.setRemark((String) req.getOrDefault("remark", ""));
        apiService.save(api);
        return R.ok();
    }

    @PostMapping("/update")
    public R<Void> update(@RequestBody Map<String, Object> req) {
        Long id = ((Number) req.get("id")).longValue();
        SysResourceApi api = apiService.findById(id);
        if (api == null) {
            return R.fail("API资源不存在");
        }
        if (req.containsKey("module")) api.setModule((String) req.get("module"));
        if (req.containsKey("path")) api.setPath((String) req.get("path"));
        if (req.containsKey("method")) api.setMethod(((String) req.get("method")).toUpperCase());
        if (req.containsKey("sortOrder")) api.setSortOrder(((Number) req.get("sortOrder")).intValue());
        if (req.containsKey("isEnabled")) api.setIsEnabled((Boolean) req.get("isEnabled"));
        if (req.containsKey("remark")) api.setRemark((String) req.get("remark"));
        apiService.update(api);
        return R.ok();
    }

    @PostMapping("/del")
    public R<Void> del(@RequestBody Map<String, List<Long>> req) {
        List<Long> ids = req.get("ids");
        if (ids == null || ids.isEmpty()) {
            return R.fail("ids不能为空");
        }
        for (Long id : ids) {
            apiService.delete(id);
        }
        return R.ok();
    }
}
