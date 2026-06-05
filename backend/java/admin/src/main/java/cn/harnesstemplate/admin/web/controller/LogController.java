package cn.harnesstemplate.admin.web.controller;

import cn.harnesstemplate.admin.application.service.system.SysApiLogQueryService;
import cn.harnesstemplate.admin.application.service.system.SysLoginLogQueryService;
import cn.harnesstemplate.admin.domain.entity.SysApiLog;
import cn.harnesstemplate.admin.domain.entity.SysLoginLog;
import cn.harnesstemplate.admin.web.dto.PageResult;
import cn.harnesstemplate.admin.web.dto.R;
import org.springframework.web.bind.annotation.*;

import java.util.*;

@RestController
@RequestMapping("/api/sys/log")
public class LogController {

    private final SysApiLogQueryService apiLogService;
    private final SysLoginLogQueryService loginLogService;

    public LogController(SysApiLogQueryService apiLogService, SysLoginLogQueryService loginLogService) {
        this.apiLogService = apiLogService;
        this.loginLogService = loginLogService;
    }

    @PostMapping("/api/list")
    public R<PageResult<Map<String, Object>>> apiLogList(@RequestBody Map<String, Object> req) {
        int page = req.containsKey("page") ? ((Number) req.get("page")).intValue() : 1;
        int pageSize = req.containsKey("pageSize") ? ((Number) req.get("pageSize")).intValue() : 10;
        String path = (String) req.getOrDefault("path", null);
        String method = (String) req.getOrDefault("method", null);

        List<SysApiLog> logs = apiLogService.page(page, pageSize, path, method);
        long total = apiLogService.count(path, method);
        List<Map<String, Object>> items = new ArrayList<>();
        for (SysApiLog l : logs) items.add(toApiLogMap(l));
        return R.ok(PageResult.of(items, total));
    }

    @PostMapping("/login/list")
    public R<PageResult<Map<String, Object>>> loginLogList(@RequestBody Map<String, Object> req) {
        int page = req.containsKey("page") ? ((Number) req.get("page")).intValue() : 1;
        int pageSize = req.containsKey("pageSize") ? ((Number) req.get("pageSize")).intValue() : 10;
        String username = (String) req.getOrDefault("username", null);

        List<SysLoginLog> logs = loginLogService.page(page, pageSize, username);
        long total = loginLogService.count(username);
        List<Map<String, Object>> items = new ArrayList<>();
        for (SysLoginLog l : logs) items.add(toLoginLogMap(l));
        return R.ok(PageResult.of(items, total));
    }

    private Map<String, Object> toApiLogMap(SysApiLog l) {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("id", l.getId());
        m.put("requestId", l.getRequestId());
        m.put("method", l.getMethod());
        m.put("module", l.getModule());
        m.put("path", l.getPath());
        m.put("statusCode", l.getStatusCode());
        m.put("success", l.getSuccess());
        m.put("costTime", l.getCostTime());
        m.put("sysUserId", l.getSysUserId());
        m.put("clientIp", l.getClientIp());
        m.put("reason", l.getReason());
        m.put("location", l.getLocation());
        m.put("userAgent", l.getUserAgent());
        m.put("browserName", l.getBrowserName());
        m.put("osName", l.getOsName());
        m.put("createdAt", l.getCreatedAt());
        return m;
    }

    private Map<String, Object> toLoginLogMap(SysLoginLog l) {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("id", l.getId());
        m.put("username", l.getUsername());
        m.put("loginIp", l.getLoginIp());
        m.put("loginMac", l.getLoginMac());
        m.put("statusCode", l.getStatusCode());
        m.put("success", l.getSuccess());
        m.put("reason", l.getReason());
        m.put("location", l.getLocation());
        m.put("sysUserId", l.getSysUserId());
        m.put("userAgent", l.getUserAgent());
        m.put("browserName", l.getBrowserName());
        m.put("osName", l.getOsName());
        m.put("createdAt", l.getCreatedAt());
        m.put("loginTime", l.getLoginTime());
        return m;
    }
}
