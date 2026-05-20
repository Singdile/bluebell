package mysql

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

// ScoreUpdateItem score更新映射结构体
type ScoreUpdateItem struct {
	PostID int64
	Delta  int64
}

// ScoreQueue 更新队列
type ScoreQueue struct {
	mu    sync.Mutex        //约定的锁，用于保护下面的字段
	items []ScoreUpdateItem //修改队列

	// 下面是初始化后的只读配置，不需要锁保护
	batchsize int           //批量更新大小
	interval  time.Duration //定时写入间隔时间
	stopChan  chan struct{} //停止信号
}

// 全局的更新队列
var scoreQueue *ScoreQueue

// InitScoreQueue 初始化 Score 队列
func InitScoreQueue() {
	scoreQueue = &ScoreQueue{
		items:     make([]ScoreUpdateItem, 0),
		batchsize: 100,
		interval:  5 * time.Second,
		stopChan:  make(chan struct{}),
	}
	go scoreQueue.consumer() //开启协程，异步执行写入mysql

}

// PushScoreUpdate 添加 score 到更新队列
func PushScoreUpdate(postID, delta int64) {
	//判断队列有没有初始化
	if scoreQueue == nil {
		zap.L().Error("Score queue not initialized")
		return
	}
	//加锁
	scoreQueue.mu.Lock()

	//修改更新队列
	scoreQueue.items = append(scoreQueue.items, ScoreUpdateItem{PostID: postID, Delta: delta})
	//释放锁
	scoreQueue.mu.Unlock()

	//判断是否达到写入批量大小 batchsize
	if len(scoreQueue.items) >= scoreQueue.batchsize {
		// 写入mysql
		go scoreQueue.flush() // 开启协程，异步执行写入mysql
	}
}

// consumer 后台消费者，定时写入
func (q *ScoreQueue) consumer() {
	ticker := time.NewTicker(q.interval)
	defer ticker.Stop()

	for {
		select { //select 是专用与chan 的控制流关键字，下面的chan 谁不阻塞，就先执行谁
		case <-ticker.C:
			q.flush() // 执行写入
		case <-q.stopChan:
			q.flush() //执行写入
			return    //退出
		}
	}
}

// flush 批量写入score到mysql中去
func (q *ScoreQueue) flush() {
	// 给队列上锁
	q.mu.Lock()

	// 取出待处理项，并释放锁
	items := q.items
	q.items = make([]ScoreUpdateItem, 0)
	q.mu.Unlock()

	// 合并相同 post_id 的更新
	merged := make(map[int64]int64)
	for _, item := range items {
		merged[item.PostID] += item.Delta
	}
	if len(merged) == 0 { //到时间，但是没有更新，提前返回
		return
	}

	// 构建单条批量 UPDATE SQL
	// UPDATE post SET score = score + id
	// WHEN post_id_1 THEN delta_1
	// WHEN ? THEN ?
	// END
	// WHERE id IN (post_id_1,?)
	sqlstr := "UPDATE post SET score = score + CASE post_id "
	args := make([]interface{}, 0, len(merged)*3) //每一个更新，对应三个？
	ids := make([]interface{}, 0, len(merged))    //对应于WHERE IN ()

	// WHEN THEN 部分
	for postid, delta := range merged {
		sqlstr += " WHEN ? THEN ?"
		args = append(args, postid, delta)
		ids = append(ids, postid)
	}
	sqlstr += " END WHERE post_id IN ("

	// WHERE IN 部分
	for i := range ids {
		if i > 0 {
			sqlstr += ","
		}
		sqlstr += "?"
	}
	sqlstr += ")"

	args = append(args, ids...)

	// 执行sql语句
	if _, err := db.Exec(sqlstr, args...); err != nil {
		zap.L().Error("Failed to exec a bulk update post scores", zap.Error(err))
	}
	zap.L().Debug("Score queue flushed", zap.Int("count", len(merged)))
}

// StopScoreQueue 关闭更新队列
func StopScoreQueue() {
	if scoreQueue != nil {
		close(scoreQueue.stopChan)
	}
}
