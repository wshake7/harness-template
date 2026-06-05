package cn.harnesstemplate.admin.application.port;

import cn.harnesstemplate.admin.domain.model.SysRole;
import cn.harnesstemplate.admin.domain.model.SysUser;

import java.util.Collections;
import java.util.List;
import java.util.Optional;

/**
 * Stub implementation used when no real repository is configured.
 * Replaced by Easy Query implementation in Milestone 3.
 */
public class StubUserLookupRepository implements UserLookupRepository {

    @Override
    public Optional<SysUser> findByUsername(String username) {
        return Optional.empty();
    }

    @Override
    public List<SysRole> findRolesByUserId(Long userId) {
        return Collections.emptyList();
    }
}
