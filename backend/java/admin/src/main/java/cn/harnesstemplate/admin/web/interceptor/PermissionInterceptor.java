package cn.harnesstemplate.admin.web.interceptor;

import cn.harnesstemplate.admin.infrastructure.auth.SessionInfo;
import cn.harnesstemplate.admin.infrastructure.permission.CasbinPermissionService;
import cn.harnesstemplate.admin.web.dto.R;
import cn.dev33.satoken.stp.StpUtil;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.MediaType;

import java.nio.charset.StandardCharsets;
import java.util.List;

/**
 * Permission interceptor matching Go's CasbinAPIMiddleware.
 * Enforces API permissions via jCasbin with OR logic across user/role subjects.
 */
public class PermissionInterceptor implements org.springframework.web.servlet.HandlerInterceptor {

    private static final Logger log = LoggerFactory.getLogger(PermissionInterceptor.class);
    private static final ObjectMapper mapper = new ObjectMapper();

    private final CasbinPermissionService casbinPermissionService;

    public PermissionInterceptor(CasbinPermissionService casbinPermissionService) {
        this.casbinPermissionService = casbinPermissionService;
    }

    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response,
                             Object handler) throws Exception {
        if (!StpUtil.isLogin()) {
            writeUnauthorized(response);
            return false;
        }

        SessionInfo sessionInfo = (SessionInfo) StpUtil.getSessionByLoginId(StpUtil.getLoginIdAsLong())
                .get(SessionInfo.SESSION_KEY);
        if (sessionInfo == null) {
            log.debug("Casbin auth skipped: session info is empty");
            writeUnauthorized(response);
            return false;
        }

        List<String> subjects = CasbinPermissionService.buildSubjects(
                sessionInfo.getId(), sessionInfo.getRoleCodes());

        String path = request.getRequestURI();
        String method = request.getMethod();

        if (!casbinPermissionService.enforce(subjects, path, method)) {
            writeUnauthorized(response);
            return false;
        }

        return true;
    }

    private void writeUnauthorized(HttpServletResponse response) throws Exception {
        response.setStatus(HttpServletResponse.SC_OK);
        response.setContentType(MediaType.APPLICATION_JSON_VALUE);
        response.setCharacterEncoding(StandardCharsets.UTF_8.name());
        response.getWriter().write(mapper.writeValueAsString(R.unauthorized()));
    }
}
