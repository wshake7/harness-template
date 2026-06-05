package cn.harnesstemplate.admin.domain.entity;

import com.easy.query.core.annotation.Table;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.experimental.FieldNameConstants;

@Data
@EqualsAndHashCode(callSuper = true)
@FieldNameConstants
@Table("knowledge_document")
public class KnowledgeDocument extends BaseAuditEntity {

    private String remark;

    private Boolean isEnabled;

    private Long collectionId;

    private String documentId;

    private String title;

    private String content;

    private String contentType;

    private String source;

    private Integer chunkIndex;

    private Integer totalChunks;

    private String vectorStatus;

    private String vectorId;

    private String metadata;

    private String indexingError;

    private Long lastIndexedAt;
}
