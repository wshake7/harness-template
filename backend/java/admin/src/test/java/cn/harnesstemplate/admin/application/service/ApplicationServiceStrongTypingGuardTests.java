package cn.harnesstemplate.admin.application.service;

import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;
import java.util.regex.Pattern;

import static org.assertj.core.api.Assertions.assertThat;

class ApplicationServiceStrongTypingGuardTests {

    private static final Path SERVICE_ROOT = Path.of("src/main/java/cn/harnesstemplate/admin/application/service");
    private static final Pattern MAGIC_QUERY_FIELD = Pattern.compile(
            "\\.eq\\(\"[A-Za-z][A-Za-z0-9]*\"|\\.set\\(\"[A-Za-z][A-Za-z0-9]*\"\\s*,|\\.column\\(\"[A-Za-z][A-Za-z0-9]*\"");

    @Test
    void easyQueryShouldNotUseStringFieldNamesInApplicationServices() throws IOException {
        try (var paths = Files.walk(SERVICE_ROOT)) {
            List<String> violations = paths
                    .filter(path -> path.toString().endsWith(".java"))
                    .map(this::findViolation)
                    .filter(result -> result != null)
                    .toList();
            assertThat(violations)
                    .withFailMessage("Found Easy Query magic field names:%n%s", String.join(System.lineSeparator(), violations))
                    .isEmpty();
        }
    }

    private String findViolation(Path path) {
        try {
            int lineNumber = 0;
            for (String line : Files.readAllLines(path)) {
                lineNumber++;
                if (MAGIC_QUERY_FIELD.matcher(line).find()) {
                    return path + ":" + lineNumber + ": " + line.trim();
                }
            }
            return null;
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }
}
