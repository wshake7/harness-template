package cn.harnesstemplate.admin.application.service.system;

import cn.harnesstemplate.admin.domain.entity.SysResourceApi;
import cn.harnesstemplate.admin.infrastructure.persistence.EasyQuerySupport;
import cn.harnesstemplate.admin.infrastructure.persistence.EntityProxies.SysResourceApiProxy;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;

@Service
public class SysResourceApiQueryService {

    private final EasyQuerySupport easyQuery;

    public SysResourceApiQueryService(EasyQuerySupport easyQuery) {
        this.easyQuery = easyQuery;
    }

    public List<SysResourceApi> listAll() {
        return easyQuery.queryable(SysResourceApiProxy.createTable())
                .where(api -> api.deletedAt().eq(0L))
                .orderBy(api -> {
                    api.sortOrder().asc();
                    api.id().asc();
                })
                .toList();
    }

    public SysResourceApi findById(Long id) {
        return easyQuery.queryable(SysResourceApiProxy.createTable())
                .where(api -> {
                    api.id().eq(id);
                    api.deletedAt().eq(0L);
                })
                .firstOrNull();
    }

    public SysResourceApi save(SysResourceApi api) {
        return easyQuery.insert(api);
    }

    public void update(SysResourceApi api) {
        easyQuery.expressionUpdatable(SysResourceApiProxy.createTable())
                .setColumns(proxy -> {
                    proxy.module().set(api.getModule());
                    proxy.path().set(api.getPath());
                    proxy.method().set(api.getMethod());
                    proxy.sortOrder().set(api.getSortOrder());
                    proxy.isEnabled().set(api.getIsEnabled());
                    proxy.remark().set(api.getRemark());
                    proxy.updatedAt().set(LocalDateTime.now());
                })
                .where(proxy -> {
                    proxy.id().eq(api.getId());
                    proxy.deletedAt().eq(0L);
                })
                .executeRows();
    }

    public void delete(Long id) {
        easyQuery.expressionUpdatable(SysResourceApiProxy.createTable())
                .setColumns(proxy -> proxy.deletedAt().set(System.currentTimeMillis()))
                .where(proxy -> {
                    proxy.id().eq(id);
                    proxy.deletedAt().eq(0L);
                })
                .executeRows();
    }
}
