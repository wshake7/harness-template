package cn.harnesstemplate.admin.config;

import cn.harnesstemplate.admin.infrastructure.permission.CasbinPermissionService;
import cn.harnesstemplate.admin.web.interceptor.LanguageInterceptor;
import cn.harnesstemplate.admin.web.interceptor.PermissionInterceptor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.InterceptorRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;

@Configuration
public class WebConfig implements WebMvcConfigurer {

    @Value("${admin.default-language:zh-CN}")
    private String defaultLanguage;

    private final CasbinPermissionService casbinPermissionService;

    public WebConfig(CasbinPermissionService casbinPermissionService) {
        this.casbinPermissionService = casbinPermissionService;
    }

    @Override
    public void addInterceptors(InterceptorRegistry registry) {
        registry.addInterceptor(new LanguageInterceptor(defaultLanguage))
                .addPathPatterns("/api/**")
                .order(0);

        registry.addInterceptor(new PermissionInterceptor(casbinPermissionService))
                .addPathPatterns("/api/**")
                .excludePathPatterns("/api/account/login/pwd",
                        "/api/account/logout",
                        "/api/encrypt/public/key",
                        "/api/events",
                        "/api/health/**",
                        "/api/public/**")
                .order(1);
    }
}
