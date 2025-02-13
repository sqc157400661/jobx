CREATE TABLE `job_cron`
(
    `id`           int           NOT NULL AUTO_INCREMENT COMMENT '主键编码',
    `app_name`  varchar(50)            DEFAULT NULL COMMENT '应用名称',
    `tenant`    varchar(50)            DEFAULT NULL COMMENT '租户信息',
    `locker`    varchar(30)   NOT NULL DEFAULT '' COMMENT '锁拥有者',
    `entry_id`     int           NOT NULL DEFAULT '0' COMMENT '定时任务id',
    `spec`         varchar(30)   NOT NULL DEFAULT '' COMMENT '定时表达式',
    `exec_type`    varchar(20)   NOT NULL DEFAULT '' COMMENT '执行任务类型，如job、func、shell、http',
    `exec_content` varchar(1000) NOT NULL DEFAULT '' COMMENT '执行任务内容',
    `status`       varchar(50)   NOT NULL DEFAULT '' COMMENT '状态',
    `create_at`    int           NOT NULL DEFAULT '0' COMMENT '创建时间',
    `update_at`    int           NOT NULL DEFAULT '0' COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_entry` (`entry_id`) USING BTREE
) DEFAULT CHARSET = utf8mb4 COMMENT='定时任务表';


CREATE TABLE `job_logs`
(
    `id`        int NOT NULL AUTO_INCREMENT COMMENT '主键编码',
    `event_id`  int NOT NULL DEFAULT '0' COMMENT '任务id',
    `result`    text CHARACTER SET utf8mb4 COMMENT '结果',
    `create_at` int NOT NULL DEFAULT '0' COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_event` (`event_id`) USING BTREE
) DEFAULT CHARSET = utf8mb4 COMMENT='任务日志表';


CREATE TABLE `job_task`
(
    `id`        int            NOT NULL AUTO_INCREMENT COMMENT '主键编码',
    `job_id`    int            NOT NULL DEFAULT '0' COMMENT '任务id',
    `name`      varchar(100)   NOT NULL DEFAULT '' COMMENT '任务名称',
    `desc` varchar(255) NOT NULL DEFAULT '' comment '任务描述',
    `description`      varchar(255)   NOT NULL DEFAULT '' COMMENT '任务描述',
    `action`    varchar(30)    NOT NULL DEFAULT '' COMMENT '对象行为名称',
    `retry`     tinyint        NOT NULL DEFAULT '3' COMMENT '允许自动重试次数',
    `retries`   tinyint        NOT NULL DEFAULT '0' COMMENT '已经自动重试的次数',
    `pause`     int                     DEFAULT NULL COMMENT '是否允许暂停',
    `phase`     varchar(30)    NOT NULL DEFAULT '' COMMENT '状态控制',
    `status`    varchar(255)   NOT NULL DEFAULT '' COMMENT '展示控制和手工控制',
    `reason`    varchar(2000)  NOT NULL DEFAULT '' COMMENT '失败原因',
    `env`       varchar(500)   NOT NULL DEFAULT '' COMMENT '配置信息',
    `input`     varchar(6000)  NOT NULL DEFAULT '' COMMENT '入参',
    `output`    varchar(10000) NOT NULL DEFAULT '' COMMENT '出参',
    `context`   varchar(5000)  NOT NULL DEFAULT '' COMMENT '上下文参数',
    `create_at` int            NOT NULL DEFAULT '0' COMMENT '创建时间',
    `update_at` int            NOT NULL DEFAULT '0' COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_job` (`job_id`) USING BTREE
) DEFAULT CHARSET = utf8mb4 COMMENT='任务节点表';


CREATE TABLE `job_token`
(
    `id`        int          NOT NULL AUTO_INCREMENT COMMENT '主键编码',
    `root_id`   int          NOT NULL DEFAULT '0' COMMENT '根任务id',
    `token`     varchar(100) NOT NULL DEFAULT '' COMMENT '令牌',
    `create_at` int          NOT NULL DEFAULT '0' COMMENT '创建时间',
    `update_at` int          NOT NULL DEFAULT '0' COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_token` (`token`) USING BTREE
) DEFAULT CHARSET = utf8mb4 COMMENT='任务token表';


CREATE TABLE `job`
(
    `id`        int           NOT NULL AUTO_INCREMENT COMMENT '主键编码',
    `biz_id`    varchar(36)            DEFAULT NULL COMMENT '业务方产生的唯一id',
    `app_name`  varchar(50)            DEFAULT NULL COMMENT '应用名称',
    `tenant`    varchar(50)            DEFAULT NULL COMMENT '租户信息',
    `root_id`   int           NOT NULL DEFAULT '0' COMMENT '主/根任务id',
    `parent_id` int           NOT NULL DEFAULT '0' COMMENT '父级任务id',
    `level`     int           NOT NULL DEFAULT '0' COMMENT '当前任务层级',
    `path`      varchar(50)            DEFAULT NULL COMMENT '任务到主/根任务的路径',
    `runnable`  tinyint       NOT NULL DEFAULT '0' COMMENT '是否是用于pipeline的任务',
    `name`      varchar(100)  NOT NULL DEFAULT '' COMMENT '任务名称',
    `desc` varchar(255) NOT NULL DEFAULT '' comment '任务描述',
    `description`      varchar(255)  NOT NULL DEFAULT '' COMMENT '任务描述',
    `owner`     varchar(50)   NOT NULL DEFAULT '' COMMENT '任务归属人',
    `pause`     tinyint       NOT NULL DEFAULT '1' COMMENT '是否允许暂停',
    `locker`    varchar(30)   NOT NULL DEFAULT '' COMMENT '锁拥有者',
    `phase`     varchar(30)   NOT NULL DEFAULT '' COMMENT '状态控制',
    `status`    varchar(255)  NOT NULL DEFAULT '' COMMENT '展示控制和手工控制',
    `reason`    varchar(6000) NOT NULL DEFAULT '' COMMENT '失败原因',
    `env`       varchar(500)  NOT NULL DEFAULT '' COMMENT '配置信息',
    `input`     varchar(6000) NOT NULL DEFAULT '' COMMENT '入参',
    `create_at` int           NOT NULL DEFAULT '0' COMMENT '创建时间',
    `update_at` int           NOT NULL DEFAULT '0' COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_biz` (`biz_id`) USING BTREE,
    KEY `idx_root` (`root_id`) USING BTREE,
    KEY `idx_parent` (`parent_id`) USING BTREE
) DEFAULT CHARSET = utf8mb4 COMMENT='任务主表';


CREATE TABLE `job_definition` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(80) NOT NULL DEFAULT '' COMMENT '名称',
  `app_name`  varchar(50)            DEFAULT NULL COMMENT '应用名称',
  `tenant` varchar(50) NOT NULL DEFAULT '' COMMENT '租户信息',
  `sort` int unsigned NOT NULL COMMENT '排序',
  `pipelines` varchar(2000) NOT NULL DEFAULT '' COMMENT '流水线json格式',
  `retry` int unsigned NOT NULL COMMENT '设定重试次数',
  `input` varchar(500) NOT NULL DEFAULT '' COMMENT '预设参数',
  `env` varchar(500) NOT NULL DEFAULT '' COMMENT '附加的env参数',
  `condition` varchar(100) NOT NULL DEFAULT '' COMMENT '处理特殊逻辑',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_name` (`name`)
) DEFAULT CHARSET=utf8mb4  COMMENT='通用任务定义表';
