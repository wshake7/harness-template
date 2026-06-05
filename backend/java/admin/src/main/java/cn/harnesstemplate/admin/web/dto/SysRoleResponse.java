package cn.harnesstemplate.admin.web.dto;

import cn.harnesstemplate.admin.domain.entity.SysRole;
import io.github.linpeilie.annotations.AutoMapper;

import java.time.LocalDateTime;
import java.util.List;

@AutoMapper(target = SysRole.class)
public record SysRoleResponse(
        Long id,
        String name,
        String code,
        Long parentId,
        String remark,
        Boolean isEnabled,
        List<Long> menuIds,
        List<Long> apiIds,
        List<SysRoleResponse> children,
        LocalDateTime createdAt,
        LocalDateTime updatedAt
) {
}
