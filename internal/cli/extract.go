// Package cli 定义命令行接口
package cli

import (
	"fmt"
	"os"

	"github.com/opskit/opskit/internal/embed"
)

var (
	extractAll bool   // 是否提取所有工具
	extractDir string // 提取目标目录
)

// runExtract 将工具二进制文件提取到指定目录
//
// 参数:
//   - toolName: 工具名称
//   - destDir: 目标目录
//
// 返回:
//   - error: 提取过程中可能的错误
func runExtract(toolName string, destDir string) error {
	mgr, err := embed.NewManager(binDir, verbose)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	path, err := mgr.ExtractTo(toolName, destDir)
	if err != nil {
		return fmt.Errorf("提取工具 %s 失败: %w", toolName, err)
	}

	fmt.Printf("✓ %s 已释放至 %s\n", toolName, path)
	return nil
}
