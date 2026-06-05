package cn.harnesstemplate.admin.application.port;

import cn.harnesstemplate.admin.domain.model.SysRole;
import cn.harnesstemplate.admin.domain.model.SysUser;
import cn.harnesstemplate.admin.infrastructure.persistence.EasyQuerySupport;
import cn.harnesstemplate.admin.infrastructure.persistence.EntityProxies.SysRoleProxy;
import cn.harnesstemplate.admin.infrastructure.persistence.EntityProxies.SysUserProxy;
import cn.harnesstemplate.admin.infrastructure.persistence.EntityProxies.SysUserRoleProxy;
import org.springframework.stereotype.Component;

import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;

@Component
public class EasyQueryUserLookupRepository implements UserLookupRepository {

    private final EasyQuerySupport easyQuery;

    public EasyQueryUserLookupRepository(EasyQuerySupport easyQuery) {
        this.easyQuery = easyQuery;
    }

    @Override
    public Optional<SysUser> findByUsername(String username) {
        var user = easyQuery.queryable(SysUserProxy.createTable())
                .where(proxy -> {
                    proxy.username().eq(username);
                    proxy.deletedAt().eq(0L);
                })
                .firstOrNull();

        if (user == null) {
            return Optional.empty();
        }

        return Optional.of(new SysUser(user.getId(), user.getUsername(), user.getPassword()));
    }

    @Override
    public List<SysRole> findRolesByUserId(Long userId) {
        var userRoles = easyQuery.queryable(SysUserRoleProxy.createTable())
                .where(proxy -> {
                    proxy.userId().eq(userId);
                    proxy.deletedAt().eq(0L);
                })
                .toList();

        if (userRoles.isEmpty()) {
            return List.of();
        }

        var roleIds = userRoles.stream()
                .map(cn.harnesstemplate.admin.domain.entity.SysUserRole::getRoleId)
                .toList();

        return easyQuery.queryable(SysRoleProxy.createTable())
                .where(proxy -> {
                    proxy.id().in(roleIds);
                    proxy.deletedAt().eq(0L);
                    proxy.isEnabled().eq(true);
                })
                .toList()
                .stream()
                .map(role -> new SysRole(role.getId(), role.getCode()))
                .collect(Collectors.toList());
    }
}
