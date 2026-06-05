package cn.harnesstemplate.admin.domain.entity;

import com.easy.query.core.annotation.Table;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.experimental.FieldNameConstants;

@Data
@EqualsAndHashCode(callSuper = true)
@FieldNameConstants
@Table("sys_role")
public class SysRole extends BaseAuditEntity {

    private String remark;

    private Boolean isEnabled;

    private String name;

    private String code;

    private Long parentId;
}
