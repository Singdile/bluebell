--撤销给 post 添加投票数据字段的操作
ALTER TABLE post DROP COLUMN vote_p;
ALTER TABLE post DROP COLUMN vote_n;
