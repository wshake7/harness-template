package cn.harnesstemplate.admin.application.service.job;

import cn.harnesstemplate.admin.domain.entity.JobExecution;
import cn.harnesstemplate.admin.domain.entity.JobSchedule;
import cn.harnesstemplate.admin.infrastructure.temporal.TemporalClientFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.LinkedHashMap;
import java.util.Map;

@Service
public class TemporalService {

    private static final Logger log = LoggerFactory.getLogger(TemporalService.class);

    private final TemporalClientFactory temporalFactory;
    private final JobExecutionQueryService executionService;

    public TemporalService(TemporalClientFactory temporalFactory, JobExecutionQueryService executionService) {
        this.temporalFactory = temporalFactory;
        this.executionService = executionService;
    }

    public boolean isAvailable() {
        return temporalFactory != null && temporalFactory.isEnabled() && temporalFactory.getScheduleClient() != null;
    }

    public void syncSchedule(JobSchedule schedule) {
        if (!isAvailable()) {
            log.debug("Temporal not available, skipping sync for {}", schedule.getJobCode());
            return;
        }
        try {
            // TODO: Implement schedule creation/update via Temporal client when Temporal is available
            log.info("Schedule sync for {} would be dispatched to Temporal", schedule.getJobCode());
        } catch (Exception e) {
            log.warn("Failed to sync schedule {}: {}", schedule.getJobCode(), e.getMessage());
        }
    }

    public void deleteSchedule(JobSchedule schedule) {
        if (!isAvailable()) return;
        try {
            var handle = temporalFactory.getScheduleClient().getHandle(
                    schedule.getTemporalScheduleId() != null ? schedule.getTemporalScheduleId() : schedule.getJobCode());
            handle.delete();
        } catch (Exception e) {
            log.warn("Failed to delete schedule: {}", e.getMessage());
        }
    }

    public void pauseSchedule(JobSchedule schedule) {
        if (!isAvailable()) return;
        String scheduleId = schedule.getTemporalScheduleId() != null ? schedule.getTemporalScheduleId() : schedule.getJobCode();
        try {
            var handle = temporalFactory.getScheduleClient().getHandle(scheduleId);
            handle.pause("Paused via admin");
        } catch (Exception e) {
            log.warn("Failed to pause schedule {}: {}", scheduleId, e.getMessage());
        }
    }

    public void unpauseSchedule(JobSchedule schedule) {
        if (!isAvailable()) return;
        String scheduleId = schedule.getTemporalScheduleId() != null ? schedule.getTemporalScheduleId() : schedule.getJobCode();
        try {
            var handle = temporalFactory.getScheduleClient().getHandle(scheduleId);
            handle.unpause("Unpaused via admin");
        } catch (Exception e) {
            log.warn("Failed to unpause schedule {}: {}", scheduleId, e.getMessage());
        }
    }

    public void triggerSchedule(JobSchedule schedule) {
        if (!isAvailable()) {
            createFailedExecution(schedule, "Temporal unavailable — execution was not dispatched");
            return;
        }
        String scheduleId = schedule.getTemporalScheduleId() != null ? schedule.getTemporalScheduleId() : schedule.getJobCode();
        try {
            var handle = temporalFactory.getScheduleClient().getHandle(scheduleId);
            handle.trigger();
        } catch (Exception e) {
            log.warn("Failed to trigger schedule {}: {}", scheduleId, e.getMessage());
            createFailedExecution(schedule, e.getMessage());
        }
    }

    public void cancelExecution(JobExecution execution) {
        if (!isAvailable()) {
            executionService.updateCanceled(execution.getId());
            return;
        }
        try {
            if (execution.getTemporalWorkflowId() == null || execution.getTemporalWorkflowId().isEmpty()) {
                executionService.updateCanceled(execution.getId());
                return;
            }
            var wfStub = temporalFactory.getWorkflowClient().newUntypedWorkflowStub(
                    execution.getTemporalWorkflowId(),
                    java.util.Optional.ofNullable(execution.getTemporalRunId()),
                    java.util.Optional.empty());
            wfStub.cancel();
            executionService.updateCanceled(execution.getId());
        } catch (Exception e) {
            log.warn("Failed to cancel execution {}: {}", execution.getTemporalWorkflowId(), e.getMessage());
            executionService.updateCanceled(execution.getId());
        }
    }

    public Map<String, Object> getWorkflowTypes() {
        Map<String, Object> types = new LinkedHashMap<>();
        types.put("PrintCountWorkflow", "PrintCountWorkflow");
        types.put("KnowledgeDocumentIndexWorkflow", "KnowledgeDocumentIndexWorkflow");
        return types;
    }

    private void createFailedExecution(JobSchedule schedule, String reason) {
        JobExecution exec = new JobExecution();
        exec.setJobCode(schedule.getJobCode());
        exec.setTemporalWorkflowId(schedule.getJobCode() + "-sim-" + System.currentTimeMillis());
        exec.setTemporalRunId("");
        exec.setTriggerTime(LocalDateTime.now());
        exec.setStartTime(LocalDateTime.now());
        exec.setStatus("FAILED");
        exec.setInputJson(schedule.getInputJson());
        exec.setRetryCount(0);
        exec.setErrorMessage(reason);
        executionService.save(exec);
    }
}
