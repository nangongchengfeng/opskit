// Package embed manages embedded binary assets
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
	Version = "0.1.0"
)

type BinaryManager struct {
	cacheDir string
	verbose  bool
	assetDir string // 存放二进制文件的目录
}

type ToolInfo struct {
	Name        string
	Version     string
	Description string
	Provides    []string
}

func NewManager(cacheDir string, verbose bool) (*BinaryManager, error) {
	if cacheDir == "" {
		var err error
		cacheDir, err = getDefaultCacheDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get default cache dir: %w", err)
		}
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache dir: %w", err)
	}

	// 找到 assets 目录
	assetDir, err := findAssetDir()
	if err != nil {
		return nil, fmt.Errorf("failed to find asset dir: %w", err)
	}

	return &BinaryManager{
		cacheDir: cacheDir,
		verbose:  verbose,
		assetDir: assetDir,
	}, nil
}

func GetManager() (*BinaryManager, error) {
	var err error
	once.Do(func() {
		globalManager, err = NewManager("", false)
	})
	return globalManager, err
}

func (m *BinaryManager) SetVerbose(verbose bool) {
	m.verbose = verbose
}

func (m *BinaryManager) CacheDir() string {
	return m.cacheDir
}

func (m *BinaryManager) GetPath(toolName string) (string, error) {
	cachedPath := filepath.Join(m.cacheDir, toolName)
	if m.isValid(cachedPath) {
		if m.verbose {
			fmt.Fprintf(os.Stderr, "[opskit] Using cached %s at %s\n", toolName, cachedPath)
		}
		return cachedPath, nil
	}

	// 优先从嵌入资源读取
	data, err := m.readEmbeddedAsset(toolName)
	if err == nil {
		if m.verbose {
			fmt.Fprintf(os.Stderr, "[opskit] Extracting embedded %s to cache\n", toolName)
		}
		if err := os.WriteFile(cachedPath, data, 0755); err != nil {
			return "", fmt.Errorf("failed to write embedded asset to cache: %w", err)
		}
		return cachedPath, nil
	}

	// 从 assets 目录复制
	srcPath := filepath.Join(m.assetDir, toolName)
	if _, err := os.Stat(srcPath); err == nil {
		if m.verbose {
			fmt.Fprintf(os.Stderr, "[opskit] Copying %s from assets to cache\n", toolName)
		}
		data, err := os.ReadFile(srcPath)
		if err != nil {
			return "", fmt.Errorf("failed to read asset: %w", err)
		}
		if err := os.WriteFile(cachedPath, data, 0755); err != nil {
			return "", fmt.Errorf("failed to write to cache: %w", err)
		}
		return cachedPath, nil
	}

	// fallback 到系统 PATH
	if path, err := execLookPath(toolName); err == nil {
		if m.verbose {
			fmt.Fprintf(os.Stderr, "[opskit] Using system %s at %s\n", toolName, path)
		}
		return path, nil
	}

	return "", fmt.Errorf("tool %s not found (not embedded, not in assets/, not on PATH)", toolName)
}

func (m *BinaryManager) ListTools() ([]ToolInfo, error) {
	return []ToolInfo{
		{
			Name:        "jq",
			Version:     "1.7.1",
			Description: "JSON processing tool",
		},
		{
			Name:        "curl",
			Version:     "8.6.0",
			Description: "HTTP client",
		},
		{
			Name:        "yq",
			Version:     "4.40.5",
			Description: "YAML/JSON processor",
		},
		{
			Name:        "busybox",
			Version:     "1.36.1",
			Description: "Toolbox providing multiple utilities",
			Provides:    []string{"telnet", "ping", "nslookup", "netstat", "nc", "wget"},
		},
	}, nil
}

func (m *BinaryManager) ExtractTo(toolName, destDir string) (string, error) {
	srcPath, err := m.GetPath(toolName)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create dest dir: %w", err)
	}

	data, err := os.ReadFile(srcPath)
	if err != nil {
		return "", fmt.Errorf("failed to read source: %w", err)
	}

	destPath := filepath.Join(destDir, toolName)
	if err := os.WriteFile(destPath, data, 0755); err != nil {
		return "", fmt.Errorf("failed to write %s: %w", destPath, err)
	}

	return destPath, nil
}

func (m *BinaryManager) Clean() error {
	return os.RemoveAll(m.cacheDir)
}

func (m *BinaryManager) isValid(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode()&0111 != 0
}

// findAssetDir 找到 assets 目录
func findAssetDir() (string, error) {
	platform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)

	// 尝试几种可能的路径
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

	// 如果都找不到，返回相对于当前工作目录的路径
	return filepath.Join("assets", platform), nil
}

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
