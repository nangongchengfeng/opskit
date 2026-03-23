package cli

import (
	"fmt"

	"github.com/opskit/opskit/internal/embed"
)


func runClean() error {
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
