package main

import (
	"Jinx/jnet"
	"fmt"
)

func main() {
	// 创建jinx的句柄
	s := jnet.NewServer()
	fmt.Println(s)

	// 链接创建和销毁的hook函数
	// 注册路由
	// 启动服务
	s.Serve()
}
