package cn.harnesstemplate.admin.application.service.system;

import cn.harnesstemplate.admin.domain.entity.SysResourceMenuApi;
import cn.harnesstemplate.admin.infrastructure.persistence.EasyQuerySupport;
import cn.harnesstemplate.admin.infrastructure.persistence.EntityProxies.SysResourceMenuApiProxy;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;

@Service
public class SysResourceMenuApiQueryService {

    private final EasyQuerySupport easyQuery;

    public SysResourceMenuApiQueryService(EasyQuerySupport easyQuery) {
        this.easyQuery = easyQuery;
    }

    public List<Long> findApiIdsByMenuId(Long menuId) {
        return easyQuery.queryable(SysResourceMenuApiProxy.createTable())
                .where(binding -> {
                    binding.menuId().eq(menuId);
                    binding.deletedAt().eq(0L);
                })
                .toList()
                .stream()
                .map(SysResourceMenuApi::getApiId)
                .toList();
    }

    public void deleteByMenuId(Long menuId) {
        easyQuery.expressionUpdatable(SysResourceMenuApiProxy.createTable())
                .setColumns(binding -> binding.deletedAt().set(System.currentTimeMillis()))
                .where(binding -> {
                    binding.menuId().eq(menuId);
                    binding.deletedAt().eq(0L);
                })
                .executeRows();
    }

    public void save(Long menuId, Long apiId) {
        SysResourceMenuApi binding = new SysResourceMenuApi();
        binding.setMenuId(menuId);
        binding.setApiId(apiId);
        binding.setCreatedAt(LocalDateTime.now());
        binding.setUpdatedAt(LocalDateTime.now());
        easyQuery.insert(binding);
    }
}
