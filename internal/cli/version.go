package cli

import (
	"fmt"
	"runtime"
)

// runVersion prints detailed version information
func runVersion() {
	fmt.Printf("OpsKit %s\n", Version)
	fmt.Printf("  Build Time: %s\n", BuildTime)
	fmt.Printf("  Commit: %s\n", Commit)
	fmt.Printf("  Go Version: %s\n", runtime.Version())
	fmt.Printf("  OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}
