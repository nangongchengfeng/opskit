//go:build !embed_assets
// +build !embed_assets

package embed

import "fmt"

// readEmbeddedAsset returns error when embed is not enabled
func (m *BinaryManager) readEmbeddedAsset(toolName string) ([]byte, error) {
	return nil, fmt.Errorf("embedded assets not enabled in this build")
}
