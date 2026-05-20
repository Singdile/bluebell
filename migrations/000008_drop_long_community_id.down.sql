-- 撤销  删除冗余的community_id,使用简单的自增id作为社区id
ALTER TABLE community ADD COLUMN community_id bigint unsigned NOT NULL COMMENT '社区唯一标识ID（业务ID）' AFTER id;

-- 从 id 复制数据到 community_id
UPDATE community SET community_id = id;

-- 添加唯一索引
ALTER TABLE community ADD UNIQUE KEY idx_community_id (community_id);
