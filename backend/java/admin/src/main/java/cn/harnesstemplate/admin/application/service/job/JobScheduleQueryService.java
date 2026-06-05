package cn.harnesstemplate.admin.application.service.job;

import cn.harnesstemplate.admin.domain.entity.JobSchedule;
import cn.harnesstemplate.admin.infrastructure.persistence.EasyQuerySupport;
import cn.harnesstemplate.admin.infrastructure.persistence.EntityProxies.JobScheduleProxy;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;

@Service
public class JobScheduleQueryService {

    private final EasyQuerySupport easyQuery;

    public JobScheduleQueryService(EasyQuerySupport easyQuery) { this.easyQuery = easyQuery; }

    public List<JobSchedule> page(int pageNum, int pageSize) {
        int offset = (pageNum - 1) * pageSize;
        return easyQuery.queryable(JobScheduleProxy.createTable())
                .where(schedule -> schedule.status().ne("DELETED"))
                .orderBy(schedule -> schedule.id().desc())
                .limit(offset, pageSize)
                .toList();
    }

    public long count() {
        return easyQuery.queryable(JobScheduleProxy.createTable())
                .where(schedule -> schedule.status().ne("DELETED"))
                .count();
    }

    public JobSchedule findById(Long id) {
        return easyQuery.queryable(JobScheduleProxy.createTable())
                .where(schedule -> {
                    schedule.id().eq(id);
                    schedule.status().ne("DELETED");
                })
                .firstOrNull();
    }

    public JobSchedule findByJobCode(String jobCode) {
        return easyQuery.queryable(JobScheduleProxy.createTable())
                .where(schedule -> {
                    schedule.jobCode().eq(jobCode);
                    schedule.status().ne("DELETED");
                })
                .firstOrNull();
    }

    public JobSchedule save(JobSchedule s) {
        return easyQuery.insert(s);
    }

    public void update(JobSchedule s) {
        easyQuery.expressionUpdatable(JobScheduleProxy.createTable())
                .setColumns(schedule -> {
                    schedule.jobName().set(s.getJobName());
                    schedule.workflowType().set(s.getWorkflowType());
                    schedule.taskQueue().set(s.getTaskQueue());
                    schedule.scheduleType().set(s.getScheduleType());
                    schedule.cronExpr().set(s.getCronExpr());
                    schedule.intervalSeconds().set(s.getIntervalSeconds());
                    schedule.startTime().set(s.getStartTime());
                    schedule.endTime().set(s.getEndTime());
                    schedule.inputJson().set(s.getInputJson());
                    schedule.status().set(s.getStatus());
                    schedule.temporalScheduleId().set(s.getTemporalScheduleId());
                    schedule.temporalWorkflowIdPrefix().set(s.getTemporalWorkflowIdPrefix());
                    schedule.description().set(s.getDescription());
                    schedule.updatedAt().set(LocalDateTime.now());
                })
                .where(schedule -> schedule.id().eq(s.getId()))
                .executeRows();
    }

    public void softDelete(Long id) {
        easyQuery.expressionUpdatable(JobScheduleProxy.createTable())
                .setColumns(schedule -> {
                    schedule.status().set("DELETED");
                    schedule.updatedAt().set(LocalDateTime.now());
                })
                .where(schedule -> schedule.id().eq(id))
                .executeRows();
    }

    public void updateStatus(Long id, String status) {
        easyQuery.expressionUpdatable(JobScheduleProxy.createTable())
                .setColumns(schedule -> {
                    schedule.status().set(status);
                    schedule.updatedAt().set(LocalDateTime.now());
                })
                .where(schedule -> schedule.id().eq(id))
                .executeRows();
    }
}
