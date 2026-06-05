package cn.harnesstemplate.admin.web.dto;

public record SysUserUpdateRequest(
        Long id,
        String nickname,
        String languageCode,
        Boolean isEnabled,
        String remark
) {
}
