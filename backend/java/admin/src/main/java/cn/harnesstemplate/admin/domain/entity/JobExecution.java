package cn.harnesstemplate.admin.domain.entity;

import com.easy.query.core.annotation.Table;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.experimental.FieldNameConstants;

import java.time.LocalDateTime;

@Data
@EqualsAndHashCode(callSuper = true)
@FieldNameConstants
@Table("job_execution")
public class JobExecution extends BaseEntity {

    private String jobCode;

    private String temporalWorkflowId;

    private String temporalRunId;

    private LocalDateTime triggerTime;

    private LocalDateTime startTime;

    private LocalDateTime endTime;

    private String status;

    private String inputJson;

    private String resultJson;

    private String errorMessage;

    private Integer retryCount;
}
