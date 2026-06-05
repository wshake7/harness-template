package cn.harnesstemplate.admin.web.controller;

import cn.harnesstemplate.admin.web.dto.R;
import cn.harnesstemplate.admin.web.filter.ServerKeyPairProvider;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.Map;

@RestController
@RequestMapping("/api/encrypt")
public class EncryptController {

    private final ServerKeyPairProvider serverKeyPairProvider;

    public EncryptController(ServerKeyPairProvider serverKeyPairProvider) {
        this.serverKeyPairProvider = serverKeyPairProvider;
    }

    @GetMapping("/public/key")
    public R<Map<String, Object>> publicKey() {
        return R.ok(Map.of("publicKey", serverKeyPairProvider.getPublicKey()));
    }
}
