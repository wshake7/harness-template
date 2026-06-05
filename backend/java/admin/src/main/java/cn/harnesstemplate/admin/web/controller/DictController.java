package cn.harnesstemplate.admin.web.controller;

import cn.harnesstemplate.admin.application.service.system.SysDictEntryQueryService;
import cn.harnesstemplate.admin.application.service.system.SysDictTypeQueryService;
import cn.harnesstemplate.admin.domain.entity.SysDictEntry;
import cn.harnesstemplate.admin.domain.entity.SysDictType;
import cn.harnesstemplate.admin.web.dto.PageResult;
import cn.harnesstemplate.admin.web.dto.R;
import org.springframework.web.bind.annotation.*;

import java.util.*;

@RestController
@RequestMapping("/api/sys/dict")
public class DictController {

    private final SysDictTypeQueryService typeService;
    private final SysDictEntryQueryService entryService;

    public DictController(SysDictTypeQueryService typeService, SysDictEntryQueryService entryService) {
        this.typeService = typeService;
        this.entryService = entryService;
    }

    @PostMapping("/type/list")
    public R<PageResult<Map<String, Object>>> typeList() {
        List<SysDictType> types = typeService.listAll();
        List<Map<String, Object>> items = new ArrayList<>();
        for (SysDictType t : types) {
            Map<String, Object> item = toTypeMap(t);
            List<SysDictEntry> entries = entryService.findByTypeId(t.getId());
            List<Map<String, Object>> entryMaps = new ArrayList<>();
            for (SysDictEntry e : entries) {
                entryMaps.add(toEntryMap(e));
            }
            item.put("entries", entryMaps);
            items.add(item);
        }
        return R.ok(PageResult.of(items, items.size()));
    }

    @PostMapping("/type/create")
    public R<Void> typeCreate(@RequestBody Map<String, Object> req) {
        SysDictType t = new SysDictType();
        t.setTypeCode((String) req.get("typeCode"));
        t.setTypeName((String) req.get("typeName"));
        t.setIsEnabled(req.containsKey("isEnabled") ? (Boolean) req.get("isEnabled") : true);
        t.setSortOrder(req.containsKey("sortOrder") ? ((Number) req.get("sortOrder")).intValue() : 0);
        t.setRemark((String) req.getOrDefault("remark", ""));
        typeService.save(t);
        return R.ok();
    }

    @PostMapping("/type/update")
    public R<Void> typeUpdate(@RequestBody Map<String, Object> req) {
        Long id = ((Number) req.get("id")).longValue();
        SysDictType t = typeService.findById(id);
        if (t == null) return R.fail("字典类型不存在");
        if (req.containsKey("typeCode")) t.setTypeCode((String) req.get("typeCode"));
        if (req.containsKey("typeName")) t.setTypeName((String) req.get("typeName"));
        if (req.containsKey("isEnabled")) t.setIsEnabled((Boolean) req.get("isEnabled"));
        if (req.containsKey("sortOrder")) t.setSortOrder(((Number) req.get("sortOrder")).intValue());
        if (req.containsKey("remark")) t.setRemark((String) req.get("remark"));
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
        List<SysDictEntry> entries = entryService.findByTypeId(typeId);
        List<Map<String, Object>> items = new ArrayList<>();
        for (SysDictEntry e : entries) items.add(toEntryMap(e));
        return R.ok(PageResult.of(items, items.size()));
    }

    @PostMapping("/entry/create")
    public R<Void> entryCreate(@RequestBody Map<String, Object> req) {
        SysDictEntry e = new SysDictEntry();
        e.setSysDictTypeId(((Number) req.get("typeId")).longValue());
        e.setEntryLabel((String) req.get("entryLabel"));
        e.setEntryValue((String) req.get("entryValue"));
        e.setLabelComponent((String) req.getOrDefault("labelComponent", ""));
        e.setLanguageCode((String) req.getOrDefault("languageCode", ""));
        e.setIsEnabled(req.containsKey("isEnabled") ? (Boolean) req.get("isEnabled") : true);
        e.setSortOrder(req.containsKey("sortOrder") ? ((Number) req.get("sortOrder")).intValue() : 0);
        e.setRemark((String) req.getOrDefault("remark", ""));
        entryService.save(e);
        return R.ok();
    }

    @PostMapping("/entry/update")
    public R<Void> entryUpdate(@RequestBody Map<String, Object> req) {
        Long id = ((Number) req.get("id")).longValue();
        SysDictEntry e = entryService.findById(id);
        if (e == null) return R.fail("字典条目不存在");
        if (req.containsKey("entryLabel")) e.setEntryLabel((String) req.get("entryLabel"));
        if (req.containsKey("entryValue")) e.setEntryValue((String) req.get("entryValue"));
        if (req.containsKey("labelComponent")) e.setLabelComponent((String) req.get("labelComponent"));
        if (req.containsKey("languageCode")) e.setLanguageCode((String) req.get("languageCode"));
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

    @PostMapping("/entry/match")
    public R<Map<String, List<Map<String, Object>>>> entryMatch(@RequestBody Map<String, List<String>> req) {
        List<String> codes = req.getOrDefault("codes", Collections.emptyList());
        Map<String, List<Map<String, Object>>> result = new LinkedHashMap<>();
        for (String code : codes) result.put(code, new ArrayList<>());
        if (codes.isEmpty()) return R.ok(result);

        List<SysDictType> types = typeService.findByCodes(codes);
        if (types.isEmpty()) return R.ok(result);

        Map<Long, String> codeByTypeId = new HashMap<>();
        List<Long> typeIds = new ArrayList<>();
        for (SysDictType t : types) {
            typeIds.add(t.getId());
            codeByTypeId.put(t.getId(), t.getTypeCode());
        }

        List<SysDictEntry> entries = entryService.findByTypeIds(typeIds);
        for (SysDictEntry e : entries) {
            String code = codeByTypeId.get(e.getSysDictTypeId());
            if (code == null) continue;
            result.get(code).add(toEntryMatchMap(e));
        }
        return R.ok(result);
    }

    private Map<String, Object> toEntryMatchMap(SysDictEntry e) {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("id", e.getId());
        m.put("labelComponent", e.getLabelComponent());
        m.put("entryLabel", e.getEntryLabel());
        m.put("entryValue", e.getEntryValue());
        return m;
    }

    private Map<String, Object> toTypeMap(SysDictType t) {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("id", t.getId());
        m.put("typeCode", t.getTypeCode());
        m.put("typeName", t.getTypeName());
        m.put("isEnabled", t.getIsEnabled());
        m.put("sortOrder", t.getSortOrder());
        m.put("remark", t.getRemark());
        m.put("createdAt", t.getCreatedAt());
        m.put("updatedAt", t.getUpdatedAt());
        return m;
    }

    private Map<String, Object> toEntryMap(SysDictEntry e) {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("id", e.getId());
        m.put("typeId", e.getSysDictTypeId());
        m.put("entryLabel", e.getEntryLabel());
        m.put("entryValue", e.getEntryValue());
        m.put("labelComponent", e.getLabelComponent());
        m.put("languageCode", e.getLanguageCode());
        m.put("isEnabled", e.getIsEnabled());
        m.put("sortOrder", e.getSortOrder());
        m.put("remark", e.getRemark());
        m.put("createdAt", e.getCreatedAt());
        m.put("updatedAt", e.getUpdatedAt());
        return m;
    }
}
