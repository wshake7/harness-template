package cn.harnesstemplate.admin.domain.entity;

import com.easy.query.core.annotation.Table;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.experimental.FieldNameConstants;

import java.time.LocalDateTime;

@Data
@EqualsAndHashCode(callSuper = true)
@FieldNameConstants
@Table("job_schedule")
public class JobSchedule extends BaseEntity {

    private String jobCode;

    private String jobName;

    private String workflowType;

    private String taskQueue;

    private String scheduleType;

    private String cronExpr;

    private Integer intervalSeconds;

    private LocalDateTime startTime;

    private LocalDateTime endTime;

    private String inputJson;

    private String status;

    private String temporalScheduleId;

    private String temporalWorkflowIdPrefix;

    private String description;
}
