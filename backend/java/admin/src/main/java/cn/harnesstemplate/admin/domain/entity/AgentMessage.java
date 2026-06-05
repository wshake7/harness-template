package cn.harnesstemplate.admin.domain.entity;

import com.easy.query.core.annotation.Table;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.experimental.FieldNameConstants;

@Data
@EqualsAndHashCode(callSuper = true)
@FieldNameConstants
@Table("agent_message")
public class AgentMessage extends BaseAuditEntity {

    private Long sessionId;

    private String role;

    private String content;

    private Integer tokenCount;

    private String metadata;
}
