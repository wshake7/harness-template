package cn.harnesstemplate.admin.web.controller;

import cn.harnesstemplate.admin.web.dto.MenuResponse;
import cn.harnesstemplate.admin.web.dto.PageResult;
import cn.harnesstemplate.admin.web.dto.R;
import cn.harnesstemplate.admin.web.dto.SysRoleCreateRequest;
import cn.harnesstemplate.admin.web.dto.SysRoleResponse;
import cn.harnesstemplate.admin.web.dto.SysRoleUpdateRequest;
import cn.harnesstemplate.admin.web.dto.SysUserCreateRequest;
import cn.harnesstemplate.admin.web.dto.SysUserResponse;
import cn.harnesstemplate.admin.web.dto.SysUserUpdateRequest;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

class ControllerRecordContractTests {

    @Test
    void sharedControllerBoundaryDtosShouldUseRecords() {
        assertThat(R.class.isRecord()).isTrue();
        assertThat(PageResult.class.isRecord()).isTrue();
        assertThat(SysUserCreateRequest.class.isRecord()).isTrue();
        assertThat(SysUserUpdateRequest.class.isRecord()).isTrue();
        assertThat(SysRoleCreateRequest.class.isRecord()).isTrue();
        assertThat(SysRoleUpdateRequest.class.isRecord()).isTrue();
        assertThat(SysUserResponse.class.isRecord()).isTrue();
        assertThat(SysRoleResponse.class.isRecord()).isTrue();
        assertThat(MenuResponse.class.isRecord()).isTrue();
        assertThat(AccountController.LoginPwdRequest.class.isRecord()).isTrue();
    }
}
