package cn.harnesstemplate.admin.infrastructure.persistence;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.jdbc.AutoConfigureTestDatabase;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.context.ApplicationContext;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import javax.sql.DataSource;
import java.sql.Connection;
import java.sql.DatabaseMetaData;
import java.sql.ResultSet;
import java.util.HashSet;
import java.util.Set;

import static org.assertj.core.api.Assertions.assertThat;

@SpringBootTest
@ActiveProfiles("test")
@Testcontainers
@AutoConfigureTestDatabase(replace = AutoConfigureTestDatabase.Replace.NONE)
class SchemaBaselineTests {

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
        registry.add("spring.flyway.enabled", () -> "true");
        registry.add("spring.flyway.locations", () -> "classpath:db/migration");
        registry.add("easy-query.enable", () -> "true");
        registry.add("easy-query.database", () -> "pgsql");
    }

    @Autowired
    private DataSource dataSource;

    @Autowired
    private ApplicationContext context;

    @Test
    void contextShouldLoad() {
        assertThat(context).isNotNull();
    }

    @Test
    void flywayShouldCreateAllExpectedTables() throws Exception {
        Set<String> expected = Set.of(
                "sys_user", "sys_role", "sys_user_role",
                "sys_resource_menu", "sys_resource_api", "sys_resource_menu_api",
                "sys_role_menu", "sys_role_api",
                "sys_dict_type", "sys_dict_entry",
                "sys_language_type", "sys_language_entry",
                "sys_casbin_model", "casbin_rule",
                "sys_data_permission",
                "job_schedule", "job_execution",
                "file_asset",
                "knowledge_collection", "knowledge_document",
                "agent_session", "agent_message",
                "sys_api_log", "sys_login_log"
        );

        Set<String> actual = new HashSet<>();
        try (Connection conn = dataSource.getConnection()) {
            DatabaseMetaData meta = conn.getMetaData();
            ResultSet tables = meta.getTables(null, "public", "%", new String[]{"TABLE"});
            while (tables.next()) {
                String tableName = tables.getString("TABLE_NAME");
                if (!tableName.startsWith("flyway_")) {
                    actual.add(tableName);
                }
            }
        }

        assertThat(actual).containsAll(expected);
    }

    @Test
    void sysUserTableShouldHaveCorrectColumns() throws Exception {
        Set<String> expectedColumns = Set.of(
                "id", "created_at", "updated_at",
                "created_by", "updated_by", "deleted_by",
                "remark", "is_enabled", "deleted_at",
                "username", "nickname", "password", "language_code"
        );

        Set<String> actual = new HashSet<>();
        try (Connection conn = dataSource.getConnection()) {
            DatabaseMetaData meta = conn.getMetaData();
            ResultSet columns = meta.getColumns(null, "public", "sys_user", null);
            while (columns.next()) {
                actual.add(columns.getString("COLUMN_NAME"));
            }
        }

        assertThat(actual).containsAll(expectedColumns);
    }

    @Test
    void flywayMigrationShouldHaveExpectedCount() throws Exception {
        JdbcTemplate jdbc = new JdbcTemplate(dataSource);
        Integer count = jdbc.queryForObject(
                "SELECT COUNT(*) FROM flyway_schema_history WHERE success = true", Integer.class);
        assertThat(count).isNotNull().isGreaterThan(0);
    }
}
