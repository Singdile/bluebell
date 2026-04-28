-- 修改email 字段，添加默认值为空字符串
ALTER TABLE `user`
MODIFY COLUMN `email` varchar(64) NOT NULL DEFAULT '' COMMENT '用户邮箱';

UPDATE `user` SET `email` = '' WHERE `email` IS NULL;
