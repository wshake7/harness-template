package cn.harnesstemplate.admin.application.port;

import cn.harnesstemplate.admin.domain.model.SysRole;
import cn.harnesstemplate.admin.domain.model.SysUser;

import java.util.List;
import java.util.Optional;

/**
 * Port for user/role lookup. Implemented by persistence layer in Milestone 3.
 */
public interface UserLookupRepository {

    Optional<SysUser> findByUsername(String username);

    List<SysRole> findRolesByUserId(Long userId);
}
