CREATE TABLE `community` (
  `id` bigint  NOT NULL AUTO_INCREMENT COMMENT '内部ID（自增主键）',
  `community_id` bigint funsigned NOT NULL COMMENT '社区唯一标识ID（业务ID）',
  `community_name` varchar(128) COLLATE utf8mb4_general_ci NOT NULL COMMENT '社区名称',
  `introduction` varchar(256) COLLATE utf8mb4_general_ci NOT NULL COMMENT '社区简介',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_community_id` (`community_id`),
  UNIQUE KEY `idx_community_name` (`community_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='社区板块表';
