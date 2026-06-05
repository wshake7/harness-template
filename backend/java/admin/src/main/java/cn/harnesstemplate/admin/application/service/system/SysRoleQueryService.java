package cn.harnesstemplate.admin.application.service.system;

import cn.harnesstemplate.admin.domain.entity.SysRole;
import cn.harnesstemplate.admin.infrastructure.persistence.EasyQuerySupport;
import cn.harnesstemplate.admin.infrastructure.persistence.EntityProxies.SysRoleProxy;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;

@Service
public class SysRoleQueryService {

    private final EasyQuerySupport easyQuery;

    public SysRoleQueryService(EasyQuerySupport easyQuery) {
        this.easyQuery = easyQuery;
    }

    public List<SysRole> page(int pageNum, int pageSize, String name) {
        int offset = (pageNum - 1) * pageSize;
        return easyQuery.queryable(SysRoleProxy.createTable())
                .where(role -> {
                    role.deletedAt().eq(0L);
                    role.name().like(name != null && !name.isEmpty(), name);
                })
                .orderBy(role -> role.id().asc())
                .limit(offset, pageSize)
                .toList();
    }

    public long count(String name) {
        return easyQuery.queryable(SysRoleProxy.createTable())
                .where(role -> {
                    role.deletedAt().eq(0L);
                    role.name().like(name != null && !name.isEmpty(), name);
                })
                .count();
    }

    public SysRole findById(Long id) {
        return easyQuery.queryable(SysRoleProxy.createTable())
                .where(role -> {
                    role.id().eq(id);
                    role.deletedAt().eq(0L);
                })
                .firstOrNull();
    }

    public SysRole findByCode(String code) {
        return easyQuery.queryable(SysRoleProxy.createTable())
                .where(role -> {
                    role.code().eq(code);
                    role.deletedAt().eq(0L);
                })
                .firstOrNull();
    }

    public SysRole save(SysRole role) {
        return easyQuery.insert(role);
    }

    public void update(SysRole role) {
        easyQuery.expressionUpdatable(SysRoleProxy.createTable())
                .setColumns(proxy -> {
                    proxy.name().set(role.getName());
                    proxy.code().set(role.getCode());
                    proxy.parentId().set(role.getParentId());
                    proxy.isEnabled().set(role.getIsEnabled());
                    proxy.remark().set(role.getRemark());
                    proxy.updatedAt().set(LocalDateTime.now());
                })
                .where(proxy -> {
                    proxy.id().eq(role.getId());
                    proxy.deletedAt().eq(0L);
                })
                .executeRows();
    }

    public void delete(Long id) {
        easyQuery.expressionUpdatable(SysRoleProxy.createTable())
                .setColumns(proxy -> proxy.deletedAt().set(System.currentTimeMillis()))
                .where(proxy -> {
                    proxy.id().eq(id);
                    proxy.deletedAt().eq(0L);
                })
                .executeRows();
    }

    public List<SysRole> listAll() {
        return easyQuery.queryable(SysRoleProxy.createTable())
                .where(role -> role.deletedAt().eq(0L))
                .orderBy(role -> role.id().asc())
                .toList();
    }

    public boolean hasChildren(Long parentId) {
        return easyQuery.queryable(SysRoleProxy.createTable())
                .where(role -> {
                    role.parentId().eq(parentId);
                    role.deletedAt().eq(0L);
                })
                .count() > 0;
    }
}
