//go:build embed_assets
// +build embed_assets

package embed

import (
	"embed"
	"fmt"
	"runtime"
)

// 直接嵌入每个文件 - 不使用通配符

//go:embed assets/linux-amd64/jq
var jqAmd64 []byte

//go:embed assets/linux-amd64/curl
var curlAmd64 []byte

//go:embed assets/linux-amd64/yq
var yqAmd64 []byte

//go:embed assets/linux-amd64/busybox
var busyboxAmd64 []byte

//go:embed assets/linux-arm64/jq
var jqArm64 []byte

//go:embed assets/linux-arm64/curl
var curlArm64 []byte

//go:embed assets/linux-arm64/yq
var yqArm64 []byte

//go:embed assets/linux-arm64/busybox
var busyboxArm64 []byte

func (m *BinaryManager) readEmbeddedAsset(toolName string) ([]byte, error) {
	platform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)

	if platform == "linux-amd64" {
		switch toolName {
		case "jq":
			return jqAmd64, nil
		case "curl":
			return curlAmd64, nil
		case "yq":
			return yqAmd64, nil
		case "busybox":
			return busyboxAmd64, nil
		}
	}

	if platform == "linux-arm64" {
		switch toolName {
		case "jq":
			return jqArm64, nil
		case "curl":
			return curlArm64, nil
		case "yq":
			return yqArm64, nil
		case "busybox":
			return busyboxArm64, nil
		}
	}

	return nil, fmt.Errorf("tool %s not found for platform %s", toolName, platform)
}
