INSERT INTO
    "public"."sys_casbin_model" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "is_enabled",
        "remark",
        "deleted_at",
        "name",
        "content"
    )
VALUES
    (
        1,
        NULL,
        NULL,
        0,
        0,
        0,
        TRUE,
        '',
        0,
        'pbac',
        '[request_definition]
                r = sub, obj, act

                [policy_definition]
                p = sub, obj, act

                [policy_effect]
                e = some(where (p.eft == allow))

                [matchers]
                m = r.sub == p.sub && keyMatch2(r.obj, p.obj) && r.act == p.act'
    );

INSERT INTO
    "public"."sys_data_permission" (
        "id",
        "created_at",
        "updated_at",
        "remark",
        "created_by",
        "updated_by",
        "deleted_by",
        "is_enabled",
        "deleted_at",
        "subject_type",
        "subject_id",
        "resource_table",
        "action",
        "action_key",
        "scope_type",
        "scope_field",
        "scope_values",
        "conditions",
        "priority"
    )
VALUES
    (
        1,
        NULL,
        NULL,
        'root can access all dict types',
        0,
        0,
        0,
        TRUE,
        0,
        'ROLE',
        1,
        'sys_dict_type',
        '["all"]',
        'all',
        'all',
        'id',
        '[]',
        '{}',
        100
    ),
    (
        2,
        NULL,
        NULL,
        'all roles can read all dict types',
        0,
        0,
        0,
        TRUE,
        0,
        'ANY_ROLE',
        0,
        'sys_dict_type',
        '["read"]',
        'read',
        'all',
        'id',
        '[]',
        '{}',
        0
    ),
    (
        3,
        NULL,
        NULL,
        'all roles can operate dict types except system:is_enabled',
        0,
        0,
        0,
        TRUE,
        0,
        'ANY_ROLE',
        0,
        'sys_dict_type',
        '["write","delete"]',
        'write,delete',
        'custom',
        'id',
        '[]',
        '{"id__not": 1}',
        0
    );

INSERT INTO
    "public"."sys_dict_type" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "is_enabled",
        "sort_order",
        "remark",
        "deleted_at",
        "type_code",
        "type_name"
    )
VALUES
    (
        1,
        NULL,
        NULL,
        0,
        0,
        0,
        TRUE,
        0,
        '',
        0,
        'system:is_enabled',
        '开关状态'
    ),
    (
        2,
        NULL,
        NULL,
        0,
        0,
        0,
        TRUE,
        0,
        '',
        0,
        'system:resource_menu_type',
        '菜单类型'
    ),
    (
        3,
        NULL,
        NULL,
        0,
        0,
        0,
        TRUE,
        0,
        '',
        0,
        'system:resource_api_type',
        'API 请求方法'
    );

INSERT INTO
    "public"."sys_dict_entry" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "sort_order",
        "is_enabled",
        "remark",
        "deleted_at",
        "label_component",
        "entry_label",
        "entry_value",
        "language_code",
        "sys_dict_type_id"
    )
VALUES
    (
        1,
        NULL,
        NULL,
        0,
        0,
        0,
        0,
        TRUE,
        '',
        0,
        '<Tag color="success">${EntryLabel}</Tag>',
        '启用',
        '1',
        '',
        1
    ),
    (
        2,
        NULL,
        NULL,
        0,
        0,
        0,
        0,
        TRUE,
        '',
        0,
        '<Tag color="error">${EntryLabel}</Tag>',
        '停用',
        '0',
        '',
        1
    ),
    (
        3,
        NULL,
        NULL,
        0,
        0,
        0,
        0,
        TRUE,
        '',
        0,
        '<Tag color="green">${EntryLabel}</Tag>',
        '目录',
        'CATALOG',
        '',
        2
    ),
    (
        4,
        NULL,
        NULL,
        0,
        0,
        0,
        1,
        TRUE,
        '',
        0,
        '<Tag color="blue">${EntryLabel}</Tag>',
        '菜单',
        'MENU',
        '',
        2
    ),
    (
        5,
        NULL,
        NULL,
        0,
        0,
        0,
        2,
        TRUE,
        '',
        0,
        '<Tag color="purple">${EntryLabel}</Tag>',
        '按钮',
        'BUTTON',
        '',
        2
    ),
    (
        6,
        NULL,
        NULL,
        0,
        0,
        0,
        3,
        TRUE,
        '',
        0,
        '<Tag color="cyan">${EntryLabel}</Tag>',
        '内嵌',
        'EMBEDDED',
        '',
        2
    ),
    (
        7,
        NULL,
        NULL,
        0,
        0,
        0,
        4,
        TRUE,
        '',
        0,
        '<Tag color="orange">${EntryLabel}</Tag>',
        '外链',
        'LINK',
        '',
        2
    ),
    (
        8,
        NULL,
        NULL,
        0,
        0,
        0,
        0,
        TRUE,
        '',
        0,
        '<Tag color="green">${EntryLabel}</Tag>',
        'GET',
        'GET',
        '',
        3
    ),
    (
        9,
        NULL,
        NULL,
        0,
        0,
        0,
        1,
        TRUE,
        '',
        0,
        '<Tag color="blue">${EntryLabel}</Tag>',
        'POST',
        'POST',
        '',
        3
    ),
    (
        10,
        NULL,
        NULL,
        0,
        0,
        0,
        2,
        TRUE,
        '',
        0,
        '<Tag color="orange">${EntryLabel}</Tag>',
        'PUT',
        'PUT',
        '',
        3
    ),
    (
        11,
        NULL,
        NULL,
        0,
        0,
        0,
        3,
        TRUE,
        '',
        0,
        '<Tag color="gold">${EntryLabel}</Tag>',
        'PATCH',
        'PATCH',
        '',
        3
    ),
    (
        12,
        NULL,
        NULL,
        0,
        0,
        0,
        4,
        TRUE,
        '',
        0,
        '<Tag color="red">${EntryLabel}</Tag>',
        'DELETE',
        'DELETE',
        '',
        3
    ),
    (
        13,
        NULL,
        NULL,
        0,
        0,
        0,
        5,
        TRUE,
        '',
        0,
        '<Tag color="purple">${EntryLabel}</Tag>',
        'OPTIONS',
        'OPTIONS',
        '',
        3
    ),
    (
        14,
        NULL,
        NULL,
        0,
        0,
        0,
        6,
        TRUE,
        '',
        0,
        '<Tag color="cyan">${EntryLabel}</Tag>',
        'HEAD',
        'HEAD',
        '',
        3
    );

INSERT INTO
    "public"."sys_role" (
        "id",
        "created_at",
        "updated_at",
        "remark",
        "created_by",
        "updated_by",
        "deleted_by",
        "is_enabled",
        "deleted_at",
        "name",
        "code",
        "parent_id"
    )
VALUES
    (
        1,
        NULL,
        NULL,
        '',
        0,
        0,
        0,
        TRUE,
        0,
        '超级管理员',
        'root',
        NULL
    );

INSERT INTO
    "public"."sys_user" (
        "id",
        "created_at",
        "updated_at",
        "remark",
        "created_by",
        "updated_by",
        "deleted_by",
        "is_enabled",
        "deleted_at",
        "username",
        "nickname",
        "password",
        "language_code"
    )
VALUES
    (
        1,
        '2026-05-06 14:13:43.158062+08',
        '2026-05-06 14:13:43.158062+08',
        '',
        0,
        0,
        0,
        TRUE,
        0,
        'root',
        '',
        '$2a$04$ASoVUxXahEpdD9.dxfwsHuUw3PqQ/yAZ0gD2KnqtMAqSGZ4VZCSVO',
        ''
    ),
    (
        2,
        '2026-05-06 14:13:43.159636+08',
        '2026-05-06 14:13:43.159636+08',
        '',
        0,
        0,
        0,
        TRUE,
        0,
        'admin',
        '',
        '$2a$04$ASoVUxXahEpdD9.dxfwsHuUw3PqQ/yAZ0gD2KnqtMAqSGZ4VZCSVO',
        ''
    );

INSERT INTO
    "public"."sys_user_role" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "user_id",
        "role_id",
        "deleted_at"
    )
VALUES
    (1, NULL, NULL, 0, 0, 0, 1, 1, 0);

INSERT INTO
    "public"."sys_resource_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "remark",
        "sort_order",
        "metadata",
        "is_enabled",
        "deleted_at",
        "menu_type",
        "path",
        "redirect",
        "alias",
        "name",
        "component",
        "parent_id",
        "tree_path"
    )
VALUES
    (
        1,
        '2026-05-09 15:44:29.247468+08',
        '2026-05-09 17:31:30.430037+08',
        1,
        1,
        0,
        '',
        0,
        '{"icon": "", "hidden": false}',
        't',
        0,
        'CATALOG',
        '/dashboard',
        '',
        '',
        '控制台1',
        '',
        NULL,
        '/1/'
    );

INSERT INTO
    "public"."sys_resource_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "remark",
        "sort_order",
        "metadata",
        "is_enabled",
        "deleted_at",
        "menu_type",
        "path",
        "redirect",
        "alias",
        "name",
        "component",
        "parent_id",
        "tree_path"
    )
VALUES
    (
        2,
        '2026-05-09 16:01:45.414731+08',
        '2026-05-09 17:54:40.486485+08',
        1,
        1,
        0,
        '',
        1,
        '{"icon": "SettingOutlined", "order": 1, "hidden": false}',
        't',
        0,
        'CATALOG',
        '/system',
        '',
        '',
        '系统管理',
        '',
        NULL,
        '/2/'
    );

INSERT INTO
    "public"."sys_resource_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "remark",
        "sort_order",
        "metadata",
        "is_enabled",
        "deleted_at",
        "menu_type",
        "path",
        "redirect",
        "alias",
        "name",
        "component",
        "parent_id",
        "tree_path"
    )
VALUES
    (
        3,
        '2026-05-09 16:13:02.748554+08',
        '2026-05-10 20:06:34.713589+08',
        1,
        1,
        0,
        '',
        -2,
        '{"icon": "", "order": -2, "hidden": false}',
        't',
        0,
        'MENU',
        '/system/resource/menu',
        '',
        '',
        '菜单管理',
        '/system/resource.menu.tsx',
        2,
        '/2/3/'
    );

INSERT INTO
    "public"."sys_resource_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "remark",
        "sort_order",
        "metadata",
        "is_enabled",
        "deleted_at",
        "menu_type",
        "path",
        "redirect",
        "alias",
        "name",
        "component",
        "parent_id",
        "tree_path"
    )
VALUES
    (
        4,
        '2026-05-10 20:10:00.656317+08',
        '2026-05-10 20:10:45.771044+08',
        1,
        1,
        0,
        '',
        3,
        '{"icon": "ReconciliationOutlined", "order": 3, "hidden": false}',
        't',
        0,
        'CATALOG',
        '/logger',
        '',
        '',
        '日志信息',
        '',
        NULL,
        '/10/'
    );

INSERT INTO
    "public"."sys_resource_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "remark",
        "sort_order",
        "metadata",
        "is_enabled",
        "deleted_at",
        "menu_type",
        "path",
        "redirect",
        "alias",
        "name",
        "component",
        "parent_id",
        "tree_path"
    )
VALUES
    (
        5,
        '2026-05-09 16:42:05.166618+08',
        '2026-05-10 19:54:52.382881+08',
        1,
        1,
        0,
        '',
        0,
        '{"icon": "", "hidden": false}',
        't',
        0,
        'MENU',
        '/system/language',
        '',
        '',
        '语言管理',
        '/system/language.tsx',
        2,
        '/2/5/'
    );

INSERT INTO
    "public"."sys_resource_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "remark",
        "sort_order",
        "metadata",
        "is_enabled",
        "deleted_at",
        "menu_type",
        "path",
        "redirect",
        "alias",
        "name",
        "component",
        "parent_id",
        "tree_path"
    )
VALUES
    (
        6,
        '2026-05-09 16:42:36.997454+08',
        '2026-05-10 19:54:57.037816+08',
        1,
        1,
        0,
        '',
        0,
        '{"icon": "", "hidden": false}',
        't',
        0,
        'MENU',
        '/system/dict',
        '',
        '',
        '字典管理',
        '/system/dict.tsx',
        2,
        '/2/6/'
    );

INSERT INTO
    "public"."sys_resource_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "remark",
        "sort_order",
        "metadata",
        "is_enabled",
        "deleted_at",
        "menu_type",
        "path",
        "redirect",
        "alias",
        "name",
        "component",
        "parent_id",
        "tree_path"
    )
VALUES
    (
        7,
        '2026-05-09 16:49:51.215257+08',
        '2026-05-10 20:28:47.004123+08',
        1,
        1,
        0,
        '',
        2,
        '{"icon": "UserOutlined", "order": 2, "hidden": false}',
        't',
        0,
        'CATALOG',
        '/account',
        '',
        '',
        '账号管理',
        '',
        NULL,
        '/7/'
    );

INSERT INTO
    "public"."sys_resource_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "remark",
        "sort_order",
        "metadata",
        "is_enabled",
        "deleted_at",
        "menu_type",
        "path",
        "redirect",
        "alias",
        "name",
        "component",
        "parent_id",
        "tree_path"
    )
VALUES
    (
        8,
        '2026-05-09 16:50:35.892181+08',
        '2026-05-09 17:55:29.888606+08',
        1,
        1,
        0,
        '',
        0,
        '{"icon": "", "hidden": false}',
        't',
        0,
        'MENU',
        '/account/role',
        '',
        '',
        '角色管理',
        '/account/role.tsx',
        7,
        '/7/8/'
    );

INSERT INTO
    "public"."sys_resource_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "remark",
        "sort_order",
        "metadata",
        "is_enabled",
        "deleted_at",
        "menu_type",
        "path",
        "redirect",
        "alias",
        "name",
        "component",
        "parent_id",
        "tree_path"
    )
VALUES
    (
        9,
        '2026-05-10 17:38:27.019193+08',
        '2026-05-10 20:06:43.984546+08',
        1,
        1,
        0,
        '',
        -1,
        '{"icon": "", "order": -1, "hidden": false}',
        't',
        0,
        'MENU',
        '/system/resource/api',
        '',
        '',
        'API管理',
        '/system/resource.api.tsx',
        2,
        '/2/9/'
    );

INSERT INTO
    "public"."sys_resource_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "remark",
        "sort_order",
        "metadata",
        "is_enabled",
        "deleted_at",
        "menu_type",
        "path",
        "redirect",
        "alias",
        "name",
        "component",
        "parent_id",
        "tree_path"
    )
VALUES
    (
        10,
        '2026-05-09 16:17:20.233627+08',
        '2026-05-10 20:10:15.368034+08',
        1,
        1,
        0,
        '',
        0,
        '{"icon": "", "order": 0, "hidden": false}',
        't',
        0,
        'MENU',
        '/logger/api/log',
        '',
        '',
        'API日志',
        '/logger/api.log.tsx',
        4,
        '/10/4/'
    );

INSERT INTO
    "public"."sys_resource_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "remark",
        "sort_order",
        "metadata",
        "is_enabled",
        "deleted_at",
        "menu_type",
        "path",
        "redirect",
        "alias",
        "name",
        "component",
        "parent_id",
        "tree_path"
    )
VALUES
    (
        11,
        '2026-05-10 20:28:36.585473+08',
        '2026-05-10 20:28:36.587573+08',
        1,
        1,
        0,
        '',
        0,
        '{"icon": "", "order": 0, "hidden": false}',
        't',
        0,
        'MENU',
        '/account/user',
        '',
        '',
        '用户管理',
        '/account/user.tsx',
        7,
        '/7/11/'
    );

INSERT INTO
    "public"."sys_resource_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "remark",
        "sort_order",
        "metadata",
        "is_enabled",
        "deleted_at",
        "menu_type",
        "path",
        "redirect",
        "alias",
        "name",
        "component",
        "parent_id",
        "tree_path"
    )
VALUES
    (
        12,
        '2026-05-13 00:00:00+08',
        '2026-05-13 00:00:00+08',
        1,
        1,
        0,
        '',
        1,
        '{"icon": "", "order": 1, "hidden": false}',
        't',
        0,
        'MENU',
        '/logger/login/log',
        '',
        '',
        '登录日志',
        '/logger/login.log.tsx',
        4,
        '/10/12/'
    );

INSERT INTO
    "public"."sys_resource_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "remark",
        "sort_order",
        "metadata",
        "is_enabled",
        "deleted_at",
        "menu_type",
        "path",
        "redirect",
        "alias",
        "name",
        "component",
        "parent_id",
        "tree_path"
    )
VALUES
    (
        13,
        '2026-05-14 00:00:00+08',
        '2026-05-14 00:00:00+08',
        1,
        1,
        0,
        '',
        4,
        '{"icon": "ScheduleOutlined", "order": 4, "hidden": false}',
        't',
        0,
        'CATALOG',
        '/job',
        '',
        '',
        '任务调度',
        '',
        NULL,
        '/13/'
    );

INSERT INTO
    "public"."sys_resource_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "remark",
        "sort_order",
        "metadata",
        "is_enabled",
        "deleted_at",
        "menu_type",
        "path",
        "redirect",
        "alias",
        "name",
        "component",
        "parent_id",
        "tree_path"
    )
VALUES
    (
        14,
        '2026-05-14 00:00:00+08',
        '2026-05-14 00:00:00+08',
        1,
        1,
        0,
        '',
        0,
        '{"icon": "", "order": 0, "hidden": false}',
        't',
        0,
        'MENU',
        '/job/schedule',
        '',
        '',
        '任务配置',
        '/job/schedule.tsx',
        13,
        '/13/14/'
    );

INSERT INTO
    "public"."sys_resource_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "remark",
        "sort_order",
        "metadata",
        "is_enabled",
        "deleted_at",
        "menu_type",
        "path",
        "redirect",
        "alias",
        "name",
        "component",
        "parent_id",
        "tree_path"
    )
VALUES
    (
        15,
        '2026-05-14 00:00:00+08',
        '2026-05-14 00:00:00+08',
        1,
        1,
        0,
        '',
        1,
        '{"icon": "", "order": 1, "hidden": false}',
        't',
        0,
        'MENU',
        '/job/execution',
        '',
        '',
        '执行记录',
        '/job/execution.tsx',
        13,
        '/13/15/'
    );

INSERT INTO
    "public"."sys_resource_api" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "remark",
        "sort_order",
        "is_enabled",
        "deleted_at",
        "module",
        "path",
        "method"
    )
VALUES
    (
        1,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'role',
        '/api/sys/role/list',
        'POST'
    ),
    (
        2,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'role',
        '/api/sys/role/tree',
        'GET'
    ),
    (
        3,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'role',
        '/api/sys/role/:id/permissions',
        'GET'
    ),
    (
        4,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'role',
        '/api/sys/role/create',
        'POST'
    ),
    (
        5,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'role',
        '/api/sys/role/update',
        'POST'
    ),
    (
        6,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'role',
        '/api/sys/role/del',
        'POST'
    ),
    (
        7,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'role',
        '/api/sys/role/permissions',
        'POST'
    ),
    (
        8,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'user',
        '/api/sys/user/list',
        'POST'
    ),
    (
        9,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'user',
        '/api/sys/user/create',
        'POST'
    ),
    (
        10,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'user',
        '/api/sys/user/update',
        'POST'
    ),
    (
        11,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'user',
        '/api/sys/user/del',
        'POST'
    ),
    (
        12,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'dict',
        '/api/sys/dict/type/list',
        'POST'
    ),
    (
        13,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'dict',
        '/api/sys/dict/type/create',
        'POST'
    ),
    (
        14,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'dict',
        '/api/sys/dict/type/update',
        'POST'
    ),
    (
        15,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'dict',
        '/api/sys/dict/type/del',
        'POST'
    ),
    (
        16,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'dict',
        '/api/sys/dict/entry/list',
        'POST'
    ),
    (
        17,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'dict',
        '/api/sys/dict/entry/match',
        'POST'
    ),
    (
        18,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'dict',
        '/api/sys/dict/entry/create',
        'POST'
    ),
    (
        19,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'dict',
        '/api/sys/dict/entry/update',
        'POST'
    ),
    (
        20,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'dict',
        '/api/sys/dict/entry/del',
        'POST'
    ),
    (
        21,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'dict',
        '/api/sys/dict/entry/batch/copy',
        'POST'
    ),
    (
        22,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'language',
        '/api/sys/language/type/list',
        'POST'
    ),
    (
        23,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'language',
        '/api/sys/language/type/create',
        'POST'
    ),
    (
        24,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'language',
        '/api/sys/language/type/update',
        'POST'
    ),
    (
        25,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'language',
        '/api/sys/language/type/del',
        'POST'
    ),
    (
        26,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'language',
        '/api/sys/language/entry/list',
        'POST'
    ),
    (
        27,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'language',
        '/api/sys/language/entry/create',
        'POST'
    ),
    (
        28,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'language',
        '/api/sys/language/entry/update',
        'POST'
    ),
    (
        29,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'language',
        '/api/sys/language/entry/del',
        'POST'
    ),
    (
        30,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'language',
        '/api/sys/language/entry/batch/create',
        'POST'
    ),
    (
        31,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'api_log',
        '/api/sys/api/log/list',
        'POST'
    ),
    (
        32,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'api_log',
        '/api/sys/api/log/detail',
        'POST'
    ),
    (
        33,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'resource_menu',
        '/api/sys/resource/menu/list',
        'POST'
    ),
    (
        34,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'resource_menu',
        '/api/sys/resource/menu/tree',
        'GET'
    ),
    (
        35,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'resource_menu',
        '/api/sys/resource/menu/create',
        'POST'
    ),
    (
        36,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'resource_menu',
        '/api/sys/resource/menu/update',
        'POST'
    ),
    (
        37,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'resource_menu',
        '/api/sys/resource/menu/del',
        'POST'
    ),
    (
        38,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'resource_api',
        '/api/sys/resource/api/list',
        'POST'
    ),
    (
        39,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'resource_api',
        '/api/sys/resource/api/create',
        'POST'
    ),
    (
        40,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'resource_api',
        '/api/sys/resource/api/update',
        'POST'
    ),
    (
        41,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'resource_api',
        '/api/sys/resource/api/del',
        'POST'
    ),
    (
        42,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'login_log',
        '/api/sys/login/log/list',
        'POST'
    ),
    (
        43,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'login_log',
        '/api/sys/login/log/detail',
        'POST'
    ),
    (
        44,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'events',
        '/api/events',
        'GET'
    ),
    (
        45,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'job_schedule',
        '/api/sys/job/schedule/list',
        'POST'
    ),
    (
        46,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'job_schedule',
        '/api/sys/job/schedule/detail',
        'POST'
    ),
    (
        47,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'job_schedule',
        '/api/sys/job/schedule/create',
        'POST'
    ),
    (
        48,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'job_schedule',
        '/api/sys/job/schedule/update',
        'POST'
    ),
    (
        49,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'job_schedule',
        '/api/sys/job/schedule/del',
        'POST'
    ),
    (
        50,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'job_schedule',
        '/api/sys/job/schedule/switch',
        'POST'
    ),
    (
        51,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'job_schedule',
        '/api/sys/job/schedule/sync',
        'POST'
    ),
    (
        52,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'job_schedule',
        '/api/sys/job/schedule/trigger',
        'POST'
    ),
    (
        53,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'job_execution',
        '/api/sys/job/execution/list',
        'POST'
    ),
    (
        54,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'job_execution',
        '/api/sys/job/execution/detail',
        'POST'
    ),
    (
        55,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'job_execution',
        '/api/sys/job/execution/cancel',
        'POST'
    ),
    (
        56,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'job_execution',
        '/api/sys/job/execution/retry',
        'POST'
    ),
    (
        57,
        NULL,
        NULL,
        0,
        0,
        0,
        '',
        0,
        TRUE,
        0,
        'job_schedule',
        '/api/sys/job/schedule/options',
        'POST'
    );

INSERT INTO
    "public"."sys_resource_menu_api" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "menu_id",
        "api_id",
        "deleted_at"
    )
VALUES
    (1, NULL, NULL, 0, 0, 0, 8, 1, 0),
    (2, NULL, NULL, 0, 0, 0, 8, 2, 0),
    (3, NULL, NULL, 0, 0, 0, 8, 3, 0),
    (4, NULL, NULL, 0, 0, 0, 8, 4, 0),
    (5, NULL, NULL, 0, 0, 0, 8, 5, 0),
    (6, NULL, NULL, 0, 0, 0, 8, 6, 0),
    (7, NULL, NULL, 0, 0, 0, 8, 7, 0),
    (8, NULL, NULL, 0, 0, 0, 11, 8, 0),
    (9, NULL, NULL, 0, 0, 0, 11, 9, 0),
    (10, NULL, NULL, 0, 0, 0, 11, 10, 0),
    (11, NULL, NULL, 0, 0, 0, 11, 11, 0),
    (12, NULL, NULL, 0, 0, 0, 6, 12, 0),
    (13, NULL, NULL, 0, 0, 0, 6, 13, 0),
    (14, NULL, NULL, 0, 0, 0, 6, 14, 0),
    (15, NULL, NULL, 0, 0, 0, 6, 15, 0),
    (16, NULL, NULL, 0, 0, 0, 6, 16, 0),
    (17, NULL, NULL, 0, 0, 0, 6, 17, 0),
    (18, NULL, NULL, 0, 0, 0, 6, 18, 0),
    (19, NULL, NULL, 0, 0, 0, 6, 19, 0),
    (20, NULL, NULL, 0, 0, 0, 6, 20, 0),
    (21, NULL, NULL, 0, 0, 0, 6, 21, 0),
    (22, NULL, NULL, 0, 0, 0, 5, 22, 0),
    (23, NULL, NULL, 0, 0, 0, 5, 23, 0),
    (24, NULL, NULL, 0, 0, 0, 5, 24, 0),
    (25, NULL, NULL, 0, 0, 0, 5, 25, 0),
    (26, NULL, NULL, 0, 0, 0, 5, 26, 0),
    (27, NULL, NULL, 0, 0, 0, 5, 27, 0),
    (28, NULL, NULL, 0, 0, 0, 5, 28, 0),
    (29, NULL, NULL, 0, 0, 0, 5, 29, 0),
    (30, NULL, NULL, 0, 0, 0, 5, 30, 0),
    (31, NULL, NULL, 0, 0, 0, 10, 31, 0),
    (32, NULL, NULL, 0, 0, 0, 10, 32, 0),
    (33, NULL, NULL, 0, 0, 0, 3, 33, 0),
    (34, NULL, NULL, 0, 0, 0, 3, 34, 0),
    (35, NULL, NULL, 0, 0, 0, 3, 35, 0),
    (36, NULL, NULL, 0, 0, 0, 3, 36, 0),
    (37, NULL, NULL, 0, 0, 0, 3, 37, 0),
    (38, NULL, NULL, 0, 0, 0, 9, 38, 0),
    (39, NULL, NULL, 0, 0, 0, 9, 39, 0),
    (40, NULL, NULL, 0, 0, 0, 9, 40, 0),
    (41, NULL, NULL, 0, 0, 0, 9, 41, 0),
    (42, NULL, NULL, 0, 0, 0, 12, 42, 0),
    (43, NULL, NULL, 0, 0, 0, 12, 43, 0),
    (44, NULL, NULL, 0, 0, 0, 14, 45, 0),
    (45, NULL, NULL, 0, 0, 0, 14, 46, 0),
    (46, NULL, NULL, 0, 0, 0, 14, 47, 0),
    (47, NULL, NULL, 0, 0, 0, 14, 48, 0),
    (48, NULL, NULL, 0, 0, 0, 14, 49, 0),
    (49, NULL, NULL, 0, 0, 0, 14, 50, 0),
    (50, NULL, NULL, 0, 0, 0, 14, 51, 0),
    (51, NULL, NULL, 0, 0, 0, 14, 52, 0),
    (52, NULL, NULL, 0, 0, 0, 14, 57, 0),
    (53, NULL, NULL, 0, 0, 0, 15, 53, 0),
    (54, NULL, NULL, 0, 0, 0, 15, 54, 0),
    (55, NULL, NULL, 0, 0, 0, 15, 55, 0),
    (56, NULL, NULL, 0, 0, 0, 15, 56, 0);

INSERT INTO
    "public"."casbin_rule" ("ptype", "v0", "v1", "v2")
VALUES
    ('p', 'role:root', '/api/sys/role/list', 'POST'),
    ('p', 'role:root', '/api/sys/role/tree', 'GET'),
    (
        'p',
        'role:root',
        '/api/sys/role/:id/permissions',
        'GET'
    ),
    ('p', 'role:root', '/api/sys/role/create', 'POST'),
    ('p', 'role:root', '/api/sys/role/update', 'POST'),
    ('p', 'role:root', '/api/sys/role/del', 'POST'),
    (
        'p',
        'role:root',
        '/api/sys/role/permissions',
        'POST'
    ),
    ('p', 'role:root', '/api/sys/user/list', 'POST'),
    ('p', 'role:root', '/api/sys/user/create', 'POST'),
    ('p', 'role:root', '/api/sys/user/update', 'POST'),
    ('p', 'role:root', '/api/sys/user/del', 'POST'),
    (
        'p',
        'role:root',
        '/api/sys/dict/type/list',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/dict/type/create',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/dict/type/update',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/dict/type/del',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/dict/entry/list',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/dict/entry/match',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/dict/entry/create',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/dict/entry/update',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/dict/entry/del',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/dict/entry/batch/copy',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/language/type/list',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/language/type/create',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/language/type/update',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/language/type/del',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/language/entry/list',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/language/entry/create',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/language/entry/update',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/language/entry/del',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/language/entry/batch/create',
        'POST'
    ),
    ('p', 'role:root', '/api/sys/api/log/list', 'POST'),
    (
        'p',
        'role:root',
        '/api/sys/api/log/detail',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/resource/menu/list',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/resource/menu/tree',
        'GET'
    ),
    (
        'p',
        'role:root',
        '/api/sys/resource/menu/create',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/resource/menu/update',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/resource/menu/del',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/resource/api/list',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/resource/api/create',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/resource/api/update',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/resource/api/del',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/login/log/list',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/login/log/detail',
        'POST'
    ),
    ('p', 'role:root', '/api/events', 'GET'),
    (
        'p',
        'role:root',
        '/api/sys/job/schedule/list',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/job/schedule/detail',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/job/schedule/create',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/job/schedule/update',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/job/schedule/del',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/job/schedule/switch',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/job/schedule/sync',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/job/schedule/trigger',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/job/execution/list',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/job/execution/detail',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/job/execution/cancel',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/job/execution/retry',
        'POST'
    ),
    (
        'p',
        'role:root',
        '/api/sys/job/schedule/options',
        'POST'
    );

INSERT INTO
    "public"."sys_role_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "role_id",
        "menu_id",
        "deleted_at"
    )
VALUES
    (
        1,
        '2026-05-12 16:54:34.859724+08',
        '2026-05-12 16:54:34.859724+08',
        1,
        1,
        0,
        1,
        2,
        0
    );

INSERT INTO
    "public"."sys_role_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "role_id",
        "menu_id",
        "deleted_at"
    )
VALUES
    (
        2,
        '2026-05-12 16:54:34.859724+08',
        '2026-05-12 16:54:34.859724+08',
        1,
        1,
        0,
        1,
        3,
        0
    );

INSERT INTO
    "public"."sys_role_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "role_id",
        "menu_id",
        "deleted_at"
    )
VALUES
    (
        3,
        '2026-05-12 16:54:34.859724+08',
        '2026-05-12 16:54:34.859724+08',
        1,
        1,
        0,
        1,
        4,
        0
    );

INSERT INTO
    "public"."sys_role_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "role_id",
        "menu_id",
        "deleted_at"
    )
VALUES
    (
        4,
        '2026-05-12 16:54:34.859724+08',
        '2026-05-12 16:54:34.859724+08',
        1,
        1,
        0,
        1,
        5,
        0
    );

INSERT INTO
    "public"."sys_role_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "role_id",
        "menu_id",
        "deleted_at"
    )
VALUES
    (
        5,
        '2026-05-12 16:54:34.859724+08',
        '2026-05-12 16:54:34.859724+08',
        1,
        1,
        0,
        1,
        6,
        0
    );

INSERT INTO
    "public"."sys_role_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "role_id",
        "menu_id",
        "deleted_at"
    )
VALUES
    (
        6,
        '2026-05-12 16:54:34.859724+08',
        '2026-05-12 16:54:34.859724+08',
        1,
        1,
        0,
        1,
        7,
        0
    );

INSERT INTO
    "public"."sys_role_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "role_id",
        "menu_id",
        "deleted_at"
    )
VALUES
    (
        7,
        '2026-05-12 16:54:34.859724+08',
        '2026-05-12 16:54:34.859724+08',
        1,
        1,
        0,
        1,
        8,
        0
    );

INSERT INTO
    "public"."sys_role_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "role_id",
        "menu_id",
        "deleted_at"
    )
VALUES
    (
        8,
        '2026-05-12 16:54:34.859724+08',
        '2026-05-12 16:54:34.859724+08',
        1,
        1,
        0,
        1,
        9,
        0
    );

INSERT INTO
    "public"."sys_role_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "role_id",
        "menu_id",
        "deleted_at"
    )
VALUES
    (
        9,
        '2026-05-12 16:54:34.859724+08',
        '2026-05-12 16:54:34.859724+08',
        1,
        1,
        0,
        1,
        10,
        0
    );

INSERT INTO
    "public"."sys_role_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "role_id",
        "menu_id",
        "deleted_at"
    )
VALUES
    (
        10,
        '2026-05-12 16:54:34.859724+08',
        '2026-05-12 16:54:34.859724+08',
        1,
        1,
        0,
        1,
        11,
        0
    );

INSERT INTO
    "public"."sys_role_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "role_id",
        "menu_id",
        "deleted_at"
    )
VALUES
    (
        11,
        '2026-05-13 00:00:00+08',
        '2026-05-13 00:00:00+08',
        1,
        1,
        0,
        1,
        12,
        0
    );

INSERT INTO
    "public"."sys_role_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "role_id",
        "menu_id",
        "deleted_at"
    )
VALUES
    (
        12,
        '2026-05-14 00:00:00+08',
        '2026-05-14 00:00:00+08',
        1,
        1,
        0,
        1,
        13,
        0
    );

INSERT INTO
    "public"."sys_role_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "role_id",
        "menu_id",
        "deleted_at"
    )
VALUES
    (
        13,
        '2026-05-14 00:00:00+08',
        '2026-05-14 00:00:00+08',
        1,
        1,
        0,
        1,
        14,
        0
    );

INSERT INTO
    "public"."sys_role_menu" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "role_id",
        "menu_id",
        "deleted_at"
    )
VALUES
    (
        14,
        '2026-05-14 00:00:00+08',
        '2026-05-14 00:00:00+08',
        1,
        1,
        0,
        1,
        15,
        0
    );

INSERT INTO
    "public"."sys_role_api" (
        "id",
        "created_at",
        "updated_at",
        "created_by",
        "updated_by",
        "deleted_by",
        "role_id",
        "api_id",
        "deleted_at"
    )
VALUES
    (1, NULL, NULL, 0, 0, 0, 1, 1, 0),
    (2, NULL, NULL, 0, 0, 0, 1, 2, 0),
    (3, NULL, NULL, 0, 0, 0, 1, 3, 0),
    (4, NULL, NULL, 0, 0, 0, 1, 4, 0),
    (5, NULL, NULL, 0, 0, 0, 1, 5, 0),
    (6, NULL, NULL, 0, 0, 0, 1, 6, 0),
    (7, NULL, NULL, 0, 0, 0, 1, 7, 0),
    (8, NULL, NULL, 0, 0, 0, 1, 8, 0),
    (9, NULL, NULL, 0, 0, 0, 1, 9, 0),
    (10, NULL, NULL, 0, 0, 0, 1, 10, 0),
    (11, NULL, NULL, 0, 0, 0, 1, 11, 0),
    (12, NULL, NULL, 0, 0, 0, 1, 12, 0),
    (13, NULL, NULL, 0, 0, 0, 1, 13, 0),
    (14, NULL, NULL, 0, 0, 0, 1, 14, 0),
    (15, NULL, NULL, 0, 0, 0, 1, 15, 0),
    (16, NULL, NULL, 0, 0, 0, 1, 16, 0),
    (17, NULL, NULL, 0, 0, 0, 1, 17, 0),
    (18, NULL, NULL, 0, 0, 0, 1, 18, 0),
    (19, NULL, NULL, 0, 0, 0, 1, 19, 0),
    (20, NULL, NULL, 0, 0, 0, 1, 20, 0),
    (21, NULL, NULL, 0, 0, 0, 1, 21, 0),
    (22, NULL, NULL, 0, 0, 0, 1, 22, 0),
    (23, NULL, NULL, 0, 0, 0, 1, 23, 0),
    (24, NULL, NULL, 0, 0, 0, 1, 24, 0),
    (25, NULL, NULL, 0, 0, 0, 1, 25, 0),
    (26, NULL, NULL, 0, 0, 0, 1, 26, 0),
    (27, NULL, NULL, 0, 0, 0, 1, 27, 0),
    (28, NULL, NULL, 0, 0, 0, 1, 28, 0),
    (29, NULL, NULL, 0, 0, 0, 1, 29, 0),
    (30, NULL, NULL, 0, 0, 0, 1, 30, 0),
    (31, NULL, NULL, 0, 0, 0, 1, 31, 0),
    (32, NULL, NULL, 0, 0, 0, 1, 32, 0),
    (33, NULL, NULL, 0, 0, 0, 1, 33, 0),
    (34, NULL, NULL, 0, 0, 0, 1, 34, 0),
    (35, NULL, NULL, 0, 0, 0, 1, 35, 0),
    (36, NULL, NULL, 0, 0, 0, 1, 36, 0),
    (37, NULL, NULL, 0, 0, 0, 1, 37, 0),
    (38, NULL, NULL, 0, 0, 0, 1, 38, 0),
    (39, NULL, NULL, 0, 0, 0, 1, 39, 0),
    (40, NULL, NULL, 0, 0, 0, 1, 40, 0),
    (41, NULL, NULL, 0, 0, 0, 1, 41, 0),
    (42, NULL, NULL, 0, 0, 0, 1, 42, 0),
    (43, NULL, NULL, 0, 0, 0, 1, 43, 0),
    (44, NULL, NULL, 0, 0, 0, 1, 44, 0),
    (45, NULL, NULL, 0, 0, 0, 1, 45, 0),
    (46, NULL, NULL, 0, 0, 0, 1, 46, 0),
    (47, NULL, NULL, 0, 0, 0, 1, 47, 0),
    (48, NULL, NULL, 0, 0, 0, 1, 48, 0),
    (49, NULL, NULL, 0, 0, 0, 1, 49, 0),
    (50, NULL, NULL, 0, 0, 0, 1, 50, 0),
    (51, NULL, NULL, 0, 0, 0, 1, 51, 0),
    (52, NULL, NULL, 0, 0, 0, 1, 52, 0),
    (53, NULL, NULL, 0, 0, 0, 1, 53, 0),
    (54, NULL, NULL, 0, 0, 0, 1, 54, 0),
    (55, NULL, NULL, 0, 0, 0, 1, 55, 0),
    (56, NULL, NULL, 0, 0, 0, 1, 56, 0),
    (57, NULL, NULL, 0, 0, 0, 1, 57, 0);

INSERT INTO
    job_schedule (
        job_code,
        job_name,
        workflow_type,
        task_queue,
        schedule_type,
        interval_seconds,
        input_json,
        status,
        temporal_schedule_id,
        temporal_workflow_id_prefix,
        description
    )
VALUES
    (
        'print_count_10s',
        '每10秒打印count',
        'PrintCountWorkflow',
        'admin',
        'INTERVAL',
        10,
        '{"count":1}',
        'ENABLED',
        'print_count_10s',
        'print-count',
        '示例任务：每10秒运行一次 PrintCountWorkflow'
    );
