package cn.harnesstemplate.admin.web.controller;

import cn.harnesstemplate.admin.application.service.AccountService;
import cn.harnesstemplate.admin.infrastructure.auth.AuthSessionService;
import cn.harnesstemplate.admin.web.dto.R;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.web.bind.annotation.*;

import java.util.Map;

@RestController
@RequestMapping("/api/account")
public class AccountController {

    private static final Logger log = LoggerFactory.getLogger(AccountController.class);

    private final AccountService accountService;
    private final AuthSessionService authSessionService;

    public AccountController(AccountService accountService, AuthSessionService authSessionService) {
        this.accountService = accountService;
        this.authSessionService = authSessionService;
    }

    @PostMapping("/login/pwd")
    public R<Map<String, Object>> loginPwd(@RequestBody LoginPwdRequest request) {
        try {
            AccountService.LoginResult result = accountService.loginByPassword(
                    request.username(), request.pwd());
            return R.ok(Map.of(
                    "token", result.getToken(),
                    "publicKey", result.getPublicKey()
            ));
        } catch (AccountService.AccountException e) {
            return R.fail(e.getMessage());
        }
    }

    @GetMapping("/logout")
    public R<Void> logout() {
        try {
            Long loginId = authSessionService.getCurrentLoginId();
            accountService.logout(loginId);
            return R.ok();
        } catch (Exception e) {
            log.debug("Logout error: {}", e.getMessage());
            return R.fail("退出登录失败");
        }
    }

    public record LoginPwdRequest(String username, String pwd) {}
}

