package cn.harnesstemplate.admin.web.controller;

import cn.harnesstemplate.admin.application.service.system.SysUserQueryService;
import cn.harnesstemplate.admin.application.service.system.SysUserRoleQueryService;
import cn.harnesstemplate.admin.domain.entity.SysUser;
import cn.harnesstemplate.admin.infrastructure.auth.SessionInfo;
import cn.harnesstemplate.admin.web.dto.*;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.web.bind.annotation.*;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/sys/user")
public class UserController {

    private final SysUserQueryService userService;
    private final SysUserRoleQueryService userRoleService;
    private final BCryptPasswordEncoder passwordEncoder = new BCryptPasswordEncoder(4);

    public UserController(SysUserQueryService userService, SysUserRoleQueryService userRoleService) {
        this.userService = userService;
        this.userRoleService = userRoleService;
    }

    @PostMapping("/list")
    public R<PageResult<SysUserResponse>> list(@RequestBody Map<String, Object> req) {
        int page = req.containsKey("page") ? ((Number) req.get("page")).intValue() : 1;
        int pageSize = req.containsKey("pageSize") ? ((Number) req.get("pageSize")).intValue() : 10;
        String username = req.containsKey("username") ? (String) req.get("username") : null;

        List<SysUser> users = userService.page(page, pageSize, username);
        long total = userService.count(username);

        List<SysUserResponse> items = new ArrayList<>();
        for (SysUser u : users) {
            SysUserResponse resp = toResponse(u);
            items.add(resp);
        }
        return R.ok(PageResult.of(items, total));
    }

    @PostMapping("/create")
    public R<Void> create(@RequestBody SysUserCreateRequest req) {
        SysUser existing = userService.findByUsername(req.username());
        if (existing != null) {
            return R.fail("用户名已存在");
        }
        SysUser user = new SysUser();
        user.setUsername(req.username());
        user.setNickname(req.nickname());
        user.setPassword(passwordEncoder.encode(req.password()));
        user.setLanguageCode(req.languageCode() != null ? req.languageCode() : "zh-CN");
        user.setIsEnabled(true);
        userService.save(user);
        if (req.roleIds() != null) {
            for (Long roleId : req.roleIds()) {
                userRoleService.save(user.getId(), roleId);
            }
        }
        return R.ok();
    }

    @PostMapping("/update")
    public R<Void> update(@RequestBody SysUserUpdateRequest req) {
        SysUser user = userService.findById(req.id());
        if (user == null) {
            return R.fail("用户不存在");
        }
        if (req.nickname() != null) user.setNickname(req.nickname());
        if (req.languageCode() != null) user.setLanguageCode(req.languageCode());
        if (req.isEnabled() != null) user.setIsEnabled(req.isEnabled());
        if (req.remark() != null) user.setRemark(req.remark());
        userService.update(user);
        return R.ok();
    }

    @PostMapping("/del")
    public R<Void> del(@RequestBody Map<String, List<Long>> req) {
        List<Long> ids = req.get("ids");
        if (ids == null || ids.isEmpty()) {
            return R.fail("ids不能为空");
        }
        for (Long id : ids) {
            userService.delete(id);
        }
        return R.ok();
    }

    private SysUserResponse toResponse(SysUser u) {
        return new SysUserResponse(
                u.getId(),
                u.getUsername(),
                u.getNickname(),
                u.getLanguageCode(),
                u.getIsEnabled(),
                null,
                null,
                u.getCreatedAt(),
                u.getUpdatedAt()
        );
    }
}
