// Package embed 管理嵌入的二进制资源和工具执行
package embed

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	globalManager *BinaryManager
	once          sync.Once
)

const (
	// Version 工具包版本号，用于缓存目录命名
	Version = "0.1.0"
)

// BinaryManager 二进制文件管理器
// 负责工具的缓存管理、路径获取、执行和释放
type BinaryManager struct {
	cacheDir string // 工具二进制文件的缓存目录
	verbose  bool   // 是否输出详细调试信息
	assetDir string // 存放二进制文件的资源目录
}

// ToolInfo 工具信息结构体
// 包含工具的名称、版本、描述和提供的子命令信息
type ToolInfo struct {
	Name        string   // 工具名称
	Version     string   // 工具版本
	Description string   // 工具描述
	Provides    []string // 该工具提供的子命令列表（如busybox提供多个命令）
}

// NewManager 创建新的二进制管理器实例
//
// 参数:
//   - cacheDir: 自定义工具缓存目录（为空则使用默认路径）
//   - verbose: 是否输出详细调试信息
//
// 返回:
//   - *BinaryManager: 管理器实例
//   - error: 创建过程中可能的错误
func NewManager(cacheDir string, verbose bool) (*BinaryManager, error) {
	if cacheDir == "" {
		var err error
		cacheDir, err = getDefaultCacheDir()
		if err != nil {
			return nil, fmt.Errorf("获取默认缓存目录失败: %w", err)
		}
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("创建缓存目录失败: %w", err)
	}

	assetDir, err := findAssetDir()
	if err != nil {
		return nil, fmt.Errorf("查找资源目录失败: %w", err)
	}

	return &BinaryManager{
		cacheDir: cacheDir,
		verbose:  verbose,
		assetDir: assetDir,
	}, nil
}

// GetManager 获取全局二进制管理器实例（单例模式）
func GetManager() (*BinaryManager, error) {
	var err error
	once.Do(func() {
		globalManager, err = NewManager("", false)
	})
	return globalManager, err
}

// SetVerbose 设置是否输出详细调试信息
func (m *BinaryManager) SetVerbose(verbose bool) {
	m.verbose = verbose
}

// CacheDir 获取当前使用的缓存目录路径
func (m *BinaryManager) CacheDir() string {
	return m.cacheDir
}

// GetPath 获取工具二进制文件的路径
//
// 参数:
//   - toolName: 工具名称
//
// 返回:
//   - string: 工具二进制文件的绝对路径
//   - error: 获取过程中可能的错误
//
// 查找优先级:
//  1. 检查缓存目录中是否有已存在的有效文件
//  2. 从嵌入的资源中读取并释放到缓存目录
//  3. 从本地assets目录复制到缓存目录
//  4. 最后尝试从系统PATH中查找
func (m *BinaryManager) GetPath(toolName string) (string, error) {
	cachedPath := filepath.Join(m.cacheDir, toolName)
	if m.isValid(cachedPath) {
		if m.verbose {
			fmt.Fprintf(os.Stderr, "[opskit] 使用缓存的 %s: %s\n", toolName, cachedPath)
		}
		return cachedPath, nil
	}

	// 优先从嵌入资源读取
	data, err := m.readEmbeddedAsset(toolName)
	if err == nil {
		if m.verbose {
			fmt.Fprintf(os.Stderr, "[opskit] 提取嵌入的 %s 到缓存\n", toolName)
		}
		if err := os.WriteFile(cachedPath, data, 0755); err != nil {
			return "", fmt.Errorf("写入嵌入资源到缓存失败: %w", err)
		}
		return cachedPath, nil
	}

	// 从assets目录复制
	srcPath := filepath.Join(m.assetDir, toolName)
	if _, err := os.Stat(srcPath); err == nil {
		if m.verbose {
			fmt.Fprintf(os.Stderr, "[opskit] 从assets目录复制 %s 到缓存\n", toolName)
		}
		data, err := os.ReadFile(srcPath)
		if err != nil {
			return "", fmt.Errorf("读取资源文件失败: %w", err)
		}
		if err := os.WriteFile(cachedPath, data, 0755); err != nil {
			return "", fmt.Errorf("写入到缓存失败: %w", err)
		}
		return cachedPath, nil
	}

	// fallback 到系统 PATH
	if path, err := execLookPath(toolName); err == nil {
		if m.verbose {
			fmt.Fprintf(os.Stderr, "[opskit] 使用系统的 %s: %s\n", toolName, path)
		}
		return path, nil
	}

	return "", fmt.Errorf("工具 %s 未找到（未嵌入、不在assets目录、也不在系统PATH中）", toolName)
}

// ListTools 列出所有内置工具的信息
//
// 返回:
//   - []ToolInfo: 工具信息列表
//   - error: 可能的错误
func (m *BinaryManager) ListTools() ([]ToolInfo, error) {
	return []ToolInfo{
		{
			Name:        "jq",
			Version:     "1.7.1",
			Description: "JSON 处理工具",
		},
		{
			Name:        "curl",
			Version:     "8.6.0",
			Description: "HTTP 客户端工具",
		},
		{
			Name:        "yq",
			Version:     "4.40.5",
			Description: "YAML/JSON 处理工具",
		},
		{
			Name:        "busybox",
			Version:     "1.36.1",
			Description: "BusyBox 工具箱（提供 telnet/ping/nslookup/netstat/nc/wget）",
			Provides:    []string{"telnet", "ping", "nslookup", "netstat", "nc", "wget"},
		},
	}, nil
}

// ExtractTo 将工具二进制文件提取到指定目录
//
// 参数:
//   - toolName: 工具名称
//   - destDir: 目标目录
//
// 返回:
//   - string: 提取后的文件路径
//   - error: 提取过程中可能的错误
func (m *BinaryManager) ExtractTo(toolName, destDir string) (string, error) {
	srcPath, err := m.GetPath(toolName)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("创建目标目录失败: %w", err)
	}

	data, err := os.ReadFile(srcPath)
	if err != nil {
		return "", fmt.Errorf("读取源文件失败: %w", err)
	}

	destPath := filepath.Join(destDir, toolName)
	if err := os.WriteFile(destPath, data, 0755); err != nil {
		return "", fmt.Errorf("写入目标文件失败: %w", err)
	}

	return destPath, nil
}

// Clean 清理工具缓存目录
func (m *BinaryManager) Clean() error {
	return os.RemoveAll(m.cacheDir)
}

// isValid 检查工具文件是否有效（存在且有执行权限）
func (m *BinaryManager) isValid(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode()&0111 != 0
}

// findAssetDir 找到二进制资源文件的存放目录
func findAssetDir() (string, error) {
	platform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)

	paths := []string{
		filepath.Join("assets", platform),
		filepath.Join("../assets", platform),
		filepath.Join("../../assets", platform),
	}

	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		paths = append(paths, filepath.Join(exeDir, "assets", platform))
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return filepath.Join("assets", platform), nil
}

// getDefaultCacheDir 获取默认缓存目录
//
// 优先级:
//  1. $OPSKIT_BIN_DIR/<version>/
//  2. $XDG_CACHE_HOME/opskit/<version>/
//  3. $HOME/.cache/opskit/<version>/
//  4. /tmp/.opskit-bin-<version>/
func getDefaultCacheDir() (string, error) {
	if dir := os.Getenv("OPSKIT_BIN_DIR"); dir != "" {
		return filepath.Join(dir, Version), nil
	}

	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return filepath.Join(os.TempDir(), ".opskit-bin-"+Version), nil
		}
		cacheHome = filepath.Join(home, ".cache")
	}

	return filepath.Join(cacheHome, "opskit", Version), nil
}

var execLookPath = func(name string) (string, error) {
	return exec.LookPath(name)
}
