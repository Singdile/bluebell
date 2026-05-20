-- 删除冗余的community_id,使用简单的自增id作为社区id
ALTER TABLE community DROP INDEX idx_community_id;

ALTER TABLE community DROP COLUMN community_id;
