// Package cli 定义命令行接口
package cli

import (
	"fmt"
	"runtime"
)

// runVersion 显示详细版本信息
//
// 输出内容包括:
//   - OpsKit 版本号
//   - 编译时间
//   - Git 提交哈希
//   - Go 版本
//   - 操作系统/架构
func runVersion() {
	fmt.Printf("OpsKit %s\n", Version)
	fmt.Printf("  编译时间: %s\n", BuildTime)
	fmt.Printf("  Git 提交: %s\n", Commit)
	fmt.Printf("  Go 版本: %s\n", runtime.Version())
	fmt.Printf("  系统/架构: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}
