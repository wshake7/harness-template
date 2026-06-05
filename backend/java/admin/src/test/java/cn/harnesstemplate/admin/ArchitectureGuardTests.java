package cn.harnesstemplate.admin;

import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;
import java.util.stream.Stream;

import static org.assertj.core.api.Assertions.assertThat;

class ArchitectureGuardTests {

    private final Path projectRoot = Path.of("").toAbsolutePath();

    @Test
    void mainCodeShouldNotUseJdbcTemplateForBusinessPersistence() throws IOException {
        Path mainJava = projectRoot.resolve("src/main/java");

        List<Path> jdbcTemplateUsers;
        try (Stream<Path> files = Files.walk(mainJava)) {
            jdbcTemplateUsers = files
                    .filter(path -> path.toString().endsWith(".java"))
                    .filter(path -> contains(path, "JdbcTemplate") || contains(path, "NamedParameterJdbcTemplate"))
                    .toList();
        }

        assertThat(jdbcTemplateUsers)
                .as("Business persistence must use easy-query, not Spring JdbcTemplate")
                .isEmpty();
    }

    @Test
    void flywayShouldBaselineExistingGoSchemaOnFirstJavaStartup() throws IOException {
        String applicationYaml = Files.readString(projectRoot.resolve("src/main/resources/application.yml"));

        assertThat(applicationYaml)
                .contains("baseline-on-migrate: true")
                .contains("baseline-version: 1");
    }

    @Test
    void aiStackShouldUseLangChain4j() throws IOException {
        String buildGradle = Files.readString(projectRoot.resolve("build.gradle"));
        String ragChatService = Files.readString(projectRoot.resolve("src/main/java/cn/harnesstemplate/admin/ai/RagChatService.java"));

        assertThat(buildGradle).contains("dev.langchain4j");
        assertThat(ragChatService).contains("dev.langchain4j");
    }

    private boolean contains(Path path, String needle) {
        try {
            return Files.readString(path).contains(needle);
        } catch (IOException e) {
            throw new IllegalStateException("Failed to read " + path, e);
        }
    }
}
