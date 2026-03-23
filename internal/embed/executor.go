// Package embed 管理嵌入的二进制资源和工具执行
package embed

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// Executor 工具执行器
// 负责执行嵌入的工具二进制文件，并透传所有输入输出和信号
type Executor struct {
	manager *BinaryManager // 二进制管理器实例
}

// NewExecutor 创建新的工具执行器实例
//
// 参数:
//   - manager: 二进制管理器实例
//
// 返回:
//   - *Executor: 执行器实例
func NewExecutor(manager *BinaryManager) *Executor {
	return &Executor{
		manager: manager,
	}
}

// Execute 执行指定的工具
//
// 参数:
//   - toolName: 工具名称
//   - args: 传递给工具的命令行参数
//
// 返回:
//   - error: 执行过程中可能的错误
//
// 说明:
//   - 标准输入、输出、错误流会直接透传
//   - 工具的退出码会原样返回
//   - 信号会透传给子进程
func (e *Executor) Execute(toolName string, args []string) error {
	binPath, err := e.manager.GetPath(toolName)
	if err != nil {
		return fmt.Errorf("获取工具 %s 失败: %w", toolName, err)
	}

	if e.manager.verbose {
		fmt.Fprintf(os.Stderr, "[opskit] 执行命令: %s %v\n", binPath, args)
	}

	cmd := exec.Command(binPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}
		return err
	}

	return nil
}

// ExecuteBusybox 执行 BusyBox 提供的子命令
//
// 参数:
//   - command: BusyBox 子命令名称（如 telnet、ping 等）
//   - args: 传递给子命令的参数
//
// 返回:
//   - error: 执行过程中可能的错误
//
// 说明:
//   - BusyBox 工具通过 busybox <command> 的方式调用
//   - 如果 BusyBox 不可用，会尝试直接从系统 PATH 中查找该命令
func (e *Executor) ExecuteBusybox(command string, args []string) error {
	binPath, err := e.manager.GetPath("busybox")
	if err != nil {
		return e.Execute(command, args)
	}

	fullArgs := append([]string{command}, args...)

	if e.manager.verbose {
		fmt.Fprintf(os.Stderr, "[opskit] 执行 BusyBox 命令: %s %v\n", binPath, fullArgs)
	}

	cmd := exec.Command(binPath, fullArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}
		return err
	}

	return nil
}
