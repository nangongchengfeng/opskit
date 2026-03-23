//go:build !embed_assets
// +build !embed_assets

package embed

import "fmt"

// readEmbeddedAsset 当嵌入资源功能未启用时返回错误
//
// 参数:
//   - toolName: 工具名称
//
// 返回:
//   - nil, error: 总是返回错误，因为嵌入资源未启用
func (m *BinaryManager) readEmbeddedAsset(toolName string) ([]byte, error) {
	return nil, fmt.Errorf("此构建未启用嵌入资源功能")
}
