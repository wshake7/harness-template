package cn.harnesstemplate.admin.ai;

import lombok.Data;

import java.util.List;

@Data
public class RagResponse {
    private String prompt;
    private int retrievedCount;
    private String augmentedPrompt;
    private String answer;
    private List<RagSource> sources;
}
