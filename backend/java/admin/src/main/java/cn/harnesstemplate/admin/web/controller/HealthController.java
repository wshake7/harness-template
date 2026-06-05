package cn.harnesstemplate.admin.web.controller;

import cn.harnesstemplate.admin.web.dto.R;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.time.Instant;
import java.util.Map;

@RestController
@RequestMapping("/api/health")
public class HealthController {

    @GetMapping("/smoke")
    public R<Map<String, Object>> smoke() {
        return R.ok(Map.of(
                "timestamp", Instant.now().toString(),
                "service", "admin-java"
        ));
    }
}
