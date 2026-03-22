package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/opskit/opskit/internal/embed"
)

var (
	extractAll bool
	extractDir string
)

// NewExtractCommand creates the extract command
func NewExtractCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract [tool...]",
		Short: "Extract tools to local directory",
		Long:  `Extract embedded tools to a local directory for standalone use.`,
		RunE:  runExtract,
	}

	cmd.Flags().BoolVar(&extractAll, "all", false, "Extract all tools")
	cmd.Flags().StringVar(&extractDir, "dir", ".", "Directory to extract to")

	return cmd
}

func runExtract(cmd *cobra.Command, args []string) error {
	mgr, err := embed.NewManager(binDir, verbose)
	if err != nil {
		return err
	}

	var toolsToExtract []string

	if extractAll {
		tools, err := mgr.ListTools()
		if err != nil {
			return err
		}
		for _, t := range tools {
			toolsToExtract = append(toolsToExtract, t.Name)
		}
	} else if len(args) > 0 {
		toolsToExtract = args
	} else {
		return fmt.Errorf("specify at least one tool or use --all")
	}

	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	for _, tool := range toolsToExtract {
		path, err := mgr.ExtractTo(tool, extractDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "✗ %s: %v\n", tool, err)
			continue
		}
		fmt.Printf("✓ %s 已释放至 %s\n", tool, path)
	}

	return nil
}
