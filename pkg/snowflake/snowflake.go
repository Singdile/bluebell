package snowflake

import (
	sf "github.com/bwmarrin/snowflake"
	"time"
)

var node *sf.Node

// 初始化雪花算法的计算节点
func Init(startTime string, machineID int64) (err error) {
	var st time.Time

	//设置项目的启动时间，用于初始化雪花算法的时间戳部分
	st, err = time.Parse("2006-01-02", startTime) //解析启动时间点
	if err != nil {
		return
	}

	//设置库的全局部分，以后生成的ID时间戳的部分 = (当前毫秒 - sf.Epoch)
	sf.Epoch = st.UnixNano() / 1000000

	node, err = sf.NewNode(machineID)

	return
}

// 利用计算节点生成唯一ID
func GenID() int64 {
	return node.Generate().Int64()
}
