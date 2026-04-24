package snowflake

import (
	"fmt"
	"testing"
)

func TestGenID(t *testing.T) {
	//初始化节点
	err := Init("2026-04-24", 1)
	if err != nil {
		t.Errorf("Init failed, err :%v\n", err)
	}

	//生成一个ID
	id := GenID()
	fmt.Printf("ID: %v\n", id)

	//断言
	if id == 0 {
		t.Error("生成的 ID 为0,逻辑有问题")
	}
}
