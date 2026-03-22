package cli

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// NewVersionCommand creates the version command
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show opskit version",
		Long:  "Show opskit version and build information.",
		Run:   runVersion,
	}
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("OpsKit %s\n", Version)
	fmt.Printf("  Build Time: %s\n", BuildTime)
	fmt.Printf("  Commit: %s\n", Commit)
	fmt.Printf("  Go Version: %s\n", runtime.Version())
	fmt.Printf("  OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}
