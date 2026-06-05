package cn.harnesstemplate.admin.application.service;

import cn.harnesstemplate.admin.application.port.UserLookupRepository;
import cn.harnesstemplate.admin.domain.model.SysRole;
import cn.harnesstemplate.admin.domain.model.SysUser;
import cn.harnesstemplate.admin.infrastructure.auth.AuthSessionService;
import cn.harnesstemplate.admin.infrastructure.auth.SessionInfo;
import cn.harnesstemplate.admin.web.filter.ServerKeyPairProvider;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

public class AccountService {

    private static final Logger log = LoggerFactory.getLogger(AccountService.class);

    private final UserLookupRepository userLookupRepository;
    private final AuthSessionService authSessionService;
    private final ServerKeyPairProvider serverKeyPairProvider;
    private final BCryptPasswordEncoder passwordEncoder;

    public AccountService(UserLookupRepository userLookupRepository,
                          AuthSessionService authSessionService,
                          ServerKeyPairProvider serverKeyPairProvider) {
        this.userLookupRepository = userLookupRepository;
        this.authSessionService = authSessionService;
        this.serverKeyPairProvider = serverKeyPairProvider;
        // Go uses bcrypt.MinCost = 4
        this.passwordEncoder = new BCryptPasswordEncoder(4);
    }

    /**
     * Password login, matching Go's PwdLogin flow.
     */
    public LoginResult loginByPassword(String username, String rawPassword) {
        Optional<SysUser> userOpt = userLookupRepository.findByUsername(username);
        if (userOpt.isEmpty()) {
            log.debug("User not found: {}", username);
            throw new AccountException("用户名或密码无效");
        }

        SysUser user = userOpt.get();
        if (!passwordEncoder.matches(rawPassword, user.getPassword())) {
            log.debug("Password mismatch for user: {}", username);
            throw new AccountException("用户名或密码无效");
        }

        List<SysRole> roles = userLookupRepository.findRolesByUserId(user.getId());
        List<String> roleCodes = new ArrayList<>();
        List<Long> roleIds = new ArrayList<>();
        for (SysRole role : roles) {
            if (role.getCode() != null && !role.getCode().isEmpty()) {
                roleCodes.add(role.getCode());
            }
            if (role.getId() != null && role.getId() > 0) {
                roleIds.add(role.getId());
            }
        }

        String token = authSessionService.login(user.getId());

        // Use the global server key pair — the frontend encrypts with the
        // global public key, and per-session keys are lost on server restart.
        String publicKeyB64 = serverKeyPairProvider.getPublicKey();

        authSessionService.saveSession(user.getId(), new SessionInfo(
                "", user.getId(), user.getUsername(), roleCodes, roleIds));

        log.info("User {} logged in successfully", username);
        return new LoginResult(token, publicKeyB64);
    }

    /**
     * Logout, matching Go's Logout.
     */
    public void logout(Long loginId) {
        authSessionService.logout(loginId);
    }

    public static class LoginResult {
        private final String token;
        private final String publicKey;

        public LoginResult(String token, String publicKey) {
            this.token = token;
            this.publicKey = publicKey;
        }

        public String getToken() { return token; }
        public String getPublicKey() { return publicKey; }
    }

    public static class AccountException extends RuntimeException {
        public AccountException(String message) {
            super(message);
        }
    }
}
