package cn.harnesstemplate.admin.web.filter;

import cn.harnesstemplate.admin.web.dto.R;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.core.Ordered;
import org.springframework.core.annotation.Order;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;

/**
 * Validates X-Timestamp header is within the allowed time window (5 minutes).
 * Matching Go's TimestampMiddleware.
 */
@Component
@Order(Ordered.HIGHEST_PRECEDENCE + 10)
public class TimestampFilter extends OncePerRequestFilter {

    private static final Logger log = LoggerFactory.getLogger(TimestampFilter.class);
    private static final long REQUEST_EXPIRE_MS = 5 * 60 * 1000; // 5 minutes
    private static final ObjectMapper objectMapper = new ObjectMapper();

    @Override
    protected void doFilterInternal(HttpServletRequest request,
                                    HttpServletResponse response,
                                    FilterChain filterChain) throws ServletException, IOException {
        String timestampHeader = request.getHeader("X-Request-Timestamp");
        if (timestampHeader == null || timestampHeader.isEmpty()) {
            timestampHeader = request.getHeader("X-Timestamp");
        }
        if (timestampHeader != null && !timestampHeader.isEmpty()) {
            try {
                long timestamp = Long.parseLong(timestampHeader);
                long now = System.currentTimeMillis();
                if (Math.abs(now - timestamp) > REQUEST_EXPIRE_MS) {
                    log.debug("Request timestamp expired: {} (now: {})", timestamp, now);
                    writeError(response, 3, "请求已过期");
                    return;
                }
            } catch (NumberFormatException e) {
                log.debug("Invalid X-Timestamp header: {}", timestampHeader);
                writeError(response, 2, "请求错误");
                return;
            }
        }
        filterChain.doFilter(request, response);
    }

    private void writeError(HttpServletResponse response, int code, String msg) throws IOException {
        response.setStatus(HttpServletResponse.SC_OK);
        response.setContentType("application/json");
        response.setCharacterEncoding("UTF-8");
        R<Void> error = new R<>(code, msg, null);
        response.getWriter().write(objectMapper.writeValueAsString(error));
    }
}
