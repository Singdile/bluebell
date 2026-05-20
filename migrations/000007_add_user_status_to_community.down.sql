-- 撤销 增添 创建者，社区状态 字段给community
ALTER TABLE community DROP COLUMN creator_id;
ALTER TABLE community DROP COLUMN status;
