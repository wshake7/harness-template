package cn.harnesstemplate.admin.application.service.system;

import cn.harnesstemplate.admin.domain.entity.SysUserRole;
import cn.harnesstemplate.admin.infrastructure.persistence.EasyQuerySupport;
import cn.harnesstemplate.admin.infrastructure.persistence.EntityProxies.SysUserRoleProxy;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;

@Service
public class SysUserRoleQueryService {

    private final EasyQuerySupport easyQuery;

    public SysUserRoleQueryService(EasyQuerySupport easyQuery) {
        this.easyQuery = easyQuery;
    }

    public List<SysUserRole> findByUserId(Long userId) {
        return easyQuery.queryable(SysUserRoleProxy.createTable())
                .where(binding -> {
                    binding.userId().eq(userId);
                    binding.deletedAt().eq(0L);
                })
                .toList();
    }

    public List<SysUserRole> findByRoleId(Long roleId) {
        return easyQuery.queryable(SysUserRoleProxy.createTable())
                .where(binding -> {
                    binding.roleId().eq(roleId);
                    binding.deletedAt().eq(0L);
                })
                .toList();
    }

    public void save(Long userId, Long roleId) {
        SysUserRole binding = new SysUserRole();
        binding.setUserId(userId);
        binding.setRoleId(roleId);
        binding.setCreatedAt(LocalDateTime.now());
        binding.setUpdatedAt(LocalDateTime.now());
        easyQuery.insert(binding);
    }

    public void deleteByUserId(Long userId) {
        easyQuery.expressionUpdatable(SysUserRoleProxy.createTable())
                .setColumns(binding -> binding.deletedAt().set(System.currentTimeMillis()))
                .where(binding -> {
                    binding.userId().eq(userId);
                    binding.deletedAt().eq(0L);
                })
                .executeRows();
    }

    public void deleteByRoleId(Long roleId) {
        easyQuery.expressionUpdatable(SysUserRoleProxy.createTable())
                .setColumns(binding -> binding.deletedAt().set(System.currentTimeMillis()))
                .where(binding -> {
                    binding.roleId().eq(roleId);
                    binding.deletedAt().eq(0L);
                })
                .executeRows();
    }

    public boolean hasUsers(Long roleId) {
        return easyQuery.queryable(SysUserRoleProxy.createTable())
                .where(binding -> {
                    binding.roleId().eq(roleId);
                    binding.deletedAt().eq(0L);
                })
                .count() > 0;
    }
}
