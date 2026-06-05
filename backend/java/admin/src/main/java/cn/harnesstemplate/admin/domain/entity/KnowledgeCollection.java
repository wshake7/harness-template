package cn.harnesstemplate.admin.domain.entity;

import com.easy.query.core.annotation.Table;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.experimental.FieldNameConstants;

@Data
@EqualsAndHashCode(callSuper = true)
@FieldNameConstants
@Table("knowledge_collection")
public class KnowledgeCollection extends BaseAuditEntity {

    private String remark;

    private Boolean isEnabled;

    private String collectionName;

    private String displayName;

    private String metricType;

    private String indexType;
}
