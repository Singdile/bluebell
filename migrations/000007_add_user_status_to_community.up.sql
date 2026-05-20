-- 增添 创建者，社区状态 字段给community
ALTER TABLE community ADD COLUMN creator_id bigint unsigned DEFAULT 0 COMMENT '创建者用户ID';
ALTER TABLE community ADD COLUMN status tinyint  DEFAULT 1 COMMENT '状态:1-正常 0-禁用';
