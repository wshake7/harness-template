package cn.harnesstemplate.admin.domain.entity;

import com.easy.query.core.annotation.Table;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.experimental.FieldNameConstants;

import java.time.LocalDateTime;

@Data
@EqualsAndHashCode(callSuper = true)
@FieldNameConstants
@Table("sys_api_log")
public class SysApiLog extends BaseEntity {

    private String requestId;

    private String method;

    private String module;

    private String path;

    private String referer;

    private String beforeChange;

    private String afterChange;

    private String formatChange;

    private String requestUri;

    private String requestBody;

    private String requestHeader;

    private String response;

    private Long costTime;

    private Long sysUserId;

    private String clientIp;

    private Integer statusCode;

    private String reason;

    private Boolean success;

    private String location;

    private String userAgent;

    private String browserName;

    private String browserVersion;

    private String clientId;

    private String clientName;

    private String osName;

    private String osVersion;
}
