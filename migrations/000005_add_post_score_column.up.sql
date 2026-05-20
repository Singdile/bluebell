 -- 添加
ALTER TABLE post ADD COLUMN score BIGINT DEFAULT 0 COMMENT '帖子分数';

-- 初始化现有帖子的 score（基于创建时间）
UPDATE post SET score = UNIX_TIMESTAMP(create_time);
