package cn.harnesstemplate.admin.web.dto;

import cn.harnesstemplate.admin.domain.entity.SysRole;
import io.github.linpeilie.annotations.AutoMapper;

import java.util.List;

@AutoMapper(target = SysRole.class)
public record SysRoleCreateRequest(
        String name,
        String code,
        Long parentId,
        String remark,
        List<Long> menuIds,
        List<Long> apiIds
) {
}
