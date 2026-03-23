// Package cli 定义命令行接口
package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/opskit/opskit/internal/embed"
)

var (
	// Version OpsKit 版本号，编译时设置
	Version = "dev"
	// BuildTime 编译时间，编译时设置
	BuildTime = "unknown"
	// Commit Git 提交哈希，编译时设置
	Commit = "unknown"

	verbose bool // 是否输出详细调试信息
	binDir  string // 自定义二进制文件缓存目录
)

// NewRootCommand 创建根命令
//
// 返回:
//   - *cobra.Command: 根命令实例
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "opskit",
		Short: "OpsKit - 嵌入式运维工具箱",
		Long: `OpsKit 是一个单二进制、零依赖、内置常用运维工具的 CLI 工具箱，专为受限环境下的故障排查设计。

它提供 jq、curl、yq 和 busybox 等工具，所有工具集成在一个文件中。
非常适合无法安装工具的受限环境使用。

使用方法:
  opskit <工具名> [参数...]    运行嵌入式工具
  opskit <命令>               运行管理命令

管理命令:
  list      列出所有内置工具
  version   显示 OpsKit 版本信息
  extract   将工具释放到本地目录
  clean     清理缓存的二进制文件
  which     显示工具缓存路径

示例:
  opskit jq '.name' data.json   - 处理 JSON 数据
  opskit curl https://example.com - 发送 HTTP 请求
  opskit list                  - 列出所有可用工具`,
		SilenceUsage:  true,
		SilenceErrors: true,
		// 完全禁用 Cobra 的参数解析，自己处理所有逻辑
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 处理全局标志
			remainingArgs := []string{}
			for i := 0; i < len(args); i++ {
				arg := args[i]
				switch arg {
				case "-v", "--verbose":
					verbose = true
				case "--bin-dir":
					if i+1 < len(args) {
						binDir = args[i+1]
						i++ // 跳过下一个参数
					}
				case "--version":
					printVersion()
					return nil
				case "-h", "--help":
					return cmd.Help()
				default:
					remainingArgs = append(remainingArgs, arg)
				}
			}

			args = remainingArgs

			if len(args) == 0 {
				return cmd.Help()
			}

			// 处理内置命令
			switch args[0] {
			case "version":
				printVersion()
				return nil
			case "list":
				return runList()
			case "extract":
				if len(args) < 2 {
					return fmt.Errorf("使用方法: opskit extract <工具名> [目标目录]")
				}
				dest := "."
				if len(args) >= 3 {
					dest = args[2]
				}
				return runExtract(args[1], dest)
			case "clean":
				return runClean()
			case "which":
				if len(args) < 2 {
					return fmt.Errorf("使用方法: opskit which <工具名>")
				}
				return runWhich(args[1])
			}

			// 否则作为工具执行
			return runTool(args[0], args[1:])
		},
	}

	return cmd
}

// printVersion 打印版本信息
func printVersion() {
	runVersion()
}

// runTool 执行工具
//
// 参数:
//   - toolName: 工具名称
//   - args: 工具参数
//
// 返回:
//   - error: 执行过程中可能的错误
func runTool(toolName string, args []string) error {
	mgr, err := embed.NewManager(binDir, verbose)
	if err != nil {
		return fmt.Errorf("初始化管理器失败: %w", err)
	}

	executor := embed.NewExecutor(mgr)

	tools, _ := mgr.ListTools()
	for _, t := range tools {
		for _, p := range t.Provides {
			if p == toolName && t.Name == "busybox" {
				return executor.ExecuteBusybox(toolName, args)
			}
		}
	}

	return executor.Execute(toolName, args)
}
