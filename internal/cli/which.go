// Package cli 定义命令行接口
package cli

import (
	"fmt"

	"github.com/opskit/opskit/internal/embed"
)

// runWhich 显示工具的缓存路径
func runWhich(toolName string) error {
	mgr, err := embed.NewManager(binDir, verbose)
	if err != nil {
		return err
	}

	path, err := mgr.GetPath(toolName)
	if err != nil {
		return err
	}

	fmt.Println(path)
	return nil
}
