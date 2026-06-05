package cn.harnesstemplate.admin.application.service.system;

import cn.harnesstemplate.admin.domain.entity.SysRoleApi;
import cn.harnesstemplate.admin.infrastructure.persistence.EasyQuerySupport;
import cn.harnesstemplate.admin.infrastructure.persistence.EntityProxies.SysRoleApiProxy;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;

@Service
public class SysRoleApiQueryService {

    private final EasyQuerySupport easyQuery;

    public SysRoleApiQueryService(EasyQuerySupport easyQuery) {
        this.easyQuery = easyQuery;
    }

    public List<Long> findApiIdsByRoleId(Long roleId) {
        return easyQuery.queryable(SysRoleApiProxy.createTable())
                .where(binding -> {
                    binding.roleId().eq(roleId);
                    binding.deletedAt().eq(0L);
                })
                .toList()
                .stream()
                .map(SysRoleApi::getApiId)
                .toList();
    }

    public boolean hasBindings(Long roleId) {
        return easyQuery.queryable(SysRoleApiProxy.createTable())
                .where(binding -> {
                    binding.roleId().eq(roleId);
                    binding.deletedAt().eq(0L);
                })
                .count() > 0;
    }

    public void deleteByRoleId(Long roleId) {
        easyQuery.expressionUpdatable(SysRoleApiProxy.createTable())
                .setColumns(binding -> binding.deletedAt().set(System.currentTimeMillis()))
                .where(binding -> {
                    binding.roleId().eq(roleId);
                    binding.deletedAt().eq(0L);
                })
                .executeRows();
    }

    public void save(Long roleId, Long apiId) {
        SysRoleApi binding = new SysRoleApi();
        binding.setRoleId(roleId);
        binding.setApiId(apiId);
        binding.setCreatedAt(LocalDateTime.now());
        binding.setUpdatedAt(LocalDateTime.now());
        easyQuery.insert(binding);
    }
}
