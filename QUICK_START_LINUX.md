# OpsKit - Linux 快速启动指南

## 复制到 Linux

将整个项目文件夹复制到 Linux 机器上：

```bash
# 如果在 Windows，可以用 scp
scp -r opskit user@linux-machine:~/

# 或者使用其他方式复制
```

## 在 Linux 上的步骤

### 1. 进入项目目录

```bash
cd ~/opskit
```

### 2. 添加执行权限

```bash
chmod +x scripts/download-tools.sh
```

### 3. 下载真实的二进制文件

```bash
# 方式一：使用脚本
./scripts/download-tools.sh

# 方式二：手动下载（如果脚本失败）
# 参考 LINUX_BUILD_GUIDE.md 中的手动下载部分
```

### 4. 修改 Makefile（关键步骤！）

编辑 `Makefile`，找到 `build`、`build-linux-amd64`、`build-linux-arm64` 这些目标，
**添加 `-tags=embed_assets`** 参数：

```makefile
.PHONY: build
build: download-tools ## Build opskit for current platform
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -tags=embed_assets -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/opskit

.PHONY: build-linux-amd64
build-linux-amd64: download-tools ## Build for Linux amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -tags=embed_assets -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/opskit
```

### 5. 构建

```bash
# 构建当前平台（linux-amd64）
make build-linux-amd64

# 或者直接用 go build
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -tags=embed_assets -v \
  -ldflags="-s -w -X main.Version=v0.1.0" \
  -o bin/opskit ./cmd/opskit
```

### 6. 测试

```bash
# 查看版本
./bin/opskit version

# 列出工具
./bin/opskit list

# 测试工具（如果系统有 jq，会使用系统的作为 fallback）
./bin/opskit jq --help
```

## 验证 assets 目录

确保下载后有这些文件：

```
assets/
├── linux-amd64/
│   ├── jq
│   ├── curl
│   ├── yq
│   └── busybox
└── linux-arm64/
    ├── jq
    ├── curl
    ├── yq
    └── busybox
```

## 文件清单

复制到 Linux 的重要文件：

```
opskit/
├── go.mod
├── go.sum
├── Makefile
├── tools.yaml
├── .goreleaser.yaml
├── QUICK_START_LINUX.md      ← 本文档
├── LINUX_BUILD_GUIDE.md      ← 详细指南
├── cmd/
│   └── opskit/
│       └── main.go
├── internal/
│   ├── cli/
│   │   ├── root.go
│   │   ├── list.go
│   │   ├── version.go
│   │   ├── extract.go
│   │   ├── clean.go
│   │   └── which.go
│   └── embed/
│       ├── manager.go
│       ├── manager_embedded.go  ← 关键！
│       ├── executor.go
│       └── assets.go
├── scripts/
│   └── download-tools.sh
└── assets/
    ├── linux-amd64/
    │   └── .gitkeep
    └── linux-arm64/
        └── .gitkeep
```

## 常见问题

### Q: 构建时提示 "no matching files found"

A: 确保 assets/linux-amd64/ 目录下有 jq/curl/yq/busybox 这些文件，并且有执行权限。

### Q: 运行时工具找不到

A: 确认构建时使用了 `-tags=embed_assets` 参数。

### Q: 更详细的说明？

A: 查看 `LINUX_BUILD_GUIDE.md` 获取完整细节。

---

祝打包顺利！
