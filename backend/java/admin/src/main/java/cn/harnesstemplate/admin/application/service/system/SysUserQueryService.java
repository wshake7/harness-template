package cn.harnesstemplate.admin.application.service.system;

import cn.harnesstemplate.admin.domain.entity.SysUser;
import cn.harnesstemplate.admin.infrastructure.persistence.EasyQuerySupport;
import cn.harnesstemplate.admin.infrastructure.persistence.EntityProxies.SysUserProxy;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;

@Service
public class SysUserQueryService {

    private final EasyQuerySupport easyQuery;

    public SysUserQueryService(EasyQuerySupport easyQuery) {
        this.easyQuery = easyQuery;
    }

    public List<SysUser> page(int pageNum, int pageSize, String username) {
        int offset = (pageNum - 1) * pageSize;
        return easyQuery.queryable(SysUserProxy.createTable())
                .where(user -> {
                    user.deletedAt().eq(0L);
                    user.username().like(username != null && !username.isEmpty(), username);
                })
                .orderBy(user -> user.id().asc())
                .limit(offset, pageSize)
                .toList();
    }

    public long count(String username) {
        return easyQuery.queryable(SysUserProxy.createTable())
                .where(user -> {
                    user.deletedAt().eq(0L);
                    user.username().like(username != null && !username.isEmpty(), username);
                })
                .count();
    }

    public SysUser findById(Long id) {
        return easyQuery.queryable(SysUserProxy.createTable())
                .where(user -> {
                    user.id().eq(id);
                    user.deletedAt().eq(0L);
                })
                .firstOrNull();
    }

    public SysUser findByUsername(String username) {
        return easyQuery.queryable(SysUserProxy.createTable())
                .where(user -> {
                    user.username().eq(username);
                    user.deletedAt().eq(0L);
                })
                .firstOrNull();
    }

    public SysUser save(SysUser user) {
        return easyQuery.insert(user);
    }

    public void update(SysUser user) {
        easyQuery.expressionUpdatable(SysUserProxy.createTable())
                .setColumns(proxy -> {
                    proxy.nickname().set(user.getNickname());
                    proxy.languageCode().set(user.getLanguageCode());
                    proxy.isEnabled().set(user.getIsEnabled());
                    proxy.remark().set(user.getRemark());
                    proxy.updatedAt().set(LocalDateTime.now());
                })
                .where(proxy -> {
                    proxy.id().eq(user.getId());
                    proxy.deletedAt().eq(0L);
                })
                .executeRows();
    }

    public void delete(Long id) {
        easyQuery.expressionUpdatable(SysUserProxy.createTable())
                .setColumns(proxy -> proxy.deletedAt().set(System.currentTimeMillis()))
                .where(proxy -> {
                    proxy.id().eq(id);
                    proxy.deletedAt().eq(0L);
                })
                .executeRows();
    }

    public List<SysUser> listAll() {
        return easyQuery.queryable(SysUserProxy.createTable())
                .where(user -> user.deletedAt().eq(0L))
                .toList();
    }
}
