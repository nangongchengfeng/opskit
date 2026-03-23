// Package cli defines the command-line interface
package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/opskit/opskit/internal/embed"
)

var (
	// Version is set at build time
	Version = "dev"
	// BuildTime is set at build time
	BuildTime = "unknown"
	// Commit is set at build time
	Commit = "unknown"

	verbose bool
	binDir  string
)

// NewRootCommand creates the root command
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "opskit",
		Short: "OpsKit - Embedded operations toolbox",
		Long: `OpsKit is a single-binary toolbox with embedded common operations tools.

It provides tools like jq, curl, yq, and busybox utilities, all in one file.
Perfect for restricted environments where you can't install tools.

Usage:
  opskit <tool> [args...]    Run an embedded tool
  opskit <command>            Run a management command

Commands:
  list      List all embedded tools
  version   Show opskit version
  extract   Extract tools to local directory
  clean     Clean cached binaries
  which     Show path to cached tool

Examples:
  opskit jq '.name' data.json
  opskit curl https://example.com
  opskit list`,
		SilenceUsage:  true,
		SilenceErrors: true,
		// 完全禁用Cobra的参数解析，自己处理所有逻辑
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
					return fmt.Errorf("usage: opskit extract <tool> [destination]")
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
					return fmt.Errorf("usage: opskit which <tool>")
				}
				return runWhich(args[1])
			}

			// 否则作为工具执行
			return runTool(args[0], args[1:])
		},
	}

	return cmd
}

// printVersion prints version information
func printVersion() {
	runVersion()
}

// runTool executes a tool
func runTool(toolName string, args []string) error {
	mgr, err := embed.NewManager(binDir, verbose)
	if err != nil {
		return fmt.Errorf("failed to initialize manager: %w", err)
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
