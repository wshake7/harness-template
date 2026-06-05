package cn.harnesstemplate.admin.domain.entity;

import com.easy.query.core.annotation.Table;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.experimental.FieldNameConstants;

@Data
@EqualsAndHashCode(callSuper = true)
@FieldNameConstants
@Table("sys_resource_api")
public class SysResourceApi extends BaseAuditEntity {

    private String remark;

    private Integer sortOrder;

    private Boolean isEnabled;

    private String module;

    private String path;

    private String method;
}
