package cn.harnesstemplate.admin.web.dto;

import cn.harnesstemplate.admin.domain.entity.SysUser;
import io.github.linpeilie.annotations.AutoMapper;

import java.util.List;

@AutoMapper(target = SysUser.class)
public record SysUserCreateRequest(
        String username,
        String nickname,
        String password,
        String languageCode,
        List<Long> roleIds
) {
}
