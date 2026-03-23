// Package cli 定义命令行接口
package cli

import (
	"fmt"

	"github.com/opskit/opskit/internal/embed"
)

// runClean 清理工具缓存目录
func runClean() error {
	mgr, err := embed.NewManager(binDir, verbose)
	if err != nil {
		return err
	}

	cacheDir := mgr.CacheDir()
	if verbose {
		fmt.Printf("正在清理缓存目录: %s\n", cacheDir)
	}

	if err := mgr.Clean(); err != nil {
		return fmt.Errorf("清理缓存失败: %w", err)
	}

	fmt.Printf("✓ 已清理缓存: %s\n", cacheDir)
	return nil
}
