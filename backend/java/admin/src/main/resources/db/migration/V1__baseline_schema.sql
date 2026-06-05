-- V1__baseline_schema.sql
-- Flyway baseline migration — mirrors Go ORM models at backend/go/admin/internal/services/orm/models/
-- Generated: 2026-05-26
-- Column naming matches Go GORM tags (create_by not create_by)
-- Purpose: Create all tables, indexes, and constraints that Java Admin service must support.
-- Notes:
--   - deleted_at uses BIGINT millisecond timestamp (0 = not deleted), matching GORM soft_delete:milli
--   - Conditional unique indexes use partial WHERE clauses matching Go's `where:deleted_at = 0`
--   - Foreign keys use ON UPDATE CASCADE ON DELETE RESTRICT
--   - JSON fields use JSONB (PostgreSQL) for efficient queries

-- =============================================================================
-- System tables: User, Role, Permissions
-- =============================================================================

CREATE TABLE sys_user (
    id              BIGSERIAL PRIMARY KEY,
    created_at      TIMESTAMPTZ,
    updated_at      TIMESTAMPTZ,
    create_by      BIGINT NOT NULL DEFAULT 0,
    update_by      BIGINT NOT NULL DEFAULT 0,
    delete_by      BIGINT NOT NULL DEFAULT 0,
    remark          TEXT DEFAULT '',
    is_enabled      BOOLEAN DEFAULT TRUE,
    deleted_at      BIGINT NOT NULL DEFAULT 0,
    username        VARCHAR(64) NOT NULL,
    nickname        VARCHAR(64) NOT NULL DEFAULT '',
    password        VARCHAR(255) NOT NULL DEFAULT '',
    language_code   VARCHAR(32) NOT NULL DEFAULT ''
);
CREATE INDEX idx_sys_user_deleted_at ON sys_user (deleted_at);
CREATE UNIQUE INDEX idx_sys_user_username_active ON sys_user (username) WHERE deleted_at = 0;

CREATE TABLE sys_role (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    create_by  BIGINT NOT NULL DEFAULT 0,
    update_by  BIGINT NOT NULL DEFAULT 0,
    delete_by  BIGINT NOT NULL DEFAULT 0,
    remark      TEXT DEFAULT '',
    is_enabled  BOOLEAN DEFAULT TRUE,
    deleted_at  BIGINT NOT NULL DEFAULT 0,
    name        VARCHAR(255) NOT NULL,
    code        VARCHAR(128) NOT NULL,
    parent_id   BIGINT
);
CREATE INDEX idx_sys_role_deleted_at ON sys_role (deleted_at);
CREATE UNIQUE INDEX idx_sys_role_code_active ON sys_role (code) WHERE deleted_at = 0;

CREATE TABLE sys_user_role (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    create_by  BIGINT NOT NULL DEFAULT 0,
    update_by  BIGINT NOT NULL DEFAULT 0,
    delete_by  BIGINT NOT NULL DEFAULT 0,
    user_id     BIGINT NOT NULL,
    role_id     BIGINT NOT NULL,
    deleted_at  BIGINT NOT NULL DEFAULT 0
);
CREATE INDEX idx_sys_user_role_user_id ON sys_user_role (user_id);
CREATE INDEX idx_sys_user_role_role_id ON sys_user_role (role_id);
CREATE INDEX idx_sys_user_role_deleted_at ON sys_user_role (deleted_at);
CREATE UNIQUE INDEX idx_sys_user_role_user_role_active ON sys_user_role (user_id, role_id) WHERE deleted_at = 0;

-- =============================================================================
-- Resource: Menu, API, and associations
-- =============================================================================

CREATE TABLE sys_resource_menu (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    create_by  BIGINT NOT NULL DEFAULT 0,
    update_by  BIGINT NOT NULL DEFAULT 0,
    delete_by  BIGINT NOT NULL DEFAULT 0,
    remark      TEXT DEFAULT '',
    sort_order  INT DEFAULT 0,
    metadata    JSONB DEFAULT '{}',
    is_enabled  BOOLEAN DEFAULT TRUE,
    deleted_at  BIGINT NOT NULL DEFAULT 0,
    menu_type   VARCHAR(32) NOT NULL DEFAULT 'MENU',
    path        VARCHAR(1024) NOT NULL,
    redirect    VARCHAR(1024) NOT NULL DEFAULT '',
    alias       VARCHAR(255) NOT NULL DEFAULT '',
    name        VARCHAR(255) NOT NULL DEFAULT '',
    component   VARCHAR(255) NOT NULL DEFAULT '',
    parent_id   BIGINT,
    tree_path   VARCHAR(1024)
);
CREATE INDEX idx_sys_resource_menu_deleted_at ON sys_resource_menu (deleted_at);

CREATE TABLE sys_resource_api (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    create_by  BIGINT NOT NULL DEFAULT 0,
    update_by  BIGINT NOT NULL DEFAULT 0,
    delete_by  BIGINT NOT NULL DEFAULT 0,
    remark      TEXT DEFAULT '',
    sort_order  INT DEFAULT 0,
    is_enabled  BOOLEAN DEFAULT TRUE,
    deleted_at  BIGINT NOT NULL DEFAULT 0,
    module      VARCHAR(128) NOT NULL DEFAULT '',
    path        VARCHAR(255) NOT NULL,
    method      VARCHAR(16) NOT NULL
);
CREATE INDEX idx_sys_resource_api_deleted_at ON sys_resource_api (deleted_at);
CREATE UNIQUE INDEX idx_sys_resource_api_method_path_active ON sys_resource_api (method, path) WHERE deleted_at = 0;

CREATE TABLE sys_resource_menu_api (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    create_by  BIGINT NOT NULL DEFAULT 0,
    update_by  BIGINT NOT NULL DEFAULT 0,
    delete_by  BIGINT NOT NULL DEFAULT 0,
    menu_id     BIGINT NOT NULL,
    api_id      BIGINT NOT NULL,
    deleted_at  BIGINT NOT NULL DEFAULT 0
);
CREATE INDEX idx_sys_resource_menu_api_menu_id ON sys_resource_menu_api (menu_id);
CREATE INDEX idx_sys_resource_menu_api_api_id ON sys_resource_menu_api (api_id);
CREATE INDEX idx_sys_resource_menu_api_deleted_at ON sys_resource_menu_api (deleted_at);
CREATE UNIQUE INDEX idx_sys_resource_menu_api_menu_api_active ON sys_resource_menu_api (menu_id, api_id) WHERE deleted_at = 0;

-- =============================================================================
-- Role-Menu / Role-API authorization
-- =============================================================================

CREATE TABLE sys_role_menu (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    create_by  BIGINT NOT NULL DEFAULT 0,
    update_by  BIGINT NOT NULL DEFAULT 0,
    delete_by  BIGINT NOT NULL DEFAULT 0,
    role_id     BIGINT NOT NULL,
    menu_id     BIGINT NOT NULL,
    deleted_at  BIGINT NOT NULL DEFAULT 0
);
CREATE INDEX idx_sys_role_menu_role_id ON sys_role_menu (role_id);
CREATE INDEX idx_sys_role_menu_menu_id ON sys_role_menu (menu_id);
CREATE INDEX idx_sys_role_menu_deleted_at ON sys_role_menu (deleted_at);
CREATE UNIQUE INDEX idx_sys_role_menu_role_menu_active ON sys_role_menu (role_id, menu_id) WHERE deleted_at = 0;

CREATE TABLE sys_role_api (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    create_by  BIGINT NOT NULL DEFAULT 0,
    update_by  BIGINT NOT NULL DEFAULT 0,
    delete_by  BIGINT NOT NULL DEFAULT 0,
    role_id     BIGINT NOT NULL,
    api_id      BIGINT NOT NULL,
    deleted_at  BIGINT NOT NULL DEFAULT 0
);
CREATE INDEX idx_sys_role_api_role_id ON sys_role_api (role_id);
CREATE INDEX idx_sys_role_api_api_id ON sys_role_api (api_id);
CREATE INDEX idx_sys_role_api_deleted_at ON sys_role_api (deleted_at);
CREATE UNIQUE INDEX idx_sys_role_api_role_api_active ON sys_role_api (role_id, api_id) WHERE deleted_at = 0;

-- =============================================================================
-- Data permissions
-- =============================================================================

CREATE TABLE sys_data_permission (
    id              BIGSERIAL PRIMARY KEY,
    created_at      TIMESTAMPTZ,
    updated_at      TIMESTAMPTZ,
    create_by      BIGINT NOT NULL DEFAULT 0,
    update_by      BIGINT NOT NULL DEFAULT 0,
    delete_by      BIGINT NOT NULL DEFAULT 0,
    remark          TEXT DEFAULT '',
    is_enabled      BOOLEAN DEFAULT TRUE,
    deleted_at      BIGINT NOT NULL DEFAULT 0,
    subject_type    VARCHAR(16) NOT NULL,
    subject_id      BIGINT NOT NULL,
    resource_table  VARCHAR(32) NOT NULL,
    action          JSONB NOT NULL DEFAULT '["read"]',
    action_key      VARCHAR(64) NOT NULL DEFAULT 'read',
    scope_type      VARCHAR(32) NOT NULL DEFAULT 'none',
    scope_field     VARCHAR(64) NOT NULL DEFAULT 'id',
    scope_values    JSONB NOT NULL DEFAULT '[]',
    conditions      JSONB NOT NULL DEFAULT '{}',
    priority        INT NOT NULL DEFAULT 0
);
CREATE INDEX idx_sys_data_permission_subject ON sys_data_permission (subject_type, subject_id);
CREATE INDEX idx_sys_data_permission_deleted_at ON sys_data_permission (deleted_at);
CREATE UNIQUE INDEX idx_sys_data_permission_subject_resource_action_active ON sys_data_permission (subject_type, subject_id, resource_table, action_key) WHERE deleted_at = 0;

-- =============================================================================
-- Dictionary
-- =============================================================================

CREATE TABLE sys_dict_type (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    create_by  BIGINT NOT NULL DEFAULT 0,
    update_by  BIGINT NOT NULL DEFAULT 0,
    delete_by  BIGINT NOT NULL DEFAULT 0,
    remark      TEXT DEFAULT '',
    is_enabled  BOOLEAN DEFAULT TRUE,
    sort_order  INT DEFAULT 0,
    deleted_at  BIGINT NOT NULL DEFAULT 0,
    type_code   VARCHAR(128) NOT NULL,
    type_name   VARCHAR(255) NOT NULL
);
CREATE INDEX idx_sys_dict_type_deleted_at ON sys_dict_type (deleted_at);
CREATE UNIQUE INDEX idx_sys_dict_type_type_code_active ON sys_dict_type (type_code) WHERE deleted_at = 0;

CREATE TABLE sys_dict_entry (
    id               BIGSERIAL PRIMARY KEY,
    created_at       TIMESTAMPTZ,
    updated_at       TIMESTAMPTZ,
    create_by       BIGINT NOT NULL DEFAULT 0,
    update_by       BIGINT NOT NULL DEFAULT 0,
    delete_by       BIGINT NOT NULL DEFAULT 0,
    remark           TEXT DEFAULT '',
    is_enabled       BOOLEAN DEFAULT TRUE,
    sort_order       INT DEFAULT 0,
    deleted_at       BIGINT NOT NULL DEFAULT 0,
    label_component  VARCHAR(255) NOT NULL DEFAULT '',
    entry_label      VARCHAR(255) NOT NULL,
    entry_value      VARCHAR(255) NOT NULL,
    language_code    VARCHAR(32) NOT NULL DEFAULT '',
    sys_dict_type_id BIGINT NOT NULL
);
CREATE INDEX idx_sys_dict_entry_dict_type_id ON sys_dict_entry (sys_dict_type_id);
CREATE INDEX idx_sys_dict_entry_deleted_at ON sys_dict_entry (deleted_at);
CREATE UNIQUE INDEX idx_sys_dict_entry_type_lang_value_active ON sys_dict_entry (sys_dict_type_id, language_code, entry_value) WHERE deleted_at = 0;

-- =============================================================================
-- Language / i18n
-- =============================================================================

CREATE TABLE sys_language_type (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    create_by  BIGINT NOT NULL DEFAULT 0,
    update_by  BIGINT NOT NULL DEFAULT 0,
    delete_by  BIGINT NOT NULL DEFAULT 0,
    sort_order  INT DEFAULT 0,
    is_enabled  BOOLEAN DEFAULT TRUE,
    deleted_at  BIGINT NOT NULL DEFAULT 0,
    type_code   VARCHAR(128) NOT NULL,
    type_name   VARCHAR(255) NOT NULL,
    is_default  BOOLEAN DEFAULT FALSE
);
CREATE INDEX idx_sys_language_type_deleted_at ON sys_language_type (deleted_at);
CREATE UNIQUE INDEX idx_sys_language_type_code_active ON sys_language_type (type_code) WHERE deleted_at = 0;
CREATE UNIQUE INDEX idx_sys_language_type_default_active ON sys_language_type (is_default) WHERE is_default = TRUE AND deleted_at = 0;

CREATE TABLE sys_language_entry (
    id                   BIGSERIAL PRIMARY KEY,
    created_at           TIMESTAMPTZ,
    updated_at           TIMESTAMPTZ,
    create_by           BIGINT NOT NULL DEFAULT 0,
    update_by           BIGINT NOT NULL DEFAULT 0,
    delete_by           BIGINT NOT NULL DEFAULT 0,
    remark               TEXT DEFAULT '',
    is_enabled           BOOLEAN DEFAULT TRUE,
    sort_order           INT DEFAULT 0,
    deleted_at           BIGINT NOT NULL DEFAULT 0,
    entry_code           VARCHAR(128) NOT NULL,
    entry_value          VARCHAR(255) NOT NULL,
    sys_language_type_id BIGINT NOT NULL
);
CREATE INDEX idx_sys_language_entry_deleted_at ON sys_language_entry (deleted_at);
CREATE UNIQUE INDEX idx_sys_language_entry_code_type_active ON sys_language_entry (entry_code, sys_language_type_id) WHERE deleted_at = 0;

-- =============================================================================
-- Logs
-- =============================================================================

CREATE TABLE sys_api_log (
    id               BIGSERIAL PRIMARY KEY,
    created_at       TIMESTAMPTZ,
    request_id       VARCHAR(128) NOT NULL,
    method           VARCHAR(16) NOT NULL,
    module           VARCHAR(255) DEFAULT '',
    path             VARCHAR(255) NOT NULL,
    referer          TEXT DEFAULT '',
    before_change    TEXT DEFAULT '',
    after_change     TEXT DEFAULT '',
    format_change    TEXT DEFAULT '',
    request_uri      TEXT DEFAULT '',
    request_body     TEXT DEFAULT '',
    request_header   TEXT DEFAULT '',
    response         TEXT DEFAULT '',
    cost_time        BIGINT DEFAULT 0,
    sys_user_id      BIGINT,
    client_ip        VARCHAR(64) DEFAULT '',
    status_code      INT DEFAULT 0,
    reason           VARCHAR(255) DEFAULT '',
    success          BOOLEAN DEFAULT FALSE,
    location         VARCHAR(255) DEFAULT '',
    user_agent       TEXT DEFAULT '',
    browser_name     VARCHAR(128) DEFAULT '',
    browser_version  VARCHAR(128) DEFAULT '',
    client_id        VARCHAR(128) DEFAULT '',
    client_name      VARCHAR(128) DEFAULT '',
    os_name          VARCHAR(128) DEFAULT '',
    os_version       VARCHAR(128) DEFAULT ''
);
CREATE UNIQUE INDEX idx_sys_api_log_request_id ON sys_api_log (request_id);
CREATE INDEX idx_sys_api_log_method_path ON sys_api_log (method, path);
CREATE INDEX idx_sys_api_log_sys_user_id ON sys_api_log (sys_user_id);
CREATE INDEX idx_sys_api_log_status_code ON sys_api_log (status_code);
CREATE INDEX idx_sys_api_log_created_at ON sys_api_log (created_at);
CREATE INDEX idx_sys_api_log_user_created_at ON sys_api_log (sys_user_id, created_at);
CREATE INDEX idx_sys_api_log_success_created_at ON sys_api_log (success, created_at);
CREATE INDEX idx_sys_api_log_status_created_at ON sys_api_log (status_code, created_at);

CREATE TABLE sys_login_log (
    id               BIGSERIAL PRIMARY KEY,
    created_at       TIMESTAMPTZ,
    username         VARCHAR(64) NOT NULL DEFAULT '',
    login_ip         VARCHAR(64) NOT NULL DEFAULT '',
    login_mac        VARCHAR(128) NOT NULL DEFAULT '',
    login_time       TIMESTAMPTZ,
    user_agent       TEXT DEFAULT '',
    browser_name     VARCHAR(128) DEFAULT '',
    browser_version  VARCHAR(128) DEFAULT '',
    client_id        VARCHAR(128) DEFAULT '',
    client_name      VARCHAR(128) DEFAULT '',
    os_name          VARCHAR(128) DEFAULT '',
    os_version       VARCHAR(128) DEFAULT '',
    sys_user_id      BIGINT,
    status_code      INT DEFAULT 0,
    success          BOOLEAN DEFAULT FALSE,
    reason           VARCHAR(255) DEFAULT '',
    location         VARCHAR(255) DEFAULT ''
);
CREATE INDEX idx_sys_login_log_username ON sys_login_log (username);
CREATE INDEX idx_sys_login_log_login_ip ON sys_login_log (login_ip);
CREATE INDEX idx_sys_login_log_sys_user_id ON sys_login_log (sys_user_id);
CREATE INDEX idx_sys_login_log_status_code ON sys_login_log (status_code);
CREATE INDEX idx_sys_login_log_login_time ON sys_login_log (login_time);
CREATE INDEX idx_sys_login_log_success_login_time ON sys_login_log (success, login_time);
CREATE INDEX idx_sys_login_log_status_login_time ON sys_login_log (status_code, login_time);
CREATE INDEX idx_sys_login_log_user_login_time ON sys_login_log (sys_user_id, login_time);

-- =============================================================================
-- Casbin
-- =============================================================================

CREATE TABLE sys_casbin_model (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    create_by  BIGINT NOT NULL DEFAULT 0,
    update_by  BIGINT NOT NULL DEFAULT 0,
    delete_by  BIGINT NOT NULL DEFAULT 0,
    remark      TEXT DEFAULT '',
    is_enabled  BOOLEAN DEFAULT TRUE,
    deleted_at  BIGINT NOT NULL DEFAULT 0,
    name        VARCHAR(255) NOT NULL,
    content     TEXT NOT NULL
);
CREATE INDEX idx_sys_casbin_model_deleted_at ON sys_casbin_model (deleted_at);
CREATE UNIQUE INDEX idx_sys_casbin_model_name_active ON sys_casbin_model (name) WHERE deleted_at = 0;

-- casbin_rule table for jCasbin JDBC adapter
CREATE TABLE casbin_rule (
    id    BIGSERIAL PRIMARY KEY,
    ptype VARCHAR(255) NOT NULL DEFAULT '',
    v0    VARCHAR(255) NOT NULL DEFAULT '',
    v1    VARCHAR(255) NOT NULL DEFAULT '',
    v2    VARCHAR(255) NOT NULL DEFAULT '',
    v3    VARCHAR(255) NOT NULL DEFAULT '',
    v4    VARCHAR(255) NOT NULL DEFAULT '',
    v5    VARCHAR(255) NOT NULL DEFAULT ''
);

-- =============================================================================
-- Job scheduling
-- =============================================================================

CREATE TABLE job_schedule (
    id                          BIGSERIAL PRIMARY KEY,
    created_at                  TIMESTAMPTZ,
    updated_at                  TIMESTAMPTZ,
    job_code                    VARCHAR(128) NOT NULL,
    job_name                    VARCHAR(255) NOT NULL,
    workflow_type               VARCHAR(255) NOT NULL,
    task_queue                  VARCHAR(255) NOT NULL,
    schedule_type               VARCHAR(32) NOT NULL,
    cron_expr                   VARCHAR(128) DEFAULT '',
    interval_seconds            INT,
    start_time                  TIMESTAMPTZ,
    end_time                    TIMESTAMPTZ,
    input_json                  JSONB,
    status                      VARCHAR(32) NOT NULL DEFAULT 'ENABLED',
    temporal_schedule_id        VARCHAR(255) DEFAULT '',
    temporal_workflow_id_prefix VARCHAR(255) DEFAULT '',
    description                 VARCHAR(512) DEFAULT ''
);
CREATE UNIQUE INDEX idx_job_schedule_job_code ON job_schedule (job_code);
CREATE INDEX idx_job_schedule_status ON job_schedule (status);
CREATE INDEX idx_job_schedule_schedule_type ON job_schedule (schedule_type);

CREATE TABLE job_execution (
    id                   BIGSERIAL PRIMARY KEY,
    created_at           TIMESTAMPTZ,
    updated_at           TIMESTAMPTZ,
    job_code             VARCHAR(128) NOT NULL,
    temporal_workflow_id VARCHAR(255) NOT NULL,
    temporal_run_id      VARCHAR(255) DEFAULT '',
    trigger_time         TIMESTAMPTZ NOT NULL,
    start_time           TIMESTAMPTZ,
    end_time             TIMESTAMPTZ,
    status               VARCHAR(32) NOT NULL DEFAULT 'RUNNING',
    input_json           JSONB,
    result_json          JSONB,
    error_message        TEXT DEFAULT '',
    retry_count          INT NOT NULL DEFAULT 0
);
CREATE INDEX idx_job_execution_job_code ON job_execution (job_code);
CREATE UNIQUE INDEX uk_job_execution_workflow_run ON job_execution (temporal_workflow_id, temporal_run_id);
CREATE INDEX idx_job_execution_status ON job_execution (status);
CREATE INDEX idx_job_execution_trigger_time ON job_execution (trigger_time);

-- =============================================================================
-- File assets / Object storage
-- =============================================================================

CREATE TABLE file_asset (
    id            BIGSERIAL PRIMARY KEY,
    created_at    TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ,
    create_by    BIGINT NOT NULL DEFAULT 0,
    update_by    BIGINT NOT NULL DEFAULT 0,
    delete_by    BIGINT NOT NULL DEFAULT 0,
    remark        TEXT DEFAULT '',
    deleted_at    BIGINT NOT NULL DEFAULT 0,
    engine        VARCHAR(32) NOT NULL,
    bucket        VARCHAR(128) NOT NULL,
    object_key    VARCHAR(512) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    content_type  VARCHAR(128) NOT NULL DEFAULT '',
    extension     VARCHAR(32) NOT NULL DEFAULT '',
    size          BIGINT NOT NULL DEFAULT 0,
    sha256        CHAR(64) NOT NULL DEFAULT '',
    biz_type      VARCHAR(64) NOT NULL DEFAULT '',
    biz_id        VARCHAR(128) NOT NULL DEFAULT '',
    metadata      JSONB DEFAULT '{}',
    status        VARCHAR(32) NOT NULL DEFAULT 'active'
);
CREATE INDEX idx_file_asset_deleted_at ON file_asset (deleted_at);
CREATE INDEX idx_file_asset_sha256 ON file_asset (sha256);
CREATE INDEX idx_file_asset_biz ON file_asset (biz_type, biz_id);
CREATE UNIQUE INDEX uk_file_asset_object_key_active ON file_asset (object_key) WHERE deleted_at = 0;

-- =============================================================================
-- Knowledge base
-- =============================================================================

CREATE TABLE knowledge_collection (
    id              BIGSERIAL PRIMARY KEY,
    created_at      TIMESTAMPTZ,
    updated_at      TIMESTAMPTZ,
    create_by      BIGINT NOT NULL DEFAULT 0,
    update_by      BIGINT NOT NULL DEFAULT 0,
    delete_by      BIGINT NOT NULL DEFAULT 0,
    remark          TEXT DEFAULT '',
    is_enabled      BOOLEAN DEFAULT TRUE,
    deleted_at      BIGINT NOT NULL DEFAULT 0,
    collection_name VARCHAR(128) NOT NULL,
    display_name    VARCHAR(255) NOT NULL,
    metric_type     VARCHAR(32) NOT NULL DEFAULT 'COSINE',
    index_type      VARCHAR(32) NOT NULL DEFAULT 'auto'
);
CREATE INDEX idx_knowledge_collection_deleted_at ON knowledge_collection (deleted_at);
CREATE UNIQUE INDEX idx_knowledge_collection_name_active ON knowledge_collection (collection_name) WHERE deleted_at = 0;

CREATE TABLE knowledge_document (
    id               BIGSERIAL PRIMARY KEY,
    created_at       TIMESTAMPTZ,
    updated_at       TIMESTAMPTZ,
    create_by       BIGINT NOT NULL DEFAULT 0,
    update_by       BIGINT NOT NULL DEFAULT 0,
    delete_by       BIGINT NOT NULL DEFAULT 0,
    remark           TEXT DEFAULT '',
    is_enabled       BOOLEAN DEFAULT TRUE,
    deleted_at       BIGINT NOT NULL DEFAULT 0,
    collection_id    BIGINT NOT NULL,
    document_id      VARCHAR(255) NOT NULL,
    title            VARCHAR(512) NOT NULL,
    content          TEXT NOT NULL,
    content_type     VARCHAR(64) NOT NULL DEFAULT 'text',
    source           VARCHAR(512) DEFAULT '',
    chunk_index      INT NOT NULL DEFAULT 0,
    total_chunks     INT NOT NULL DEFAULT 1,
    vector_status    VARCHAR(32) NOT NULL DEFAULT 'pending',
    vector_id        VARCHAR(255) DEFAULT '',
    metadata         JSONB DEFAULT '{}',
    indexing_error   TEXT DEFAULT '',
    last_indexed_at  BIGINT DEFAULT 0
);
CREATE INDEX idx_knowledge_document_collection_id ON knowledge_document (collection_id);
CREATE INDEX idx_knowledge_document_deleted_at ON knowledge_document (deleted_at);
CREATE UNIQUE INDEX idx_knowledge_document_doc_id_active ON knowledge_document (document_id) WHERE deleted_at = 0;

-- =============================================================================
-- AI Agent
-- =============================================================================

CREATE TABLE agent_session (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    create_by  BIGINT NOT NULL DEFAULT 0,
    update_by  BIGINT NOT NULL DEFAULT 0,
    delete_by  BIGINT NOT NULL DEFAULT 0,
    remark      TEXT DEFAULT '',
    deleted_at  BIGINT NOT NULL DEFAULT 0,
    session_id  VARCHAR(64) NOT NULL,
    title       VARCHAR(255) DEFAULT '',
    status      VARCHAR(32) NOT NULL DEFAULT 'active',
    metadata    JSONB DEFAULT '{}'
);
CREATE INDEX idx_agent_session_deleted_at ON agent_session (deleted_at);
CREATE UNIQUE INDEX idx_agent_session_session_id_active ON agent_session (session_id) WHERE deleted_at = 0;

CREATE TABLE agent_message (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    create_by  BIGINT NOT NULL DEFAULT 0,
    update_by  BIGINT NOT NULL DEFAULT 0,
    delete_by  BIGINT NOT NULL DEFAULT 0,
    deleted_at  BIGINT NOT NULL DEFAULT 0,
    session_id  BIGINT NOT NULL,
    role        VARCHAR(32) NOT NULL,
    content     TEXT NOT NULL,
    token_count INT NOT NULL DEFAULT 0,
    metadata    JSONB DEFAULT '{}'
);
CREATE INDEX idx_agent_message_session_id ON agent_message (session_id);
CREATE INDEX idx_agent_message_deleted_at ON agent_message (deleted_at);

-- =============================================================================
-- Foreign key constraints
-- =============================================================================

ALTER TABLE sys_user_role ADD CONSTRAINT fk_sys_user_role_user
    FOREIGN KEY (user_id) REFERENCES sys_user (id) ON UPDATE CASCADE ON DELETE RESTRICT;
ALTER TABLE sys_user_role ADD CONSTRAINT fk_sys_user_role_role
    FOREIGN KEY (role_id) REFERENCES sys_role (id) ON UPDATE CASCADE ON DELETE RESTRICT;

ALTER TABLE sys_resource_menu_api ADD CONSTRAINT fk_menu_api_menu
    FOREIGN KEY (menu_id) REFERENCES sys_resource_menu (id) ON UPDATE CASCADE ON DELETE RESTRICT;
ALTER TABLE sys_resource_menu_api ADD CONSTRAINT fk_menu_api_api
    FOREIGN KEY (api_id) REFERENCES sys_resource_api (id) ON UPDATE CASCADE ON DELETE RESTRICT;

ALTER TABLE sys_role_menu ADD CONSTRAINT fk_role_menu_role
    FOREIGN KEY (role_id) REFERENCES sys_role (id) ON UPDATE CASCADE ON DELETE RESTRICT;
ALTER TABLE sys_role_menu ADD CONSTRAINT fk_role_menu_menu
    FOREIGN KEY (menu_id) REFERENCES sys_resource_menu (id) ON UPDATE CASCADE ON DELETE RESTRICT;

ALTER TABLE sys_role_api ADD CONSTRAINT fk_role_api_role
    FOREIGN KEY (role_id) REFERENCES sys_role (id) ON UPDATE CASCADE ON DELETE RESTRICT;
ALTER TABLE sys_role_api ADD CONSTRAINT fk_role_api_api
    FOREIGN KEY (api_id) REFERENCES sys_resource_api (id) ON UPDATE CASCADE ON DELETE RESTRICT;

ALTER TABLE sys_dict_entry ADD CONSTRAINT fk_dict_entry_type
    FOREIGN KEY (sys_dict_type_id) REFERENCES sys_dict_type (id) ON UPDATE CASCADE ON DELETE RESTRICT;

ALTER TABLE sys_language_entry ADD CONSTRAINT fk_language_entry_type
    FOREIGN KEY (sys_language_type_id) REFERENCES sys_language_type (id) ON UPDATE CASCADE ON DELETE RESTRICT;

ALTER TABLE sys_api_log ADD CONSTRAINT fk_api_log_user
    FOREIGN KEY (sys_user_id) REFERENCES sys_user (id) ON UPDATE CASCADE ON DELETE SET NULL;

ALTER TABLE sys_login_log ADD CONSTRAINT fk_login_log_user
    FOREIGN KEY (sys_user_id) REFERENCES sys_user (id) ON UPDATE CASCADE ON DELETE SET NULL;

ALTER TABLE knowledge_document ADD CONSTRAINT fk_knowledge_document_collection
    FOREIGN KEY (collection_id) REFERENCES knowledge_collection (id) ON UPDATE CASCADE ON DELETE RESTRICT;

ALTER TABLE agent_message ADD CONSTRAINT fk_agent_message_session
    FOREIGN KEY (session_id) REFERENCES agent_session (id) ON UPDATE CASCADE ON DELETE RESTRICT;
