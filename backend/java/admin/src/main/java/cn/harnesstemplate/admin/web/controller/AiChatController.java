package cn.harnesstemplate.admin.web.controller;

import cn.harnesstemplate.admin.ai.RagChatService;
import cn.harnesstemplate.admin.ai.RagRequest;
import cn.harnesstemplate.admin.web.dto.R;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/ai/chat")
public class AiChatController {

    private final RagChatService ragChatService;

    public AiChatController(RagChatService ragChatService) {
        this.ragChatService = ragChatService;
    }

    @PostMapping("/rag")
    public R<?> chat(@RequestBody RagRequest request) {
        if (request.getCollectionName() == null || request.getCollectionName().isEmpty()) {
            return R.fail("知识库名称不能为空");
        }
        if (request.getPrompt() == null || request.getPrompt().isEmpty()) {
            return R.fail("提示词不能为空");
        }
        var response = ragChatService.chat(request);
        return R.ok(response);
    }
}
