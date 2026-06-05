package cn.harnesstemplate.admin.domain.entity;

import com.easy.query.core.annotation.Table;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.experimental.FieldNameConstants;

@Data
@EqualsAndHashCode(callSuper = true)
@FieldNameConstants
@Table("sys_language_type")
public class SysLanguageType extends BaseAuditEntity {

    private Boolean isEnabled;

    private Integer sortOrder;

    private String typeCode;

    private String typeName;

    private Boolean isDefault;
}
