package cn.harnesstemplate.admin.config;

import cn.harnesstemplate.admin.application.port.UserLookupRepository;
import cn.harnesstemplate.admin.application.service.AccountService;
import cn.harnesstemplate.admin.infrastructure.auth.AuthSessionService;
import cn.harnesstemplate.admin.infrastructure.auth.CryptoService;
import cn.harnesstemplate.admin.infrastructure.permission.CasbinPermissionService;
import cn.harnesstemplate.admin.web.filter.ServerKeyPairProvider;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import javax.sql.DataSource;

@Configuration
public class AuthConfig {

    @Bean
    public CryptoService cryptoService() {
        return new CryptoService();
    }

    @Bean
    public AuthSessionService authSessionService() {
        return new AuthSessionService();
    }

    @Bean
    public CasbinPermissionService casbinPermissionService(DataSource dataSource) {
        return new CasbinPermissionService(dataSource);
    }

    @Bean
    public AccountService accountService(UserLookupRepository userLookupRepository,
                                         AuthSessionService authSessionService,
                                         ServerKeyPairProvider serverKeyPairProvider) {
        return new AccountService(userLookupRepository, authSessionService, serverKeyPairProvider);
    }
}
