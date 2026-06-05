package cn.harnesstemplate.admin.ai;

import lombok.Data;

import java.util.List;

@Data
public class RagRequest {
    private String collectionName;
    private String prompt;
    private int topK = 5;
    private int maxTokens = 2048;
}
