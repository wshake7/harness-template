package cn.harnesstemplate.admin.web.interceptor;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.web.servlet.HandlerInterceptor;

/**
 * Sets the language from X-Language header, falling back to default.
 * Matching Go's LanguageMiddleware.
 */
public class LanguageInterceptor implements HandlerInterceptor {

    private static final Logger log = LoggerFactory.getLogger(LanguageInterceptor.class);
    public static final String LANGUAGE_ATTRIBUTE = "language";

    private final String defaultLanguage;

    public LanguageInterceptor(String defaultLanguage) {
        this.defaultLanguage = defaultLanguage;
    }

    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response,
                             Object handler) {
        String language = request.getHeader("X-Language");
        if (language == null || language.trim().isEmpty()) {
            language = defaultLanguage;
        } else {
            language = language.trim();
        }
        request.setAttribute(LANGUAGE_ATTRIBUTE, language);
        log.debug("Request language: {}", language);
        return true;
    }
}
