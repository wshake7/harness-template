package cn.harnesstemplate.admin.application.service.job;

import cn.harnesstemplate.admin.domain.entity.JobExecution;
import cn.harnesstemplate.admin.domain.entity.JobSchedule;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.jdbc.AutoConfigureTestDatabase;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import static org.assertj.core.api.Assertions.assertThat;

@SpringBootTest
@ActiveProfiles("test-db")
@Testcontainers
@AutoConfigureTestDatabase(replace = AutoConfigureTestDatabase.Replace.NONE)
class JobScheduleQueryServiceTests {

    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:16-alpine")
            .withDatabaseName("testdb")
            .withUsername("test")
            .withPassword("test");

    @DynamicPropertySource
    static void configureProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", postgres::getJdbcUrl);
        registry.add("spring.datasource.username", postgres::getUsername);
        registry.add("spring.datasource.password", postgres::getPassword);
        registry.add("spring.datasource.driver-class-name", () -> "org.postgresql.Driver");
    }

    @Autowired
    private JobScheduleQueryService scheduleService;

    @Autowired
    private JobExecutionQueryService executionService;

    @BeforeEach
    void setUp() {
        var all = scheduleService.page(1, 100);
        for (var s : all) scheduleService.softDelete(s.getId());
    }

    @Test
    void shouldCreateAndFindSchedule() {
        JobSchedule s = new JobSchedule();
        s.setJobCode("test-job");
        s.setJobName("Test Job");
        s.setWorkflowType("PrintCountWorkflow");
        s.setTaskQueue("admin-java");
        s.setScheduleType("CRON");
        s.setCronExpr("0 */5 * * *");
        s.setStatus("ENABLED");
        scheduleService.save(s);

        assertThat(s.getId()).isNotNull();

        JobSchedule found = scheduleService.findById(s.getId());
        assertThat(found.getJobCode()).isEqualTo("test-job");
    }

    @Test
    void shouldSoftDeleteSchedule() {
        JobSchedule s = new JobSchedule();
        s.setJobCode("to-delete");
        s.setJobName("To Delete");
        s.setWorkflowType("PrintCountWorkflow");
        s.setTaskQueue("admin-java");
        s.setScheduleType("ONCE");
        s.setStatus("ENABLED");
        scheduleService.save(s);

        scheduleService.softDelete(s.getId());
        assertThat(scheduleService.findById(s.getId())).isNull();
    }

    @Test
    void shouldCreateExecution() {
        JobExecution e = new JobExecution();
        e.setJobCode("test-job");
        e.setTemporalWorkflowId("test-job-dispatch-001");
        e.setTemporalRunId("");
        e.setStatus("RUNNING");
        e.setRetryCount(0);
        executionService.save(e);

        assertThat(e.getId()).isNotNull();
        assertThat(e.getStatus()).isEqualTo("RUNNING");
    }

    @Test
    void shouldUpdateExecutionStatus() {
        JobExecution e = new JobExecution();
        e.setJobCode("test-job");
        e.setTemporalWorkflowId("test-job-001");
        e.setTemporalRunId("");
        e.setStatus("RUNNING");
        e.setRetryCount(0);
        executionService.save(e);

        executionService.updateStatus(e.getId(), "SUCCESS", null, "{}");
        JobExecution found = executionService.findById(e.getId());
        assertThat(found.getStatus()).isEqualTo("SUCCESS");
    }
}
