package cn.harnesstemplate.admin.domain.entity;

import com.easy.query.core.annotation.Column;
import com.easy.query.core.annotation.Table;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.experimental.FieldNameConstants;

@Data
@EqualsAndHashCode(callSuper = true)
@FieldNameConstants
@Table("sys_user")
public class SysUser extends BaseAuditEntity {

    private String remark;

    private Boolean isEnabled;

    private String username;

    private String nickname;

    private String password;

    private String languageCode;
}
