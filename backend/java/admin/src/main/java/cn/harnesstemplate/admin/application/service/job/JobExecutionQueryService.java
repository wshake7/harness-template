package cn.harnesstemplate.admin.application.service.job;

import cn.harnesstemplate.admin.domain.entity.JobExecution;
import cn.harnesstemplate.admin.infrastructure.persistence.EasyQuerySupport;
import cn.harnesstemplate.admin.infrastructure.persistence.EntityProxies.JobExecutionProxy;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;

@Service
public class JobExecutionQueryService {

    private final EasyQuerySupport easyQuery;

    public JobExecutionQueryService(EasyQuerySupport easyQuery) { this.easyQuery = easyQuery; }

    public List<JobExecution> page(int pageNum, int pageSize, String jobCode) {
        int offset = (pageNum - 1) * pageSize;
        return easyQuery.queryable(JobExecutionProxy.createTable())
                .where(job -> job.jobCode().eq(jobCode != null && !jobCode.isEmpty(), jobCode))
                .orderBy(job -> job.id().desc())
                .limit(offset, pageSize)
                .toList();
    }

    public long count(String jobCode) {
        return easyQuery.queryable(JobExecutionProxy.createTable())
                .where(job -> job.jobCode().eq(jobCode != null && !jobCode.isEmpty(), jobCode))
                .count();
    }

    public JobExecution findById(Long id) {
        return easyQuery.queryable(JobExecutionProxy.createTable())
                .where(job -> job.id().eq(id))
                .firstOrNull();
    }

    public JobExecution save(JobExecution e) {
        if (e.getTriggerTime() == null) {
            e.setTriggerTime(LocalDateTime.now());
        }
        return easyQuery.insert(e);
    }

    public void updateStatus(Long id, String status, String errorMessage, String resultJson) {
        easyQuery.expressionUpdatable(JobExecutionProxy.createTable())
                .setColumns(job -> {
                    job.status().set(status);
                    job.errorMessage().set(errorMessage);
                    job.resultJson().set(resultJson);
                    job.endTime().set(LocalDateTime.now());
                    job.updatedAt().set(LocalDateTime.now());
                })
                .where(job -> job.id().eq(id))
                .executeRows();
    }

    public void updateCanceled(Long id) {
        easyQuery.expressionUpdatable(JobExecutionProxy.createTable())
                .setColumns(job -> {
                    job.status().set("CANCELED");
                    job.endTime().set(LocalDateTime.now());
                    job.updatedAt().set(LocalDateTime.now());
                })
                .where(job -> job.id().eq(id))
                .executeRows();
    }
}
