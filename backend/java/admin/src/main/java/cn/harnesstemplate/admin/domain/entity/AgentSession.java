package cn.harnesstemplate.admin.domain.entity;

import com.easy.query.core.annotation.Table;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.experimental.FieldNameConstants;

@Data
@EqualsAndHashCode(callSuper = true)
@FieldNameConstants
@Table("agent_session")
public class AgentSession extends BaseAuditEntity {

    private String remark;

    private String sessionId;

    private String title;

    private String status;

    private String metadata;
}
