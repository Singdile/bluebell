-- 为post添加投票信息字段 vote_p, vote_n
ALTER TABLE post ADD COLUMN vote_p INT DEFAULT 0 COMMENT '赞成票数';
ALTER TABLE post ADD COLUMN vote_n INT DEFAULT 0 COMMENT '反对票数';
