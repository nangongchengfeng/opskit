package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/opskit/opskit/internal/embed"
)


func runList() error {
	mgr, err := embed.NewManager(binDir, verbose)
	if err != nil {
		return err
	}

	tools, err := mgr.ListTools()
	if err != nil {
		return err
	}

	fmt.Printf("\nOpsKit %s — 内置工具列表\n\n", Version)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TOOL\tVERSION\tDESCRIPTION")
	fmt.Fprintln(w, "----\t-------\t-----------")

	for _, t := range tools {
		fmt.Fprintf(w, "%s\t%s\t%s\n", t.Name, t.Version, t.Description)
		for _, p := range t.Provides {
			fmt.Fprintf(w, "  ↳ %s\t\t(via %s)\n", p, t.Name)
		}
	}

	w.Flush()

	fmt.Printf("\n缓存目录: %s\n", mgr.CacheDir())
	fmt.Println("\n提示: 使用 'opskit extract <tool>' 可将工具释放到本地目录单独使用。")

	return nil
}
