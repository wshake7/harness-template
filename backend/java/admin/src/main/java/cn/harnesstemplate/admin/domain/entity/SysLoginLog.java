package cn.harnesstemplate.admin.domain.entity;

import com.easy.query.core.annotation.Table;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.experimental.FieldNameConstants;

import java.time.LocalDateTime;

@Data
@EqualsAndHashCode(callSuper = true)
@FieldNameConstants
@Table("sys_login_log")
public class SysLoginLog extends BaseEntity {

    private String username;

    private String loginIp;

    private String loginMac;

    private LocalDateTime loginTime;

    private String userAgent;

    private String browserName;

    private String browserVersion;

    private String clientId;

    private String clientName;

    private String osName;

    private String osVersion;

    private Long sysUserId;

    private Integer statusCode;

    private Boolean success;

    private String reason;

    private String location;
}
