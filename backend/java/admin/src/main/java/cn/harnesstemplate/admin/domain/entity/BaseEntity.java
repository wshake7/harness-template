package cn.harnesstemplate.admin.domain.entity;

import com.easy.query.core.annotation.Column;
import com.easy.query.core.annotation.UpdateIgnore;
import lombok.Data;
import lombok.experimental.FieldNameConstants;

import java.io.Serializable;
import java.time.LocalDateTime;

@Data
@FieldNameConstants
public abstract class BaseEntity implements Serializable, Cloneable {

    private static final long serialVersionUID = 1L;

    @Column(value = "id", primaryKey = true, generatedKey = true)
    private Long id;

    @Column("created_at")
    @UpdateIgnore
    private LocalDateTime createdAt;

    @Column("updated_at")
    private LocalDateTime updatedAt;
}
