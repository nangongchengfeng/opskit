package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/opskit/opskit/internal/embed"
)

// NewCleanCommand creates the clean command
func NewCleanCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Clean cached binaries",
		Long:  "Remove all cached binaries from the opskit cache directory.",
		RunE:  runClean,
	}
}

func runClean(cmd *cobra.Command, args []string) error {
	mgr, err := embed.NewManager(binDir, verbose)
	if err != nil {
		return err
	}

	cacheDir := mgr.CacheDir()
	if verbose {
		fmt.Printf("Cleaning cache directory: %s\n", cacheDir)
	}

	if err := mgr.Clean(); err != nil {
		return fmt.Errorf("failed to clean cache: %w", err)
	}

	fmt.Printf("✓ 已清理缓存: %s\n", cacheDir)
	return nil
}
