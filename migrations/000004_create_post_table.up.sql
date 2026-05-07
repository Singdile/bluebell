CREATE TABLE `post` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '内部ID（自增主键）',
  `post_id` bigint NOT NULL COMMENT '帖子唯一标识ID（业务ID，雪花算法生成）',
  `title` varchar(128) COLLATE utf8mb4_general_ci NOT NULL COMMENT '帖子标题',
  `content` varchar(8192) COLLATE utf8mb4_general_ci NOT NULL COMMENT '帖子内容',
  `author_id` bigint NOT NULL COMMENT '作者ID（关联user.user_id）',
  `community_id` bigint unsigned NOT NULL COMMENT '所属社区ID（关联community.community_id）',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '帖子状态：1-正常，2-审核中，3-已删除，4-已封禁',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_post_id` (`post_id`),
  KEY `idx_author_id` (`author_id`),
  KEY `idx_community_id` (`community_id`),
  KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='帖子表';
