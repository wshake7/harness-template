package cn.harnesstemplate.admin.web.controller;

import cn.harnesstemplate.admin.application.service.system.*;
import cn.harnesstemplate.admin.domain.entity.SysRole;
import cn.harnesstemplate.admin.web.dto.*;
import org.springframework.web.bind.annotation.*;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/sys/role")
public class RoleController {

    private final SysRoleQueryService roleService;
    private final SysUserRoleQueryService userRoleService;
    private final SysRoleMenuQueryService roleMenuService;
    private final SysRoleApiQueryService roleApiService;

    public RoleController(SysRoleQueryService roleService, SysUserRoleQueryService userRoleService,
                          SysRoleMenuQueryService roleMenuService, SysRoleApiQueryService roleApiService) {
        this.roleService = roleService;
        this.userRoleService = userRoleService;
        this.roleMenuService = roleMenuService;
        this.roleApiService = roleApiService;
    }

    @PostMapping("/list")
    public R<PageResult<SysRoleResponse>> list(@RequestBody Map<String, Object> req) {
        int page = req.containsKey("page") ? ((Number) req.get("page")).intValue() : 1;
        int pageSize = req.containsKey("pageSize") ? ((Number) req.get("pageSize")).intValue() : 10;
        boolean noPaging = req.containsKey("noPaging") && (Boolean) req.get("noPaging");

        if (noPaging) {
            List<SysRole> all = roleService.listAll();
            List<SysRoleResponse> tree = buildTree(all);
            return R.ok(PageResult.of(tree, all.size()));
        }

        String name = req.containsKey("name") ? (String) req.get("name") : null;
        List<SysRole> roles = roleService.page(page, pageSize, name);
        long total = roleService.count(name);
        List<SysRoleResponse> items = new ArrayList<>();
        for (SysRole r : roles) {
            items.add(toResponse(r));
        }
        return R.ok(PageResult.of(items, total));
    }

    @GetMapping("/tree")
    public R<List<SysRoleResponse>> tree() {
        List<SysRole> all = roleService.listAll();
        return R.ok(buildTree(all));
    }

    @GetMapping("/{id}/permissions")
    public R<Map<String, List<Long>>> permissions(@PathVariable Long id) {
        List<Long> menuIds = roleMenuService.findMenuIdsByRoleId(id);
        List<Long> apiIds = roleApiService.findApiIdsByRoleId(id);
        return R.ok(Map.of("menuIds", menuIds, "apiIds", apiIds));
    }

    @PostMapping("/create")
    public R<Void> create(@RequestBody SysRoleCreateRequest req) {
        SysRole existing = roleService.findByCode(req.code());
        if (existing != null) {
            return R.fail("角色编码已存在");
        }
        SysRole role = new SysRole();
        role.setName(req.name());
        role.setCode(req.code());
        role.setParentId(req.parentId());
        role.setIsEnabled(true);
        role.setRemark(req.remark());
        roleService.save(role);
        if (req.menuIds() != null) {
            for (Long menuId : req.menuIds()) {
                roleMenuService.save(role.getId(), menuId);
            }
        }
        if (req.apiIds() != null) {
            for (Long apiId : req.apiIds()) {
                roleApiService.save(role.getId(), apiId);
            }
        }
        return R.ok();
    }

    @PostMapping("/update")
    public R<Void> update(@RequestBody SysRoleUpdateRequest req) {
        SysRole role = roleService.findById(req.id());
        if (role == null) {
            return R.fail("角色不存在");
        }
        if (req.name() != null) role.setName(req.name());
        if (req.code() != null) role.setCode(req.code());
        if (req.parentId() != null) role.setParentId(req.parentId());
        if (req.isEnabled() != null) role.setIsEnabled(req.isEnabled());
        if (req.remark() != null) role.setRemark(req.remark());
        roleService.update(role);
        return R.ok();
    }

    @PostMapping("/del")
    public R<Void> del(@RequestBody Map<String, List<Long>> req) {
        List<Long> ids = req.get("ids");
        if (ids == null || ids.isEmpty()) {
            return R.fail("ids不能为空");
        }
        for (Long id : ids) {
            if (roleService.hasChildren(id)) {
                return R.fail("存在子角色，不能删除");
            }
            if (userRoleService.hasUsers(id)) {
                return R.fail("角色已绑定用户，不能删除");
            }
            if (roleMenuService.hasBindings(id)) {
                return R.fail("角色已绑定菜单权限，不能删除");
            }
            if (roleApiService.hasBindings(id)) {
                return R.fail("角色已绑定API权限，不能删除");
            }
        }
        for (Long id : ids) {
            roleService.delete(id);
        }
        return R.ok();
    }

    @PostMapping("/permissions")
    public R<Void> savePermissions(@RequestBody Map<String, Object> req) {
        Long id = ((Number) req.get("id")).longValue();
        @SuppressWarnings("unchecked")
        List<Integer> menuIdsRaw = (List<Integer>) req.get("menuIds");
        @SuppressWarnings("unchecked")
        List<Integer> apiIdsRaw = (List<Integer>) req.get("apiIds");

        roleMenuService.deleteByRoleId(id);
        roleApiService.deleteByRoleId(id);

        if (menuIdsRaw != null) {
            for (Integer mid : menuIdsRaw) {
                roleMenuService.save(id, mid.longValue());
            }
        }
        if (apiIdsRaw != null) {
            for (Integer aid : apiIdsRaw) {
                roleApiService.save(id, aid.longValue());
            }
        }
        return R.ok();
    }

    private List<SysRoleResponse> buildTree(List<SysRole> all) {
        List<SysRoleResponse> roots = new ArrayList<>();
        for (SysRole r : all) {
            if (r.getParentId() == null || r.getParentId() == 0) {
                roots.add(buildNode(r, all));
            }
        }
        return roots;
    }

    private SysRoleResponse buildNode(SysRole role, List<SysRole> all) {
        List<SysRoleResponse> children = new ArrayList<>();
        for (SysRole r : all) {
            if (r.getParentId() != null && r.getParentId().equals(role.getId())) {
                children.add(buildNode(r, all));
            }
        }
        return toResponse(role, children);
    }

    private SysRoleResponse toResponse(SysRole r) {
        return toResponse(r, null);
    }

    private SysRoleResponse toResponse(SysRole role, List<SysRoleResponse> children) {
        return new SysRoleResponse(
                role.getId(),
                role.getName(),
                role.getCode(),
                role.getParentId(),
                role.getRemark(),
                role.getIsEnabled(),
                null,
                null,
                children,
                role.getCreatedAt(),
                role.getUpdatedAt()
        );
    }
}
