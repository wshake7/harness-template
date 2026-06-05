package cn.harnesstemplate.admin.web.controller;

import cn.harnesstemplate.admin.infrastructure.auth.CryptoService;
import cn.harnesstemplate.admin.web.filter.ServerKeyPairProvider;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.annotation.PreDestroy;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.servlet.mvc.method.annotation.SseEmitter;

import java.io.IOException;
import java.util.Map;
import java.util.concurrent.CopyOnWriteArrayList;

@RestController
public class EventController {

    private static final Logger log = LoggerFactory.getLogger(EventController.class);
    private static final ObjectMapper objectMapper = new ObjectMapper();

    private final CryptoService cryptoService;
    private final ServerKeyPairProvider serverKeyPairProvider;
    private final CopyOnWriteArrayList<SseEmitter> activeEmitters = new CopyOnWriteArrayList<>();

    public EventController(CryptoService cryptoService,
                           ServerKeyPairProvider serverKeyPairProvider) {
        this.cryptoService = cryptoService;
        this.serverKeyPairProvider = serverKeyPairProvider;
    }

    @GetMapping("/api/events")
    public SseEmitter events(HttpServletRequest request, HttpServletResponse response) {
        String encryptedKey = request.getHeader("X-Request-Encrypted-Key");
        if (encryptedKey == null || encryptedKey.isEmpty()) {
            throw new IllegalArgumentException("Missing X-Request-Encrypted-Key header");
        }

        String aesKeyBase64;
        try {
            aesKeyBase64 = cryptoService.rsaDecrypt(encryptedKey,
                    CryptoService.parsePrivateKeyPem(serverKeyPairProvider.getPrivateKeyPem()));
        } catch (Exception e) {
            log.debug("RSA decrypt of SSE request key failed: {}", e.getMessage());
            throw new IllegalArgumentException("Failed to decrypt request key", e);
        }

        response.setHeader("Cache-Control", "no-cache");
        response.setHeader("Connection", "keep-alive");
        response.setHeader("X-Response-Is-Encrypt", "true");

        SseEmitter emitter = new SseEmitter(Long.MAX_VALUE);
        activeEmitters.add(emitter);

        Thread.ofVirtual().start(() -> {
            int count = 0;
            try {
                while (true) {
                    Thread.sleep(60_000);
                    count++;

                    String plainData = objectMapper.writeValueAsString(Map.of("count", count));
                    CryptoService.EncryptResult result = cryptoService.aesEncrypt(plainData, aesKeyBase64, "");
                    String eventData = objectMapper.writeValueAsString(Map.of("payload", result.combined));

                    emitter.send(SseEmitter.event()
                            .name("count")
                            .data(eventData));
                }
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            } catch (IOException e) {
                log.debug("SSE client disconnected");
            } catch (Exception e) {
                log.debug("SSE send error: {}", e.getMessage());
            } finally {
                activeEmitters.remove(emitter);
            }
        });

        emitter.onCompletion(() -> activeEmitters.remove(emitter));
        emitter.onTimeout(() -> activeEmitters.remove(emitter));

        return emitter;
    }

    @PreDestroy
    public void shutdown() {
        for (SseEmitter emitter : activeEmitters) {
            try {
                emitter.complete();
            } catch (Exception ignored) {
            }
        }
        activeEmitters.clear();
    }
}
