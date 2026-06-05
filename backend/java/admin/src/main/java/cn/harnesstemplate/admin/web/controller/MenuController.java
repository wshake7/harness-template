package cn.harnesstemplate.admin.web.controller;

import cn.harnesstemplate.admin.application.service.system.SysResourceMenuApiQueryService;
import cn.harnesstemplate.admin.application.service.system.SysResourceMenuQueryService;
import cn.harnesstemplate.admin.domain.entity.SysResourceMenu;
import cn.harnesstemplate.admin.web.dto.*;
import org.springframework.web.bind.annotation.*;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/sys/resource/menu")
public class MenuController {

    private final SysResourceMenuQueryService menuService;
    private final SysResourceMenuApiQueryService menuApiService;

    public MenuController(SysResourceMenuQueryService menuService,
                          SysResourceMenuApiQueryService menuApiService) {
        this.menuService = menuService;
        this.menuApiService = menuApiService;
    }

    @PostMapping("/list")
    public R<PageResult<MenuResponse>> list(@RequestBody Map<String, Object> req) {
        boolean noPaging = req.containsKey("noPaging") && (Boolean) req.get("noPaging");
        List<SysResourceMenu> all = menuService.listAll();

        if (noPaging) {
            List<MenuResponse> tree = buildTree(all);
            return R.ok(PageResult.of(tree, all.size()));
        }

        // Non-tree: filter by top-level and return flat list
        List<MenuResponse> items = new ArrayList<>();
        for (SysResourceMenu m : all) {
            if (m.getParentId() == null || m.getParentId() == 0) {
                items.add(toResponse(m, null));
            }
        }
        return R.ok(PageResult.of(items, all.size()));
    }

    @GetMapping("/tree")
    public R<List<MenuResponse>> tree() {
        List<SysResourceMenu> all = menuService.listAll();
        return R.ok(buildTree(all));
    }

    @PostMapping("/create")
    public R<Void> create(@RequestBody Map<String, Object> req) {
        SysResourceMenu menu = new SysResourceMenu();
        if (req.containsKey("parentId") && req.get("parentId") != null)
            menu.setParentId(((Number) req.get("parentId")).longValue());
        menu.setMenuType((String) req.get("menuType"));
        menu.setPath((String) req.getOrDefault("path", ""));
        menu.setRedirect((String) req.getOrDefault("redirect", ""));
        menu.setAlias((String) req.getOrDefault("alias", ""));
        menu.setName((String) req.get("name"));
        menu.setComponent((String) req.getOrDefault("component", ""));
        menu.setMetadata(req.containsKey("metadata") ? req.get("metadata").toString() : "{}");
        menu.setSortOrder(req.containsKey("sortOrder") ? ((Number) req.get("sortOrder")).intValue() : 0);
        menu.setIsEnabled(req.containsKey("isEnabled") ? (Boolean) req.get("isEnabled") : true);
        menu.setRemark((String) req.getOrDefault("remark", ""));
        menu.setTreePath("");

        // Compute tree path from parent
        if (menu.getParentId() != null && menu.getParentId() > 0) {
            SysResourceMenu parent = menuService.findById(menu.getParentId());
            if (parent != null && parent.getTreePath() != null) {
                menu.setTreePath(parent.getTreePath() + menu.getParentId() + "/");
            }
        }

        menuService.save(menu);

        // Update tree path with own id
        menu.setTreePath((menu.getTreePath() != null ? menu.getTreePath() : "/") + menu.getId() + "/");
        menuService.update(menu);

        // Save API associations
        if (req.containsKey("apiIds") && req.get("apiIds") != null) {
            @SuppressWarnings("unchecked")
            List<Integer> apiIds = (List<Integer>) req.get("apiIds");
            for (Integer apiId : apiIds) {
                menuApiService.save(menu.getId(), apiId.longValue());
            }
        }

        return R.ok();
    }

    @PostMapping("/update")
    public R<Void> update(@RequestBody Map<String, Object> req) {
        Long id = ((Number) req.get("id")).longValue();
        SysResourceMenu menu = menuService.findById(id);
        if (menu == null) {
            return R.fail("菜单不存在");
        }

        if (req.containsKey("parentId")) menu.setParentId(req.get("parentId") != null ? ((Number) req.get("parentId")).longValue() : null);
        if (req.containsKey("menuType")) menu.setMenuType((String) req.get("menuType"));
        if (req.containsKey("path")) menu.setPath((String) req.get("path"));
        if (req.containsKey("redirect")) menu.setRedirect((String) req.get("redirect"));
        if (req.containsKey("alias")) menu.setAlias((String) req.get("alias"));
        if (req.containsKey("name")) menu.setName((String) req.get("name"));
        if (req.containsKey("component")) menu.setComponent((String) req.get("component"));
        if (req.containsKey("metadata")) menu.setMetadata(req.get("metadata").toString());
        if (req.containsKey("sortOrder")) menu.setSortOrder(((Number) req.get("sortOrder")).intValue());
        if (req.containsKey("isEnabled")) menu.setIsEnabled((Boolean) req.get("isEnabled"));
        if (req.containsKey("remark")) menu.setRemark((String) req.get("remark"));

        menuService.update(menu);

        // Update API associations if provided
        if (req.containsKey("apiIds") && req.get("apiIds") != null) {
            menuApiService.deleteByMenuId(id);
            @SuppressWarnings("unchecked")
            List<Integer> apiIds = (List<Integer>) req.get("apiIds");
            for (Integer apiId : apiIds) {
                menuApiService.save(id, apiId.longValue());
            }
        }

        return R.ok();
    }

    @PostMapping("/del")
    public R<Void> del(@RequestBody Map<String, List<Long>> req) {
        List<Long> ids = req.get("ids");
        if (ids == null || ids.isEmpty()) {
            return R.fail("ids不能为空");
        }
        for (Long id : ids) {
            if (menuService.hasChildren(id)) {
                return R.fail("存在子节点，不能删除");
            }
        }
        for (Long id : ids) {
            menuApiService.deleteByMenuId(id);
            menuService.delete(id);
        }
        return R.ok();
    }

    private List<MenuResponse> buildTree(List<SysResourceMenu> all) {
        List<MenuResponse> roots = new ArrayList<>();
        for (SysResourceMenu m : all) {
            if (m.getParentId() == null || m.getParentId() == 0) {
                roots.add(buildNode(m, all));
            }
        }
        return roots;
    }

    private MenuResponse buildNode(SysResourceMenu menu, List<SysResourceMenu> all) {
        List<MenuResponse> children = new ArrayList<>();
        for (SysResourceMenu m : all) {
            if (m.getParentId() != null && m.getParentId().equals(menu.getId())) {
                children.add(buildNode(m, all));
            }
        }
        return toResponse(menu, children);
    }

    private MenuResponse toResponse(SysResourceMenu m, List<MenuResponse> children) {
        return new MenuResponse(
                m.getId(),
                m.getMenuType(),
                m.getPath(),
                m.getRedirect(),
                m.getAlias(),
                m.getName(),
                m.getComponent(),
                m.getParentId(),
                m.getTreePath(),
                m.getSortOrder(),
                m.getIsEnabled(),
                m.getRemark(),
                children,
                menuApiService.findApiIdsByMenuId(m.getId()),
                m.getCreatedAt(),
                m.getUpdatedAt()
        );
    }
}
