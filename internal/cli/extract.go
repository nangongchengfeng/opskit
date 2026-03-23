package cli

import (
	"fmt"
	"os"

	"github.com/opskit/opskit/internal/embed"
)

var (
	extractAll bool
	extractDir string
)


func runExtract(toolName string, destDir string) error {
	mgr, err := embed.NewManager(binDir, verbose)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	path, err := mgr.ExtractTo(toolName, destDir)
	if err != nil {
		return fmt.Errorf("failed to extract %s: %w", toolName, err)
	}

	fmt.Printf("✓ %s 已释放至 %s\n", toolName, path)
	return nil
}
