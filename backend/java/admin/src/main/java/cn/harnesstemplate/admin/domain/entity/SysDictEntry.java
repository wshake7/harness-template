package cn.harnesstemplate.admin.domain.entity;

import com.easy.query.core.annotation.Table;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.experimental.FieldNameConstants;

@Data
@EqualsAndHashCode(callSuper = true)
@FieldNameConstants
@Table("sys_dict_entry")
public class SysDictEntry extends BaseAuditEntity {

    private String remark;

    private Boolean isEnabled;

    private Integer sortOrder;

    private String labelComponent;

    private String entryLabel;

    private String entryValue;

    private String languageCode;

    private Long sysDictTypeId;
}
