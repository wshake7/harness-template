package cn.harnesstemplate.admin.web.dto;

import cn.harnesstemplate.admin.domain.entity.SysUser;
import io.github.linpeilie.annotations.AutoMapper;

import java.time.LocalDateTime;
import java.util.List;

@AutoMapper(target = SysUser.class)
public record SysUserResponse(
        Long id,
        String username,
        String nickname,
        String languageCode,
        Boolean isEnabled,
        List<String> roleCodes,
        List<Long> roleIds,
        LocalDateTime createdAt,
        LocalDateTime updatedAt
) {
}
