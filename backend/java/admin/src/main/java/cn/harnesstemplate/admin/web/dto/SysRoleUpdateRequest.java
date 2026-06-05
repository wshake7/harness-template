package cn.harnesstemplate.admin.web.dto;

public record SysRoleUpdateRequest(
        Long id,
        String name,
        String code,
        Long parentId,
        Boolean isEnabled,
        String remark
) {
}
