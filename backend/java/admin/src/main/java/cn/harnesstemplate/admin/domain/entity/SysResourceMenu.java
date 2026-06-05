package cn.harnesstemplate.admin.domain.entity;

import com.easy.query.core.annotation.Table;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.experimental.FieldNameConstants;

@Data
@EqualsAndHashCode(callSuper = true)
@FieldNameConstants
@Table("sys_resource_menu")
public class SysResourceMenu extends BaseAuditEntity {

    private String remark;

    private Integer sortOrder;

    private String metadata;

    private Boolean isEnabled;

    private String menuType;

    private String path;

    private String redirect;

    private String alias;

    private String name;

    private String component;

    private Long parentId;

    private String treePath;
}
