-- 移除 email 字段的默认值
ALTER TABLE `user`
MODIFY COLUMN `email` varchar(64) NOT NULL COMMENT '用户邮箱';
