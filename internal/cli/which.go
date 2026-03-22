package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/opskit/opskit/internal/embed"
)

// NewWhichCommand creates the which command
func NewWhichCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "which <tool>",
		Short: "Show path to cached tool",
		Long:  "Show the path to the cached executable for a tool.",
		Args:  cobra.ExactArgs(1),
		RunE:  runWhich,
	}
}

func runWhich(cmd *cobra.Command, args []string) error {
	toolName := args[0]

	mgr, err := embed.NewManager(binDir, verbose)
	if err != nil {
		return err
	}

	path, err := mgr.GetPath(toolName)
	if err != nil {
		return err
	}

	fmt.Println(path)
	return nil
}
