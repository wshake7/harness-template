package cn.harnesstemplate.admin.domain.entity;

import com.easy.query.core.annotation.Table;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.experimental.FieldNameConstants;

@Data
@EqualsAndHashCode(callSuper = true)
@FieldNameConstants
@Table("file_asset")
public class FileAsset extends BaseAuditEntity {

    private String remark;

    private String engine;

    private String bucket;

    private String objectKey;

    private String originalName;

    private String contentType;

    private String extension;

    private Long size;

    private String sha256;

    private String bizType;

    private String bizId;

    private String metadata;

    private String status;
}
