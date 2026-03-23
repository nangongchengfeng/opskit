//go:build embed_assets
// +build embed_assets

package embed

import (
	"embed"
	"fmt"
	"runtime"
)

// Embedded assets - only included when build tag "embed_assets" is set
// 嵌入的资源文件 - 仅在编译时设置 "embed_assets" 标签时才会包含

//go:embed assets
var assetsFS embed.FS

// readEmbeddedAsset 从嵌入的资源中读取工具二进制文件
//
// 参数:
//   - toolName: 工具名称
//
// 返回:
//   - []byte: 二进制文件内容
//   - error: 读取过程中可能的错误
func (m *BinaryManager) readEmbeddedAsset(toolName string) ([]byte, error) {
	platform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	path := fmt.Sprintf("assets/%s/%s", platform, toolName)

	return assetsFS.ReadFile(path)
}
