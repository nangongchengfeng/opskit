package embed

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// Executor handles executing embedded binaries
type Executor struct {
	manager *BinaryManager
}

// NewExecutor creates a new Executor
func NewExecutor(manager *BinaryManager) *Executor {
	return &Executor{
		manager: manager,
	}
}

// Execute executes a tool with the given arguments
func (e *Executor) Execute(toolName string, args []string) error {
	binPath, err := e.manager.GetPath(toolName)
	if err != nil {
		return fmt.Errorf("failed to get tool %s: %w", toolName, err)
	}

	if e.manager.verbose {
		fmt.Fprintf(os.Stderr, "[opskit] Executing: %s %v\n", binPath, args)
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

// ExecuteBusybox executes a busybox command
func (e *Executor) ExecuteBusybox(command string, args []string) error {
	binPath, err := e.manager.GetPath("busybox")
	if err != nil {
		// Try direct command if busybox not available
		return e.Execute(command, args)
	}

	fullArgs := append([]string{command}, args...)

	if e.manager.verbose {
		fmt.Fprintf(os.Stderr, "[opskit] Executing busybox: %s %v\n", binPath, fullArgs)
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
