package cn.harnesstemplate.admin.web.controller;

import cn.harnesstemplate.admin.application.service.job.JobScheduleQueryService;
import cn.harnesstemplate.admin.application.service.job.TemporalService;
import cn.harnesstemplate.admin.domain.entity.JobSchedule;
import cn.harnesstemplate.admin.web.dto.PageResult;
import cn.harnesstemplate.admin.web.dto.R;
import org.springframework.web.bind.annotation.*;

import java.util.*;

@RestController
@RequestMapping("/api/sys/job/schedule")
public class JobScheduleController {

    private final JobScheduleQueryService service;
    private final TemporalService temporalService;

    public JobScheduleController(JobScheduleQueryService service, TemporalService temporalService) {
        this.service = service;
        this.temporalService = temporalService;
    }

    @PostMapping("/list")
    public R<PageResult<Map<String, Object>>> list(@RequestBody Map<String, Object> req) {
        int page = req.containsKey("page") ? ((Number) req.get("page")).intValue() : 1;
        int pageSize = req.containsKey("pageSize") ? ((Number) req.get("pageSize")).intValue() : 10;
        List<JobSchedule> schedules = service.page(page, pageSize);
        long total = service.count();
        List<Map<String, Object>> items = new ArrayList<>();
        for (JobSchedule s : schedules) items.add(toMap(s));
        return R.ok(PageResult.of(items, total));
    }

    @PostMapping("/detail")
    public R<Map<String, Object>> detail(@RequestBody Map<String, Long> req) {
        JobSchedule s = service.findById(req.get("id"));
        if (s == null) return R.fail("任务不存在");
        return R.ok(toMap(s));
    }

    @PostMapping("/create")
    public R<Void> create(@RequestBody Map<String, Object> req) {
        JobSchedule s = new JobSchedule();
        s.setJobCode((String) req.get("jobCode"));
        s.setJobName((String) req.get("jobName"));
        s.setWorkflowType((String) req.get("workflowType"));
        s.setTaskQueue((String) req.getOrDefault("taskQueue", "admin-java"));
        s.setScheduleType((String) req.get("scheduleType"));
        s.setCronExpr((String) req.getOrDefault("cronExpr", null));
        s.setIntervalSeconds(req.containsKey("intervalSeconds") ? ((Number) req.get("intervalSeconds")).intValue() : null);
        s.setStatus((String) req.getOrDefault("status", "ENABLED"));
        s.setTemporalScheduleId((String) req.getOrDefault("temporalScheduleId", null));
        s.setTemporalWorkflowIdPrefix((String) req.getOrDefault("temporalWorkflowIdPrefix", null));
        s.setInputJson(req.containsKey("inputJson") ? req.get("inputJson").toString() : "{}");
        s.setDescription((String) req.getOrDefault("description", ""));
        service.save(s);
        if ("ENABLED".equals(s.getStatus())) {
            temporalService.syncSchedule(s);
        }
        return R.ok();
    }

    @PostMapping("/update")
    public R<Void> update(@RequestBody Map<String, Object> req) {
        Long id = ((Number) req.get("id")).longValue();
        JobSchedule s = service.findById(id);
        if (s == null) return R.fail("任务不存在");
        if (req.containsKey("jobName")) s.setJobName((String) req.get("jobName"));
        if (req.containsKey("workflowType")) s.setWorkflowType((String) req.get("workflowType"));
        if (req.containsKey("taskQueue")) s.setTaskQueue((String) req.get("taskQueue"));
        if (req.containsKey("scheduleType")) s.setScheduleType((String) req.get("scheduleType"));
        if (req.containsKey("cronExpr")) s.setCronExpr((String) req.get("cronExpr"));
        if (req.containsKey("intervalSeconds")) s.setIntervalSeconds(((Number) req.get("intervalSeconds")).intValue());
        if (req.containsKey("status")) s.setStatus((String) req.get("status"));
        if (req.containsKey("temporalScheduleId")) s.setTemporalScheduleId((String) req.get("temporalScheduleId"));
        if (req.containsKey("temporalWorkflowIdPrefix")) s.setTemporalWorkflowIdPrefix((String) req.get("temporalWorkflowIdPrefix"));
        if (req.containsKey("inputJson")) s.setInputJson(req.get("inputJson").toString());
        if (req.containsKey("description")) s.setDescription((String) req.get("description"));
        service.update(s);
        temporalService.syncSchedule(s);
        return R.ok();
    }

    @PostMapping("/del")
    public R<Void> del(@RequestBody Map<String, Long> req) {
        Long id = req.get("id");
        JobSchedule s = service.findById(id);
        if (s != null) {
            temporalService.deleteSchedule(s);
            service.softDelete(id);
        }
        return R.ok();
    }

    @PostMapping("/switch")
    public R<Void> switchStatus(@RequestBody Map<String, Object> req) {
        Long id = ((Number) req.get("id")).longValue();
        boolean enabled = (Boolean) req.get("enabled");
        JobSchedule s = service.findById(id);
        if (s == null) return R.fail("任务不存在");
        service.updateStatus(id, enabled ? "ENABLED" : "DISABLED");
        if (enabled) {
            temporalService.unpauseSchedule(s);
            temporalService.syncSchedule(s);
        } else {
            temporalService.pauseSchedule(s);
        }
        return R.ok();
    }

    @PostMapping("/sync")
    public R<Void> sync(@RequestBody Map<String, Long> req) {
        Long id = req.get("id");
        JobSchedule s = service.findById(id);
        if (s == null) return R.fail("任务不存在");
        temporalService.deleteSchedule(s);
        temporalService.syncSchedule(s);
        return R.ok();
    }

    @PostMapping("/trigger")
    public R<Void> trigger(@RequestBody Map<String, Long> req) {
        Long id = req.get("id");
        JobSchedule s = service.findById(id);
        if (s == null) return R.fail("任务不存在");
        temporalService.triggerSchedule(s);
        return R.ok();
    }

    @PostMapping("/options")
    public R<Map<String, Object>> options() {
        Map<String, Object> result = new LinkedHashMap<>();
        result.put("workflowTypes", Map.of(
                "PrintCountWorkflow", "PrintCountWorkflow",
                "KnowledgeDocumentIndexWorkflow", "KnowledgeDocumentIndexWorkflow"));
        result.put("taskQueues", List.of("admin-java", "admin"));
        result.put("defaultTaskQueue", "admin-java");
        return R.ok(result);
    }

    private Map<String, Object> toMap(JobSchedule s) {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("id", s.getId());
        m.put("jobCode", s.getJobCode());
        m.put("jobName", s.getJobName());
        m.put("workflowType", s.getWorkflowType());
        m.put("taskQueue", s.getTaskQueue());
        m.put("scheduleType", s.getScheduleType());
        m.put("cronExpr", s.getCronExpr());
        m.put("intervalSeconds", s.getIntervalSeconds());
        m.put("startTime", s.getStartTime());
        m.put("endTime", s.getEndTime());
        m.put("inputJson", s.getInputJson());
        m.put("status", s.getStatus());
        m.put("temporalScheduleId", s.getTemporalScheduleId());
        m.put("temporalWorkflowIdPrefix", s.getTemporalWorkflowIdPrefix());
        m.put("description", s.getDescription());
        m.put("createdAt", s.getCreatedAt());
        m.put("updatedAt", s.getUpdatedAt());
        return m;
    }
}
