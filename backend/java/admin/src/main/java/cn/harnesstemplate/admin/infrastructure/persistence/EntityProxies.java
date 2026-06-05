package cn.harnesstemplate.admin.infrastructure.persistence;

import cn.harnesstemplate.admin.domain.entity.*;
import com.easy.query.core.proxy.AbstractProxyEntity;
import com.easy.query.core.proxy.columns.types.SQLBooleanTypeColumn;
import com.easy.query.core.proxy.columns.types.SQLIntegerTypeColumn;
import com.easy.query.core.proxy.columns.types.SQLLocalDateTimeTypeColumn;
import com.easy.query.core.proxy.columns.types.SQLLongTypeColumn;
import com.easy.query.core.proxy.columns.types.SQLStringTypeColumn;

public final class EntityProxies {

    private EntityProxies() {
    }

    public abstract static class BaseEntityProxy<TProxy extends AbstractProxyEntity<TProxy, TEntity>, TEntity>
            extends AbstractProxyEntity<TProxy, TEntity> {
    }

    public static final class SysUserProxy extends BaseEntityProxy<SysUserProxy, SysUser> {
        public static SysUserProxy createTable() { return new SysUserProxy(); }
        @Override public Class<SysUser> getEntityClass() { return SysUser.class; }
        public SQLLongTypeColumn<SysUserProxy> id() { return getLongTypeColumn("id"); }
        public SQLLongTypeColumn<SysUserProxy> deletedAt() { return getLongTypeColumn("deletedAt"); }
        public SQLStringTypeColumn<SysUserProxy> username() { return getStringTypeColumn("username"); }
        public SQLStringTypeColumn<SysUserProxy> nickname() { return getStringTypeColumn("nickname"); }
        public SQLStringTypeColumn<SysUserProxy> languageCode() { return getStringTypeColumn("languageCode"); }
        public SQLBooleanTypeColumn<SysUserProxy> isEnabled() { return getBooleanTypeColumn("isEnabled"); }
        public SQLStringTypeColumn<SysUserProxy> remark() { return getStringTypeColumn("remark"); }
        public SQLLocalDateTimeTypeColumn<SysUserProxy> updatedAt() { return getLocalDateTimeTypeColumn("updatedAt"); }
    }

    public static final class SysRoleProxy extends BaseEntityProxy<SysRoleProxy, SysRole> {
        public static SysRoleProxy createTable() { return new SysRoleProxy(); }
        @Override public Class<SysRole> getEntityClass() { return SysRole.class; }
        public SQLLongTypeColumn<SysRoleProxy> id() { return getLongTypeColumn("id"); }
        public SQLLongTypeColumn<SysRoleProxy> deletedAt() { return getLongTypeColumn("deletedAt"); }
        public SQLStringTypeColumn<SysRoleProxy> name() { return getStringTypeColumn("name"); }
        public SQLStringTypeColumn<SysRoleProxy> code() { return getStringTypeColumn("code"); }
        public SQLLongTypeColumn<SysRoleProxy> parentId() { return getLongTypeColumn("parentId"); }
        public SQLBooleanTypeColumn<SysRoleProxy> isEnabled() { return getBooleanTypeColumn("isEnabled"); }
        public SQLStringTypeColumn<SysRoleProxy> remark() { return getStringTypeColumn("remark"); }
        public SQLLocalDateTimeTypeColumn<SysRoleProxy> updatedAt() { return getLocalDateTimeTypeColumn("updatedAt"); }
    }

    public static final class SysUserRoleProxy extends BaseEntityProxy<SysUserRoleProxy, SysUserRole> {
        public static SysUserRoleProxy createTable() { return new SysUserRoleProxy(); }
        @Override public Class<SysUserRole> getEntityClass() { return SysUserRole.class; }
        public SQLLongTypeColumn<SysUserRoleProxy> userId() { return getLongTypeColumn("userId"); }
        public SQLLongTypeColumn<SysUserRoleProxy> roleId() { return getLongTypeColumn("roleId"); }
        public SQLLongTypeColumn<SysUserRoleProxy> deletedAt() { return getLongTypeColumn("deletedAt"); }
    }

    public static final class SysRoleApiProxy extends BaseEntityProxy<SysRoleApiProxy, SysRoleApi> {
        public static SysRoleApiProxy createTable() { return new SysRoleApiProxy(); }
        @Override public Class<SysRoleApi> getEntityClass() { return SysRoleApi.class; }
        public SQLLongTypeColumn<SysRoleApiProxy> roleId() { return getLongTypeColumn("roleId"); }
        public SQLLongTypeColumn<SysRoleApiProxy> apiId() { return getLongTypeColumn("apiId"); }
        public SQLLongTypeColumn<SysRoleApiProxy> deletedAt() { return getLongTypeColumn("deletedAt"); }
    }

    public static final class SysRoleMenuProxy extends BaseEntityProxy<SysRoleMenuProxy, SysRoleMenu> {
        public static SysRoleMenuProxy createTable() { return new SysRoleMenuProxy(); }
        @Override public Class<SysRoleMenu> getEntityClass() { return SysRoleMenu.class; }
        public SQLLongTypeColumn<SysRoleMenuProxy> roleId() { return getLongTypeColumn("roleId"); }
        public SQLLongTypeColumn<SysRoleMenuProxy> menuId() { return getLongTypeColumn("menuId"); }
        public SQLLongTypeColumn<SysRoleMenuProxy> deletedAt() { return getLongTypeColumn("deletedAt"); }
    }

    public static final class SysResourceMenuApiProxy extends BaseEntityProxy<SysResourceMenuApiProxy, SysResourceMenuApi> {
        public static SysResourceMenuApiProxy createTable() { return new SysResourceMenuApiProxy(); }
        @Override public Class<SysResourceMenuApi> getEntityClass() { return SysResourceMenuApi.class; }
        public SQLLongTypeColumn<SysResourceMenuApiProxy> menuId() { return getLongTypeColumn("menuId"); }
        public SQLLongTypeColumn<SysResourceMenuApiProxy> apiId() { return getLongTypeColumn("apiId"); }
        public SQLLongTypeColumn<SysResourceMenuApiProxy> deletedAt() { return getLongTypeColumn("deletedAt"); }
    }

    public static final class SysResourceApiProxy extends BaseEntityProxy<SysResourceApiProxy, SysResourceApi> {
        public static SysResourceApiProxy createTable() { return new SysResourceApiProxy(); }
        @Override public Class<SysResourceApi> getEntityClass() { return SysResourceApi.class; }
        public SQLLongTypeColumn<SysResourceApiProxy> id() { return getLongTypeColumn("id"); }
        public SQLIntegerTypeColumn<SysResourceApiProxy> sortOrder() { return getIntegerTypeColumn("sortOrder"); }
        public SQLLongTypeColumn<SysResourceApiProxy> deletedAt() { return getLongTypeColumn("deletedAt"); }
        public SQLStringTypeColumn<SysResourceApiProxy> module() { return getStringTypeColumn("module"); }
        public SQLStringTypeColumn<SysResourceApiProxy> path() { return getStringTypeColumn("path"); }
        public SQLStringTypeColumn<SysResourceApiProxy> method() { return getStringTypeColumn("method"); }
        public SQLBooleanTypeColumn<SysResourceApiProxy> isEnabled() { return getBooleanTypeColumn("isEnabled"); }
        public SQLStringTypeColumn<SysResourceApiProxy> remark() { return getStringTypeColumn("remark"); }
        public SQLLocalDateTimeTypeColumn<SysResourceApiProxy> updatedAt() { return getLocalDateTimeTypeColumn("updatedAt"); }
    }

    public static final class SysResourceMenuProxy extends BaseEntityProxy<SysResourceMenuProxy, SysResourceMenu> {
        public static SysResourceMenuProxy createTable() { return new SysResourceMenuProxy(); }
        @Override public Class<SysResourceMenu> getEntityClass() { return SysResourceMenu.class; }
        public SQLLongTypeColumn<SysResourceMenuProxy> id() { return getLongTypeColumn("id"); }
        public SQLIntegerTypeColumn<SysResourceMenuProxy> sortOrder() { return getIntegerTypeColumn("sortOrder"); }
        public SQLLongTypeColumn<SysResourceMenuProxy> deletedAt() { return getLongTypeColumn("deletedAt"); }
        public SQLStringTypeColumn<SysResourceMenuProxy> metadata() { return getStringTypeColumn("metadata"); }
        public SQLBooleanTypeColumn<SysResourceMenuProxy> isEnabled() { return getBooleanTypeColumn("isEnabled"); }
        public SQLStringTypeColumn<SysResourceMenuProxy> menuType() { return getStringTypeColumn("menuType"); }
        public SQLStringTypeColumn<SysResourceMenuProxy> path() { return getStringTypeColumn("path"); }
        public SQLStringTypeColumn<SysResourceMenuProxy> redirect() { return getStringTypeColumn("redirect"); }
        public SQLStringTypeColumn<SysResourceMenuProxy> alias() { return getStringTypeColumn("alias"); }
        public SQLStringTypeColumn<SysResourceMenuProxy> name() { return getStringTypeColumn("name"); }
        public SQLStringTypeColumn<SysResourceMenuProxy> component() { return getStringTypeColumn("component"); }
        public SQLLongTypeColumn<SysResourceMenuProxy> parentId() { return getLongTypeColumn("parentId"); }
        public SQLStringTypeColumn<SysResourceMenuProxy> treePath() { return getStringTypeColumn("treePath"); }
        public SQLStringTypeColumn<SysResourceMenuProxy> remark() { return getStringTypeColumn("remark"); }
        public SQLLocalDateTimeTypeColumn<SysResourceMenuProxy> updatedAt() { return getLocalDateTimeTypeColumn("updatedAt"); }
    }

    public static final class SysLoginLogProxy extends BaseEntityProxy<SysLoginLogProxy, SysLoginLog> {
        public static SysLoginLogProxy createTable() { return new SysLoginLogProxy(); }
        @Override public Class<SysLoginLog> getEntityClass() { return SysLoginLog.class; }
        public SQLLongTypeColumn<SysLoginLogProxy> id() { return getLongTypeColumn("id"); }
        public SQLStringTypeColumn<SysLoginLogProxy> username() { return getStringTypeColumn("username"); }
    }

    public static final class SysApiLogProxy extends BaseEntityProxy<SysApiLogProxy, SysApiLog> {
        public static SysApiLogProxy createTable() { return new SysApiLogProxy(); }
        @Override public Class<SysApiLog> getEntityClass() { return SysApiLog.class; }
        public SQLLongTypeColumn<SysApiLogProxy> id() { return getLongTypeColumn("id"); }
        public SQLStringTypeColumn<SysApiLogProxy> method() { return getStringTypeColumn("method"); }
        public SQLStringTypeColumn<SysApiLogProxy> path() { return getStringTypeColumn("path"); }
    }

    public static final class SysDictTypeProxy extends BaseEntityProxy<SysDictTypeProxy, SysDictType> {
        public static SysDictTypeProxy createTable() { return new SysDictTypeProxy(); }
        @Override public Class<SysDictType> getEntityClass() { return SysDictType.class; }
        public SQLLongTypeColumn<SysDictTypeProxy> id() { return getLongTypeColumn("id"); }
        public SQLLongTypeColumn<SysDictTypeProxy> deletedAt() { return getLongTypeColumn("deletedAt"); }
        public SQLIntegerTypeColumn<SysDictTypeProxy> sortOrder() { return getIntegerTypeColumn("sortOrder"); }
        public SQLStringTypeColumn<SysDictTypeProxy> typeCode() { return getStringTypeColumn("typeCode"); }
        public SQLStringTypeColumn<SysDictTypeProxy> typeName() { return getStringTypeColumn("typeName"); }
        public SQLBooleanTypeColumn<SysDictTypeProxy> isEnabled() { return getBooleanTypeColumn("isEnabled"); }
        public SQLStringTypeColumn<SysDictTypeProxy> remark() { return getStringTypeColumn("remark"); }
        public SQLLocalDateTimeTypeColumn<SysDictTypeProxy> updatedAt() { return getLocalDateTimeTypeColumn("updatedAt"); }
    }

    public static final class SysDictEntryProxy extends BaseEntityProxy<SysDictEntryProxy, SysDictEntry> {
        public static SysDictEntryProxy createTable() { return new SysDictEntryProxy(); }
        @Override public Class<SysDictEntry> getEntityClass() { return SysDictEntry.class; }
        public SQLLongTypeColumn<SysDictEntryProxy> id() { return getLongTypeColumn("id"); }
        public SQLLongTypeColumn<SysDictEntryProxy> deletedAt() { return getLongTypeColumn("deletedAt"); }
        public SQLIntegerTypeColumn<SysDictEntryProxy> sortOrder() { return getIntegerTypeColumn("sortOrder"); }
        public SQLStringTypeColumn<SysDictEntryProxy> labelComponent() { return getStringTypeColumn("labelComponent"); }
        public SQLStringTypeColumn<SysDictEntryProxy> entryLabel() { return getStringTypeColumn("entryLabel"); }
        public SQLStringTypeColumn<SysDictEntryProxy> entryValue() { return getStringTypeColumn("entryValue"); }
        public SQLStringTypeColumn<SysDictEntryProxy> languageCode() { return getStringTypeColumn("languageCode"); }
        public SQLLongTypeColumn<SysDictEntryProxy> sysDictTypeId() { return getLongTypeColumn("sysDictTypeId"); }
        public SQLBooleanTypeColumn<SysDictEntryProxy> isEnabled() { return getBooleanTypeColumn("isEnabled"); }
        public SQLStringTypeColumn<SysDictEntryProxy> remark() { return getStringTypeColumn("remark"); }
        public SQLLocalDateTimeTypeColumn<SysDictEntryProxy> updatedAt() { return getLocalDateTimeTypeColumn("updatedAt"); }
    }

    public static final class SysLanguageTypeProxy extends BaseEntityProxy<SysLanguageTypeProxy, SysLanguageType> {
        public static SysLanguageTypeProxy createTable() { return new SysLanguageTypeProxy(); }
        @Override public Class<SysLanguageType> getEntityClass() { return SysLanguageType.class; }
        public SQLLongTypeColumn<SysLanguageTypeProxy> id() { return getLongTypeColumn("id"); }
        public SQLLongTypeColumn<SysLanguageTypeProxy> deletedAt() { return getLongTypeColumn("deletedAt"); }
        public SQLIntegerTypeColumn<SysLanguageTypeProxy> sortOrder() { return getIntegerTypeColumn("sortOrder"); }
        public SQLStringTypeColumn<SysLanguageTypeProxy> typeCode() { return getStringTypeColumn("typeCode"); }
        public SQLStringTypeColumn<SysLanguageTypeProxy> typeName() { return getStringTypeColumn("typeName"); }
        public SQLBooleanTypeColumn<SysLanguageTypeProxy> isDefault() { return getBooleanTypeColumn("isDefault"); }
        public SQLBooleanTypeColumn<SysLanguageTypeProxy> isEnabled() { return getBooleanTypeColumn("isEnabled"); }
        public SQLLocalDateTimeTypeColumn<SysLanguageTypeProxy> updatedAt() { return getLocalDateTimeTypeColumn("updatedAt"); }
    }

    public static final class SysLanguageEntryProxy extends BaseEntityProxy<SysLanguageEntryProxy, SysLanguageEntry> {
        public static SysLanguageEntryProxy createTable() { return new SysLanguageEntryProxy(); }
        @Override public Class<SysLanguageEntry> getEntityClass() { return SysLanguageEntry.class; }
        public SQLLongTypeColumn<SysLanguageEntryProxy> id() { return getLongTypeColumn("id"); }
        public SQLLongTypeColumn<SysLanguageEntryProxy> deletedAt() { return getLongTypeColumn("deletedAt"); }
        public SQLIntegerTypeColumn<SysLanguageEntryProxy> sortOrder() { return getIntegerTypeColumn("sortOrder"); }
        public SQLStringTypeColumn<SysLanguageEntryProxy> entryCode() { return getStringTypeColumn("entryCode"); }
        public SQLStringTypeColumn<SysLanguageEntryProxy> entryValue() { return getStringTypeColumn("entryValue"); }
        public SQLLongTypeColumn<SysLanguageEntryProxy> sysLanguageTypeId() { return getLongTypeColumn("sysLanguageTypeId"); }
        public SQLBooleanTypeColumn<SysLanguageEntryProxy> isEnabled() { return getBooleanTypeColumn("isEnabled"); }
        public SQLStringTypeColumn<SysLanguageEntryProxy> remark() { return getStringTypeColumn("remark"); }
        public SQLLocalDateTimeTypeColumn<SysLanguageEntryProxy> updatedAt() { return getLocalDateTimeTypeColumn("updatedAt"); }
    }

    public static final class FileAssetProxy extends BaseEntityProxy<FileAssetProxy, FileAsset> {
        public static FileAssetProxy createTable() { return new FileAssetProxy(); }
        @Override public Class<FileAsset> getEntityClass() { return FileAsset.class; }
        public SQLLongTypeColumn<FileAssetProxy> id() { return getLongTypeColumn("id"); }
        public SQLLongTypeColumn<FileAssetProxy> deletedAt() { return getLongTypeColumn("deletedAt"); }
        public SQLStringTypeColumn<FileAssetProxy> status() { return getStringTypeColumn("status"); }
        public SQLLocalDateTimeTypeColumn<FileAssetProxy> updatedAt() { return getLocalDateTimeTypeColumn("updatedAt"); }
    }

    public static final class KnowledgeCollectionProxy extends BaseEntityProxy<KnowledgeCollectionProxy, KnowledgeCollection> {
        public static KnowledgeCollectionProxy createTable() { return new KnowledgeCollectionProxy(); }
        @Override public Class<KnowledgeCollection> getEntityClass() { return KnowledgeCollection.class; }
        public SQLLongTypeColumn<KnowledgeCollectionProxy> id() { return getLongTypeColumn("id"); }
        public SQLLongTypeColumn<KnowledgeCollectionProxy> deletedAt() { return getLongTypeColumn("deletedAt"); }
        public SQLStringTypeColumn<KnowledgeCollectionProxy> collectionName() { return getStringTypeColumn("collectionName"); }
        public SQLStringTypeColumn<KnowledgeCollectionProxy> displayName() { return getStringTypeColumn("displayName"); }
        public SQLStringTypeColumn<KnowledgeCollectionProxy> metricType() { return getStringTypeColumn("metricType"); }
        public SQLStringTypeColumn<KnowledgeCollectionProxy> indexType() { return getStringTypeColumn("indexType"); }
        public SQLBooleanTypeColumn<KnowledgeCollectionProxy> isEnabled() { return getBooleanTypeColumn("isEnabled"); }
        public SQLStringTypeColumn<KnowledgeCollectionProxy> remark() { return getStringTypeColumn("remark"); }
        public SQLLocalDateTimeTypeColumn<KnowledgeCollectionProxy> updatedAt() { return getLocalDateTimeTypeColumn("updatedAt"); }
    }

    public static final class KnowledgeDocumentProxy extends BaseEntityProxy<KnowledgeDocumentProxy, KnowledgeDocument> {
        public static KnowledgeDocumentProxy createTable() { return new KnowledgeDocumentProxy(); }
        @Override public Class<KnowledgeDocument> getEntityClass() { return KnowledgeDocument.class; }
        public SQLLongTypeColumn<KnowledgeDocumentProxy> id() { return getLongTypeColumn("id"); }
        public SQLLongTypeColumn<KnowledgeDocumentProxy> deletedAt() { return getLongTypeColumn("deletedAt"); }
        public SQLLongTypeColumn<KnowledgeDocumentProxy> collectionId() { return getLongTypeColumn("collectionId"); }
        public SQLStringTypeColumn<KnowledgeDocumentProxy> title() { return getStringTypeColumn("title"); }
        public SQLStringTypeColumn<KnowledgeDocumentProxy> content() { return getStringTypeColumn("content"); }
        public SQLStringTypeColumn<KnowledgeDocumentProxy> contentType() { return getStringTypeColumn("contentType"); }
        public SQLStringTypeColumn<KnowledgeDocumentProxy> vectorStatus() { return getStringTypeColumn("vectorStatus"); }
        public SQLStringTypeColumn<KnowledgeDocumentProxy> vectorId() { return getStringTypeColumn("vectorId"); }
        public SQLStringTypeColumn<KnowledgeDocumentProxy> metadata() { return getStringTypeColumn("metadata"); }
        public SQLStringTypeColumn<KnowledgeDocumentProxy> indexingError() { return getStringTypeColumn("indexingError"); }
        public SQLLongTypeColumn<KnowledgeDocumentProxy> lastIndexedAt() { return getLongTypeColumn("lastIndexedAt"); }
        public SQLBooleanTypeColumn<KnowledgeDocumentProxy> isEnabled() { return getBooleanTypeColumn("isEnabled"); }
        public SQLStringTypeColumn<KnowledgeDocumentProxy> remark() { return getStringTypeColumn("remark"); }
        public SQLLocalDateTimeTypeColumn<KnowledgeDocumentProxy> updatedAt() { return getLocalDateTimeTypeColumn("updatedAt"); }
    }

    public static final class JobExecutionProxy extends BaseEntityProxy<JobExecutionProxy, JobExecution> {
        public static JobExecutionProxy createTable() { return new JobExecutionProxy(); }
        @Override public Class<JobExecution> getEntityClass() { return JobExecution.class; }
        public SQLLongTypeColumn<JobExecutionProxy> id() { return getLongTypeColumn("id"); }
        public SQLStringTypeColumn<JobExecutionProxy> jobCode() { return getStringTypeColumn("jobCode"); }
        public SQLStringTypeColumn<JobExecutionProxy> status() { return getStringTypeColumn("status"); }
        public SQLStringTypeColumn<JobExecutionProxy> resultJson() { return getStringTypeColumn("resultJson"); }
        public SQLStringTypeColumn<JobExecutionProxy> errorMessage() { return getStringTypeColumn("errorMessage"); }
        public SQLLocalDateTimeTypeColumn<JobExecutionProxy> endTime() { return getLocalDateTimeTypeColumn("endTime"); }
        public SQLLocalDateTimeTypeColumn<JobExecutionProxy> updatedAt() { return getLocalDateTimeTypeColumn("updatedAt"); }
    }

    public static final class JobScheduleProxy extends BaseEntityProxy<JobScheduleProxy, JobSchedule> {
        public static JobScheduleProxy createTable() { return new JobScheduleProxy(); }
        @Override public Class<JobSchedule> getEntityClass() { return JobSchedule.class; }
        public SQLLongTypeColumn<JobScheduleProxy> id() { return getLongTypeColumn("id"); }
        public SQLStringTypeColumn<JobScheduleProxy> jobCode() { return getStringTypeColumn("jobCode"); }
        public SQLStringTypeColumn<JobScheduleProxy> jobName() { return getStringTypeColumn("jobName"); }
        public SQLStringTypeColumn<JobScheduleProxy> workflowType() { return getStringTypeColumn("workflowType"); }
        public SQLStringTypeColumn<JobScheduleProxy> taskQueue() { return getStringTypeColumn("taskQueue"); }
        public SQLStringTypeColumn<JobScheduleProxy> scheduleType() { return getStringTypeColumn("scheduleType"); }
        public SQLStringTypeColumn<JobScheduleProxy> cronExpr() { return getStringTypeColumn("cronExpr"); }
        public SQLIntegerTypeColumn<JobScheduleProxy> intervalSeconds() { return getIntegerTypeColumn("intervalSeconds"); }
        public SQLLocalDateTimeTypeColumn<JobScheduleProxy> startTime() { return getLocalDateTimeTypeColumn("startTime"); }
        public SQLLocalDateTimeTypeColumn<JobScheduleProxy> endTime() { return getLocalDateTimeTypeColumn("endTime"); }
        public SQLStringTypeColumn<JobScheduleProxy> inputJson() { return getStringTypeColumn("inputJson"); }
        public SQLStringTypeColumn<JobScheduleProxy> status() { return getStringTypeColumn("status"); }
        public SQLStringTypeColumn<JobScheduleProxy> temporalScheduleId() { return getStringTypeColumn("temporalScheduleId"); }
        public SQLStringTypeColumn<JobScheduleProxy> temporalWorkflowIdPrefix() { return getStringTypeColumn("temporalWorkflowIdPrefix"); }
        public SQLStringTypeColumn<JobScheduleProxy> description() { return getStringTypeColumn("description"); }
        public SQLLocalDateTimeTypeColumn<JobScheduleProxy> updatedAt() { return getLocalDateTimeTypeColumn("updatedAt"); }
    }
}
