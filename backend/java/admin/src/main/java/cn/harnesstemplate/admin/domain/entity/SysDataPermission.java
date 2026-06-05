package cn.harnesstemplate.admin.domain.entity;

import com.easy.query.core.annotation.Table;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.experimental.FieldNameConstants;

@Data
@EqualsAndHashCode(callSuper = true)
@FieldNameConstants
@Table("sys_data_permission")
public class SysDataPermission extends BaseAuditEntity {

    private String remark;

    private Boolean isEnabled;

    private String subjectType;

    private Long subjectId;

    private String resourceTable;

    private String action;

    private String actionKey;

    private String scopeType;

    private String scopeField;

    private String scopeValues;

    private String conditions;

    private Integer priority;
}
