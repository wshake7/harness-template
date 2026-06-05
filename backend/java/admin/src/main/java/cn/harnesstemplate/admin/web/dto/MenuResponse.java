package cn.harnesstemplate.admin.web.dto;

import cn.harnesstemplate.admin.domain.entity.SysResourceMenu;
import io.github.linpeilie.annotations.AutoMapper;

import java.time.LocalDateTime;
import java.util.List;

@AutoMapper(target = SysResourceMenu.class)
public record MenuResponse(
        Long id,
        String menuType,
        String path,
        String redirect,
        String alias,
        String name,
        String component,
        Long parentId,
        String treePath,
        Integer sortOrder,
        Boolean isEnabled,
        String remark,
        List<MenuResponse> children,
        List<Long> apiIds,
        LocalDateTime createdAt,
        LocalDateTime updatedAt
) {
}
