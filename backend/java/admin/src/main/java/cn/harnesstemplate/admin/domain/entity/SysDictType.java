package cn.harnesstemplate.admin.domain.entity;

import com.easy.query.core.annotation.Table;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.experimental.FieldNameConstants;

@Data
@EqualsAndHashCode(callSuper = true)
@FieldNameConstants
@Table("sys_dict_type")
public class SysDictType extends BaseAuditEntity {

    private String remark;

    private Boolean isEnabled;

    private Integer sortOrder;

    private String typeCode;

    private String typeName;
}
