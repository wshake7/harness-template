package cn.harnesstemplate.admin.web.filter;

import cn.harnesstemplate.admin.infrastructure.auth.CryptoService;
import cn.harnesstemplate.admin.web.dto.R;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.ServletInputStream;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletRequestWrapper;
import jakarta.servlet.http.HttpServletResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.core.Ordered;
import org.springframework.core.annotation.Order;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;
import org.springframework.web.util.ContentCachingResponseWrapper;

import java.io.*;
import java.nio.charset.StandardCharsets;
import java.util.LinkedHashMap;
import java.util.Map;

/**
 * Decrypts encrypted request bodies and encrypts responses.
 * Matching Go's EncryptMiddleware.
 */
@Component
@Order(Ordered.HIGHEST_PRECEDENCE + 20)
public class EncryptFilter extends OncePerRequestFilter {

    private static final Logger log = LoggerFactory.getLogger(EncryptFilter.class);
    private static final ObjectMapper objectMapper = new ObjectMapper();

    private final CryptoService cryptoService;
    private final ServerKeyPairProvider serverKeyPairProvider;

    public EncryptFilter(CryptoService cryptoService,
                         ServerKeyPairProvider serverKeyPairProvider) {
        this.cryptoService = cryptoService;
        this.serverKeyPairProvider = serverKeyPairProvider;
    }

    @Override
    protected void doFilterInternal(HttpServletRequest request,
                                    HttpServletResponse response,
                                    FilterChain filterChain) throws ServletException, IOException {
        // SSE endpoints handle their own per-event encryption; buffering the stream would break them.
        if (request.getRequestURI().endsWith("/events")) {
            filterChain.doFilter(request, response);
            return;
        }

        String encryptedKey = request.getHeader("X-Request-Encrypted-Key");

        if (encryptedKey == null || encryptedKey.isEmpty()) {
            // Not an encrypted request, pass through
            filterChain.doFilter(request, response);
            return;
        }

        // Skip encryption for multipart uploads
        String contentType = request.getContentType();
        if (contentType != null && contentType.startsWith("multipart/form-data")) {
            filterChain.doFilter(request, response);
            return;
        }

        // Get private key: from session if authenticated, otherwise from global server key
        // Read token directly from request header because this filter may run before
        // Sa-Token's SaTokenContextFilter sets up the context.
        String privateKeyPem = resolvePrivateKey(request);

        // Decrypt AES key
        String aesKeyBase64;
        try {
            aesKeyBase64 = cryptoService.rsaDecrypt(encryptedKey,
                    CryptoService.parsePrivateKeyPem(privateKeyPem));
        } catch (Exception e) {
            log.debug("RSA decrypt failed: {}", e.getMessage());
            writeError(response, 5, "请求错误");
            return;
        }

        // Build AAD
        String aad = buildAadFromRequest(request);

        // Cache the original body once so downstream code never sees an already-consumed stream.
        HttpServletRequest requestToUse = request;
        byte[] rawBody = readBodyBytes(request);
        if (rawBody.length > 0) {
            String sign = firstHeader(request, "X-Request-Signature", "X-Sign");
            try {
                if (sign != null && !sign.isEmpty()) {
                    String ciphertextBody = new String(rawBody, StandardCharsets.UTF_8);
                    byte[] decrypted = cryptoService.aesDecryptCiphertextAndTag(
                            ciphertextBody, sign, aesKeyBase64, aad);
                    requestToUse = new CachedBodyRequestWrapper(request, decrypted);
                } else {
                    requestToUse = new CachedBodyRequestWrapper(request, rawBody);
                }
            } catch (Exception e) {
                log.debug("Request body decryption failed: {}", e.getMessage());
                writeError(response, 2, "请求错误");
                return;
            }
        }

        // Wrap response for encryption
        ContentCachingResponseWrapper wrappedResponse = new ContentCachingResponseWrapper(response);

        try {
            filterChain.doFilter(requestToUse, wrappedResponse);

            byte[] responseBody = wrappedResponse.getContentAsByteArray();
            if (responseBody.length > 0) {
                String plainResponse = new String(responseBody, StandardCharsets.UTF_8);
                CryptoService.EncryptResult result = cryptoService.aesEncrypt(
                        plainResponse, aesKeyBase64, "");
                byte[] encryptedResponse = result.combined.getBytes(StandardCharsets.UTF_8);

                wrappedResponse.resetBuffer();
                wrappedResponse.setHeader("X-Response-Is-Encrypt", "true");
                wrappedResponse.setContentType("application/json");
                wrappedResponse.setCharacterEncoding("UTF-8");
                wrappedResponse.setContentLength(encryptedResponse.length);
                wrappedResponse.getOutputStream().write(encryptedResponse);
            }
        } catch (Exception e) {
            log.error("Response encryption failed: {}", e.getMessage());
        } finally {
            wrappedResponse.copyBodyToResponse();
        }
    }

    /**
     * Always use the global server key pair for decryption.
     * Per-session keys are not used because they are lost on server restart,
     * causing a mismatch with the frontend's cached public key.
     */
    private String resolvePrivateKey(HttpServletRequest request) {
        return serverKeyPairProvider.getPrivateKeyPem();
    }

    private void writeError(HttpServletResponse response, int code, String msg) throws IOException {
        response.setStatus(HttpServletResponse.SC_OK);
        response.setContentType("application/json");
        response.setCharacterEncoding("UTF-8");
        R<Void> error = new R<>(code, msg, null);
        response.getWriter().write(objectMapper.writeValueAsString(error));
    }

    private byte[] readBodyBytes(HttpServletRequest request) throws IOException {
        try (InputStream inputStream = request.getInputStream();
             ByteArrayOutputStream outputStream = new ByteArrayOutputStream()) {
            inputStream.transferTo(outputStream);
            return outputStream.toByteArray();
        }
    }

    private String buildAadFromRequest(HttpServletRequest request) {
        Map<String, String> params = new LinkedHashMap<>();

        String requestId = firstHeader(request, "X-Request-ID");
        if (requestId != null && !requestId.isEmpty()) {
            params.put("X-Request-ID", requestId);
        }

        String timestamp = firstHeader(request, "X-Request-Timestamp", "X-Timestamp");
        if (timestamp != null && !timestamp.isEmpty()) {
            params.put(timestampHeaderName(request), timestamp);
        }

        // Add query params
        for (Map.Entry<String, String[]> e : request.getParameterMap().entrySet()) {
            String[] values = e.getValue();
            if (values != null && values.length > 0) {
                params.put(e.getKey(), values[0]);
            }
        }

        return cryptoService.buildAad(params);
    }

    private String firstHeader(HttpServletRequest request, String... names) {
        for (String name : names) {
            String value = request.getHeader(name);
            if (value != null && !value.isEmpty()) {
                return value;
            }
        }
        return null;
    }

    private String timestampHeaderName(HttpServletRequest request) {
        String frontendTimestamp = request.getHeader("X-Request-Timestamp");
        if (frontendTimestamp != null && !frontendTimestamp.isEmpty()) {
            return "X-Request-Timestamp";
        }
        return "X-Timestamp";
    }

    // ---- Request wrapper with cached body ----

    static class CachedBodyRequestWrapper extends HttpServletRequestWrapper {
        private final byte[] body;

        CachedBodyRequestWrapper(HttpServletRequest request, byte[] body) {
            super(request);
            this.body = body;
        }

        @Override
        public ServletInputStream getInputStream() {
            ByteArrayInputStream bis = new ByteArrayInputStream(body);
            return new ServletInputStream() {
                @Override
                public int read() { return bis.read(); }
                @Override
                public boolean isFinished() { return bis.available() == 0; }
                @Override
                public boolean isReady() { return true; }
                @Override
                public void setReadListener(jakarta.servlet.ReadListener listener) {}
            };
        }

        @Override
        public BufferedReader getReader() {
            return new BufferedReader(new InputStreamReader(getInputStream(), StandardCharsets.UTF_8));
        }

        @Override
        public int getContentLength() {
            return body.length;
        }

        @Override
        public long getContentLengthLong() {
            return body.length;
        }
    }
}
