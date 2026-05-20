package mysql

import (
	"go.uber.org/zap"
	"sync"
	"time"
)

// VoteUpdate 投票更新数据
type VoteUpdateItem struct {
	PostID int64
	VoteP  int64
	VoteN  int64
}

// VoteUpdateQueue 投票队列，用于异步保存投票信息
type VoteUpdateQueue struct {
	mu    sync.Mutex       //用于保护实际的投票信息队列
	items []VoteUpdateItem //实际的投票信息队列

	batchsize int           //每批写入的大小
	interval  time.Duration //写入间隔时间
	stopChan  chan struct{} //停止信号
}

var voteQueue *VoteUpdateQueue

// InitVoteUpdateQueue 初始化投票信息队列
func InitVoteUpdateQueue() {
	voteQueue = &VoteUpdateQueue{
		items: make([]VoteUpdateItem, 0),

		batchsize: 100,
		interval:  5 * time.Second,
		stopChan:  make(chan struct{}),
	}
	go voteQueue.consumer()
}

// consumer 后台消费，用于定时/定量 写入
func (q *VoteUpdateQueue) consumer() {
	ticker := time.NewTicker(q.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C: //定时写入
			q.flush()
		case <-q.stopChan: //停止写入
			q.flush()
			return
		}
	}
}

// flush 批量写入投票信息到mysql表中去
func (q *VoteUpdateQueue) flush() {
	// 给队列上锁
	q.mu.Lock()
	// 取出队列中的待处理数据，并释放队列的锁
	items := q.items
	q.items = make([]VoteUpdateItem, 0) //清空队列
	q.mu.Unlock()

	// 合并: 保留每个 post_id 的最后一次更新
	merged := make(map[int64]VoteUpdateItem)
	for _, item := range items {
		merged[item.PostID] = item
	}
	if len(merged) == 0 { //没有修改
		return
	}

	// 构建sql语句
	sqlstr := "UPDATE post SET vote_p = ?, vote_n = ? WHERE post_id = ?"

	// 执行sql语句
	// 开启事务
	tx, err := db.Begin()
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return
	}
	defer tx.Rollback() //当commit提交成功，Rollback 回滚无效; 当不成功，则能自动回滚

	// 准备mysql 预处理模板，提高mysql处理性能
	stmt, err := tx.Prepare(sqlstr)
	if err != nil {
		zap.L().Error("Failed to prepare statement", zap.Error(err))
	}
	defer stmt.Close()

	// 执行sql
	for postid, item := range merged {
		_, err := stmt.Exec(item.VoteP, item.VoteN, postid)
		if err != nil {
			zap.L().Error("Failed to update post votes", zap.Error(err), zap.Int64("post_id", postid))
		}
	}

	if err := tx.Commit(); err != nil {
		zap.L().Error("Failed to commit update post voets commit transcation", zap.Error(err))
	}
	zap.L().Debug("Vote queue flushed", zap.Int("count", len(merged)))

}

// PushVoteUpdate 添加投票信息到更新队列
func PushVoteUpdate(postid, voteP, voteN int64) {
	// 判断队列有没有初始化
	if voteQueue == nil {
		zap.L().Error("vote queue not initialized")
		return
	}

	// 队列加锁
	voteQueue.mu.Lock()

	// 添加信息到队列
	voteQueue.items = append(voteQueue.items, VoteUpdateItem{
		PostID: postid,
		VoteP:  voteP,
		VoteN:  voteN,
	})

	// 释放锁
	voteQueue.mu.Unlock()

	// 判断是否达到batchsize，批量写入
	if len(voteQueue.items) >= voteQueue.batchsize {
		// 写入mysql
		go voteQueue.flush() //开启协程
	}
}

// StopVoteUpdateQueue 关闭更新队列
func StopVoteUpdateQueue() {
	if voteQueue != nil {
		close(voteQueue.stopChan)
	}
}
