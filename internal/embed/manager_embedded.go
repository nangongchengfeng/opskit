//go:build embed_assets
// +build embed_assets

package embed

import (
	"embed"
	"fmt"
	"runtime"
)

// Embedded assets - only included when build tag "embed_assets" is set
//go:embed assets
var assetsFS embed.FS

func (m *BinaryManager) readEmbeddedAsset(toolName string) ([]byte, error) {
	platform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	path := fmt.Sprintf("assets/%s/%s", platform, toolName)

	return assetsFS.ReadFile(path)
}
