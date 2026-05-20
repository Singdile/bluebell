package snowflake

import (
	"testing"
)

// TestGenID 测试能否生成ID
func TestGenID(t *testing.T) {
	//初始化节点
	err := Init("2026-04-24", 1)
	if err != nil {
		t.Fatalf("Init failed, err :%v\n", err)
	}

	//执行函数
	id := GenID()

	//结果分析
	if id <= 0 {
		t.Errorf("ID should be positive, but got %d", id)
	}
}

// TestGenIDUnique 测试生成ID是否唯一
func TestGenIDUnique(t *testing.T) {
	//准备数据
	//初始化节点
	err := Init("2026-04-24", 1)
	if err != nil {
		t.Fatalf("Init failed, err :%v\n", err)
	}
	//执行函数
	ids := make(map[int64]bool)
	for _ = range 1000 {
		id := GenID()
		if ids[id] {
			t.Errorf("duplicate ID: %d", id)

		}
		ids[id] = true

	}

}
