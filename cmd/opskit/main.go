// Package main 是 OpsKit 的主入口文件
package main

import (
	"os"

	"github.com/opskit/opskit/internal/cli"
)

var (
	Version   = "dev"     // 版本号，编译时设置
	BuildTime = "unknown" // 编译时间，编译时设置
	Commit    = "unknown" // Git 提交哈希，编译时设置
)

// main 函数是程序的入口点
func main() {
	// 注入编译时变量
	cli.Version = Version
	cli.BuildTime = BuildTime
	cli.Commit = Commit

	cmd := cli.NewRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
