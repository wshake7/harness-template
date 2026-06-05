package cn.harnesstemplate.admin.web.controller;

import cn.harnesstemplate.admin.application.service.job.JobExecutionQueryService;
import cn.harnesstemplate.admin.application.service.job.JobScheduleQueryService;
import cn.harnesstemplate.admin.application.service.job.TemporalService;
import cn.harnesstemplate.admin.domain.entity.JobExecution;
import cn.harnesstemplate.admin.domain.entity.JobSchedule;
import cn.harnesstemplate.admin.web.dto.PageResult;
import cn.harnesstemplate.admin.web.dto.R;
import org.springframework.web.bind.annotation.*;

import java.util.*;

@RestController
@RequestMapping("/api/sys/job/execution")
public class JobExecutionController {

    private final JobExecutionQueryService service;
    private final JobScheduleQueryService scheduleService;
    private final TemporalService temporalService;

    public JobExecutionController(JobExecutionQueryService service, JobScheduleQueryService scheduleService,
                                  TemporalService temporalService) {
        this.service = service;
        this.scheduleService = scheduleService;
        this.temporalService = temporalService;
    }

    @PostMapping("/list")
    public R<PageResult<Map<String, Object>>> list(@RequestBody Map<String, Object> req) {
        int page = req.containsKey("page") ? ((Number) req.get("page")).intValue() : 1;
        int pageSize = req.containsKey("pageSize") ? ((Number) req.get("pageSize")).intValue() : 10;
        String jobCode = (String) req.getOrDefault("jobCode", null);
        List<JobExecution> executions = service.page(page, pageSize, jobCode);
        long total = service.count(jobCode);
        List<Map<String, Object>> items = new ArrayList<>();
        for (JobExecution e : executions) items.add(toMap(e));
        return R.ok(PageResult.of(items, total));
    }

    @PostMapping("/detail")
    public R<Map<String, Object>> detail(@RequestBody Map<String, Long> req) {
        JobExecution e = service.findById(req.get("id"));
        if (e == null) return R.fail("执行记录不存在");
        return R.ok(toMap(e));
    }

    @PostMapping("/cancel")
    public R<Void> cancel(@RequestBody Map<String, Long> req) {
        Long id = req.get("id");
        JobExecution e = service.findById(id);
        if (e == null) return R.fail("执行记录不存在");
        if (!"RUNNING".equals(e.getStatus())) return R.fail("只能取消运行中的任务");
        temporalService.cancelExecution(e);
        return R.ok();
    }

    @PostMapping("/retry")
    public R<Void> retry(@RequestBody Map<String, Long> req) {
        Long id = req.get("id");
        JobExecution e = service.findById(id);
        if (e == null) return R.fail("执行记录不存在");
        if (!List.of("FAILED", "CANCELED", "TIMEOUT").contains(e.getStatus())) {
            return R.fail("只能重试失败、取消或超时的任务");
        }
        JobSchedule schedule = scheduleService.findByJobCode(e.getJobCode());
        if (schedule == null) return R.fail("关联的任务调度不存在");
        temporalService.triggerSchedule(schedule);
        return R.ok();
    }

    private Map<String, Object> toMap(JobExecution e) {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("id", e.getId());
        m.put("jobCode", e.getJobCode());
        m.put("temporalWorkflowId", e.getTemporalWorkflowId());
        m.put("temporalRunId", e.getTemporalRunId());
        m.put("triggerTime", e.getTriggerTime());
        m.put("startTime", e.getStartTime());
        m.put("endTime", e.getEndTime());
        m.put("status", e.getStatus());
        m.put("inputJson", e.getInputJson());
        m.put("resultJson", e.getResultJson());
        m.put("errorMessage", e.getErrorMessage());
        m.put("retryCount", e.getRetryCount());
        m.put("createdAt", e.getCreatedAt());
        m.put("updatedAt", e.getUpdatedAt());
        return m;
    }
}
