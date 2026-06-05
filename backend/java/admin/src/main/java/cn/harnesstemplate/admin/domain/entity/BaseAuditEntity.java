package cn.harnesstemplate.admin.domain.entity;

import com.easy.query.core.annotation.Column;
import com.easy.query.core.annotation.LogicDelete;
import com.easy.query.core.annotation.UpdateIgnore;
import com.easy.query.core.basic.extension.logicdel.LogicDeleteStrategyEnum;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.experimental.FieldNameConstants;

@Data
@EqualsAndHashCode(callSuper = true)
@FieldNameConstants
public abstract class BaseAuditEntity extends BaseEntity {

    @Column("created_by")
    @UpdateIgnore
    private Long createdBy;

    @Column("updated_by")
    private Long updatedBy;

    @Column("deleted_by")
    @UpdateIgnore
    private Long deletedBy;

    @Column("deleted_at")
    @LogicDelete(strategy = LogicDeleteStrategyEnum.DELETE_LONG_TIMESTAMP)
    @UpdateIgnore
    private Long deletedAt;
}
