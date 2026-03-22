// Package cli defines the command-line interface
package cli

import (
	"fmt"
	"os"

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
	showVer bool
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
		Version:       Version,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if showVer {
				fmt.Printf("OpsKit %s\n", Version)
				fmt.Printf("Build: %s\n", BuildTime)
				fmt.Printf("Commit: %s\n", Commit)
				os.Exit(0)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			return runTool(args[0], args[1:])
		},
		DisableFlagParsing: false,
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
	}

	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	cmd.PersistentFlags().StringVar(&binDir, "bin-dir", "", "Directory to store/use cached binaries")
	cmd.PersistentFlags().BoolVar(&showVer, "version", false, "Show version information")

	cmd.AddCommand(NewListCommand())
	cmd.AddCommand(NewVersionCommand())
	cmd.AddCommand(NewExtractCommand())
	cmd.AddCommand(NewCleanCommand())
	cmd.AddCommand(NewWhichCommand())

	return cmd
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
