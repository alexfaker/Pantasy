package main

import (
	"github.com/alexfaker/Pantasy/config"
	"github.com/alexfaker/Pantasy/router"
)

func main() {

	config.Initialize()

	router.Run()
}

// NodeCode 获取节点随机唯一编号
