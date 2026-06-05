package cn.harnesstemplate.admin.application.service.system;

import cn.harnesstemplate.admin.domain.entity.SysResourceMenu;
import cn.harnesstemplate.admin.infrastructure.persistence.EasyQuerySupport;
import cn.harnesstemplate.admin.infrastructure.persistence.EntityProxies.SysResourceMenuProxy;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;

@Service
public class SysResourceMenuQueryService {

    private final EasyQuerySupport easyQuery;

    public SysResourceMenuQueryService(EasyQuerySupport easyQuery) {
        this.easyQuery = easyQuery;
    }

    public List<SysResourceMenu> listAll() {
        return easyQuery.queryable(SysResourceMenuProxy.createTable())
                .where(menu -> menu.deletedAt().eq(0L))
                .orderBy(menu -> {
                    menu.sortOrder().asc();
                    menu.id().asc();
                })
                .toList();
    }

    public SysResourceMenu findById(Long id) {
        return easyQuery.queryable(SysResourceMenuProxy.createTable())
                .where(menu -> {
                    menu.id().eq(id);
                    menu.deletedAt().eq(0L);
                })
                .firstOrNull();
    }

    public SysResourceMenu save(SysResourceMenu menu) {
        return easyQuery.insert(menu);
    }

    public void update(SysResourceMenu menu) {
        easyQuery.expressionUpdatable(SysResourceMenuProxy.createTable())
                .setColumns(proxy -> {
                    proxy.parentId().set(menu.getParentId());
                    proxy.menuType().set(menu.getMenuType());
                    proxy.path().set(menu.getPath());
                    proxy.redirect().set(menu.getRedirect());
                    proxy.alias().set(menu.getAlias());
                    proxy.name().set(menu.getName());
                    proxy.component().set(menu.getComponent());
                    proxy.metadata().set(menu.getMetadata());
                    proxy.sortOrder().set(menu.getSortOrder());
                    proxy.isEnabled().set(menu.getIsEnabled());
                    proxy.remark().set(menu.getRemark());
                    proxy.treePath().set(menu.getTreePath());
                    proxy.updatedAt().set(LocalDateTime.now());
                })
                .where(proxy -> {
                    proxy.id().eq(menu.getId());
                    proxy.deletedAt().eq(0L);
                })
                .executeRows();
    }

    public void delete(Long id) {
        easyQuery.expressionUpdatable(SysResourceMenuProxy.createTable())
                .setColumns(proxy -> proxy.deletedAt().set(System.currentTimeMillis()))
                .where(proxy -> {
                    proxy.id().eq(id);
                    proxy.deletedAt().eq(0L);
                })
                .executeRows();
    }

    public boolean hasChildren(Long parentId) {
        return easyQuery.queryable(SysResourceMenuProxy.createTable())
                .where(menu -> {
                    menu.parentId().eq(parentId);
                    menu.deletedAt().eq(0L);
                })
                .count() > 0;
    }
}
