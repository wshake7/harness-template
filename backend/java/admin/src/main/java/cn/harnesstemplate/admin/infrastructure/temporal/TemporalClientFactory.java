package cn.harnesstemplate.admin.infrastructure.temporal;

import cn.harnesstemplate.admin.config.AdminProperties;
import io.temporal.client.WorkflowClient;
import io.temporal.client.WorkflowClientOptions;
import io.temporal.client.schedules.ScheduleClient;
import io.temporal.client.schedules.ScheduleClientOptions;
import io.temporal.serviceclient.WorkflowServiceStubs;
import io.temporal.serviceclient.WorkflowServiceStubsOptions;
import io.temporal.worker.Worker;
import io.temporal.worker.WorkerFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;

public class TemporalClientFactory {

    private static final Logger log = LoggerFactory.getLogger(TemporalClientFactory.class);

    private final AdminProperties.TemporalConfig config;
    private WorkflowServiceStubs service;
    private WorkflowClient workflowClient;
    private ScheduleClient scheduleClient;
    private WorkerFactory workerFactory;
    private Worker worker;

    public TemporalClientFactory(AdminProperties.TemporalConfig config) {
        this.config = config;
    }

    public boolean isEnabled() {
        return config.isEnabled();
    }

    public void start() {
        if (!config.isEnabled()) {
            log.info("Temporal is disabled, skipping client initialization");
            return;
        }

        var stubsOptions = WorkflowServiceStubsOptions.newBuilder()
                .setTarget(config.getHostPort())
                .build();
        service = WorkflowServiceStubs.newServiceStubs(stubsOptions);

        var clientOptions = WorkflowClientOptions.newBuilder()
                .setNamespace(config.getNamespace())
                .setIdentity(config.getIdentity() != null ? config.getIdentity() : "")
                .build();
        workflowClient = WorkflowClient.newInstance(service, clientOptions);

        var scheduleOptions = ScheduleClientOptions.newBuilder()
                .setNamespace(config.getNamespace())
                .build();
        scheduleClient = ScheduleClient.newInstance(service, scheduleOptions);

        if (config.isWorkerEnabled()) {
            workerFactory = WorkerFactory.newInstance(workflowClient);
            worker = workerFactory.newWorker(config.getTaskQueue());
            workerFactory.start();
            log.info("Temporal worker started on task queue: {}", config.getTaskQueue());
        }

        log.info("Temporal client connected to {} (namespace: {})", config.getHostPort(), config.getNamespace());
    }

    public WorkflowClient getWorkflowClient() {
        return workflowClient;
    }

    public ScheduleClient getScheduleClient() {
        return scheduleClient;
    }

    public Worker getWorker() {
        return worker;
    }

    public void registerWorkflow(Class<?> workflowClass) {
        if (worker != null) {
            worker.registerWorkflowImplementationTypes(workflowClass);
        }
    }

    public void registerActivities(Object activities) {
        if (worker != null) {
            worker.registerActivitiesImplementations(activities);
        }
    }

    public void close() {
        if (workerFactory != null) {
            try { workerFactory.shutdown(); } catch (Exception ignored) {}
        }
        if (workflowClient != null) {
            try {
                workflowClient.getWorkflowServiceStubs().shutdown();
            } catch (Exception ignored) {}
        }
    }
}
