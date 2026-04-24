CREATE TABLE    `user` (
       `id` bigint(20) NOT NULL AUTO_INCREMENT,
       `user_id` bigint(20) NOT NULL,
       `username` varchar(64)   NOT NULL,
       `password` varchar(64)  NOT NULL COMMENT '哈希加密之后的密码',
       `email` varchar(64)  NOT NULL,
       `gender` tinyint(4) NOT NULL DEFAULT '0' COMMENT '性别:0 未知, 1 男, 2 女',
       `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
       `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
       PRIMARY KEY (`id`),
       UNIQUE KEY `idx_username` (`username`) USING BTREE,
       UNIQUE KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
