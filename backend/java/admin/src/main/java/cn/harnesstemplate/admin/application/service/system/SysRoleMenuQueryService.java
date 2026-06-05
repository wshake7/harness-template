package cn.harnesstemplate.admin.application.service.system;

import cn.harnesstemplate.admin.domain.entity.SysRoleMenu;
import cn.harnesstemplate.admin.infrastructure.persistence.EasyQuerySupport;
import cn.harnesstemplate.admin.infrastructure.persistence.EntityProxies.SysRoleMenuProxy;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;

@Service
public class SysRoleMenuQueryService {

    private final EasyQuerySupport easyQuery;

    public SysRoleMenuQueryService(EasyQuerySupport easyQuery) {
        this.easyQuery = easyQuery;
    }

    public List<Long> findMenuIdsByRoleId(Long roleId) {
        return easyQuery.queryable(SysRoleMenuProxy.createTable())
                .where(binding -> {
                    binding.roleId().eq(roleId);
                    binding.deletedAt().eq(0L);
                })
                .toList()
                .stream()
                .map(SysRoleMenu::getMenuId)
                .toList();
    }

    public boolean hasBindings(Long roleId) {
        return easyQuery.queryable(SysRoleMenuProxy.createTable())
                .where(binding -> {
                    binding.roleId().eq(roleId);
                    binding.deletedAt().eq(0L);
                })
                .count() > 0;
    }

    public void deleteByRoleId(Long roleId) {
        easyQuery.expressionUpdatable(SysRoleMenuProxy.createTable())
                .setColumns(binding -> binding.deletedAt().set(System.currentTimeMillis()))
                .where(binding -> {
                    binding.roleId().eq(roleId);
                    binding.deletedAt().eq(0L);
                })
                .executeRows();
    }

    public void save(Long roleId, Long menuId) {
        SysRoleMenu binding = new SysRoleMenu();
        binding.setRoleId(roleId);
        binding.setMenuId(menuId);
        binding.setCreatedAt(LocalDateTime.now());
        binding.setUpdatedAt(LocalDateTime.now());
        easyQuery.insert(binding);
    }
}
