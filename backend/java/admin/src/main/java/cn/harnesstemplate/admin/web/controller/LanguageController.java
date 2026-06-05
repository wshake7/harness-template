package cn.harnesstemplate.admin.web.controller;

import cn.harnesstemplate.admin.application.service.system.SysLanguageEntryQueryService;
import cn.harnesstemplate.admin.application.service.system.SysLanguageTypeQueryService;
import cn.harnesstemplate.admin.domain.entity.SysLanguageEntry;
import cn.harnesstemplate.admin.domain.entity.SysLanguageType;
import cn.harnesstemplate.admin.web.dto.PageResult;
import cn.harnesstemplate.admin.web.dto.R;
import org.springframework.web.bind.annotation.*;

import java.util.*;

@RestController
@RequestMapping("/api/sys/language")
public class LanguageController {

    private final SysLanguageTypeQueryService typeService;
    private final SysLanguageEntryQueryService entryService;

    public LanguageController(SysLanguageTypeQueryService typeService, SysLanguageEntryQueryService entryService) {
        this.typeService = typeService;
        this.entryService = entryService;
    }

    @PostMapping("/type/list")
    public R<PageResult<Map<String, Object>>> typeList() {
        List<SysLanguageType> types = typeService.listAll();
        List<Map<String, Object>> items = new ArrayList<>();
        for (SysLanguageType t : types) {
            Map<String, Object> item = toTypeMap(t);
            List<SysLanguageEntry> entries = entryService.findByTypeId(t.getId());
            List<Map<String, Object>> entryMaps = new ArrayList<>();
            for (SysLanguageEntry e : entries) {
                entryMaps.add(toEntryMap(e));
            }
            item.put("entries", entryMaps);
            items.add(item);
        }
        return R.ok(PageResult.of(items, items.size()));
    }

    @PostMapping("/type/create")
    public R<Void> typeCreate(@RequestBody Map<String, Object> req) {
        SysLanguageType t = new SysLanguageType();
        t.setTypeCode((String) req.get("typeCode"));
        t.setTypeName((String) req.get("typeName"));
        t.setIsDefault(req.containsKey("isDefault") ? (Boolean) req.get("isDefault") : false);
        t.setIsEnabled(true);
        t.setSortOrder(req.containsKey("sortOrder") ? ((Number) req.get("sortOrder")).intValue() : 0);
        typeService.save(t);
        return R.ok();
    }

    @PostMapping("/type/update")
    public R<Void> typeUpdate(@RequestBody Map<String, Object> req) {
        Long id = ((Number) req.get("id")).longValue();
        SysLanguageType t = typeService.findById(id);
        if (t == null) return R.fail("语言类型不存在");
        if (req.containsKey("typeCode")) t.setTypeCode((String) req.get("typeCode"));
        if (req.containsKey("typeName")) t.setTypeName((String) req.get("typeName"));
        if (req.containsKey("isDefault")) t.setIsDefault((Boolean) req.get("isDefault"));
        if (req.containsKey("isEnabled")) t.setIsEnabled((Boolean) req.get("isEnabled"));
        if (req.containsKey("sortOrder")) t.setSortOrder(((Number) req.get("sortOrder")).intValue());
        typeService.update(t);
        return R.ok();
    }

    @PostMapping("/type/del")
    public R<Void> typeDel(@RequestBody Map<String, List<Long>> req) {
        for (Long id : req.get("ids")) typeService.delete(id);
        return R.ok();
    }

    @PostMapping("/entry/list")
    public R<PageResult<Map<String, Object>>> entryList(@RequestBody Map<String, Object> req) {
        Long typeId = ((Number) req.get("typeId")).longValue();
        List<SysLanguageEntry> entries = entryService.findByTypeId(typeId);
        List<Map<String, Object>> items = new ArrayList<>();
        for (SysLanguageEntry e : entries) items.add(toEntryMap(e));
        return R.ok(PageResult.of(items, items.size()));
    }

    @PostMapping("/entry/create")
    public R<Void> entryCreate(@RequestBody Map<String, Object> req) {
        SysLanguageEntry e = new SysLanguageEntry();
        e.setSysLanguageTypeId(((Number) req.get("typeId")).longValue());
        e.setEntryCode((String) req.get("entryCode"));
        e.setEntryValue((String) req.get("entryValue"));
        e.setIsEnabled(req.containsKey("isEnabled") ? (Boolean) req.get("isEnabled") : true);
        e.setSortOrder(req.containsKey("sortOrder") ? ((Number) req.get("sortOrder")).intValue() : 0);
        e.setRemark((String) req.getOrDefault("remark", ""));
        entryService.save(e);
        return R.ok();
    }

    @PostMapping("/entry/update")
    public R<Void> entryUpdate(@RequestBody Map<String, Object> req) {
        Long id = ((Number) req.get("id")).longValue();
        SysLanguageEntry e = entryService.findById(id);
        if (e == null) return R.fail("语言条目不存在");
        if (req.containsKey("entryCode")) e.setEntryCode((String) req.get("entryCode"));
        if (req.containsKey("entryValue")) e.setEntryValue((String) req.get("entryValue"));
        if (req.containsKey("isEnabled")) e.setIsEnabled((Boolean) req.get("isEnabled"));
        if (req.containsKey("sortOrder")) e.setSortOrder(((Number) req.get("sortOrder")).intValue());
        if (req.containsKey("remark")) e.setRemark((String) req.get("remark"));
        entryService.update(e);
        return R.ok();
    }

    @PostMapping("/entry/del")
    public R<Void> entryDel(@RequestBody Map<String, List<Long>> req) {
        for (Long id : req.get("ids")) entryService.delete(id);
        return R.ok();
    }

    private Map<String, Object> toTypeMap(SysLanguageType t) {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("id", t.getId());
        m.put("typeCode", t.getTypeCode());
        m.put("typeName", t.getTypeName());
        m.put("isDefault", t.getIsDefault());
        m.put("isEnabled", t.getIsEnabled());
        m.put("sortOrder", t.getSortOrder());
        m.put("createdAt", t.getCreatedAt());
        m.put("updatedAt", t.getUpdatedAt());
        return m;
    }

    private Map<String, Object> toEntryMap(SysLanguageEntry e) {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("id", e.getId());
        m.put("typeId", e.getSysLanguageTypeId());
        m.put("entryCode", e.getEntryCode());
        m.put("entryValue", e.getEntryValue());
        m.put("isEnabled", e.getIsEnabled());
        m.put("sortOrder", e.getSortOrder());
        m.put("remark", e.getRemark());
        m.put("createdAt", e.getCreatedAt());
        m.put("updatedAt", e.getUpdatedAt());
        return m;
    }
}
