# Linux 打包测试指南

本指南说明如何在 Linux 上下载真实的二进制文件并打包 opskit。

## 快速开始

### 1. 准备环境

```bash
# 克隆或复制项目到 Linux
cd opskit

# 确保有执行权限
chmod +x scripts/download-tools.sh
chmod +x scripts/*.sh
```

### 2. 下载工具二进制

```bash
# 方式一：使用 Makefile
make download-tools

# 方式二：直接运行脚本
./scripts/download-tools.sh
```

### 3. 构建（带嵌入资产）

```bash
# 先确认 assets 目录有文件
ls -la assets/linux-amd64/
ls -la assets/linux-arm64/

# 方式一：使用 Makefile（需要修改以支持 embed_assets tag）
# 编辑 Makefile，在 build 目标添加 -tags=embed_assets

# 方式二：直接使用 go build
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -tags=embed_assets -v \
  -ldflags="-s -w -X main.Version=v0.1.0" \
  -o bin/opskit-linux-amd64 \
  ./cmd/opskit
```

## 手动下载二进制（如果脚本失败）

如果自动下载脚本不工作，可以手动下载：

### 下载 jq

```bash
# linux-amd64
curl -fSL https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64 \
  -o assets/linux-amd64/jq
chmod +x assets/linux-amd64/jq

# linux-arm64
curl -fSL https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-arm64 \
  -o assets/linux-arm64/jq
chmod +x assets/linux-arm64/jq
```

### 下载 yq

```bash
# linux-amd64
curl -fSL https://github.com/mikefarah/yq/releases/download/v4.40.5/yq_linux_amd64 \
  -o assets/linux-amd64/yq
chmod +x assets/linux-amd64/yq

# linux-arm64
curl -fSL https://github.com/mikefarah/yq/releases/download/v4.40.5/yq_linux_arm64 \
  -o assets/linux-arm64/yq
chmod +x assets/linux-arm64/yq
```

### 下载 busybox

```bash
# linux-amd64
curl -fSL https://busybox.net/downloads/binaries/1.36.1-x86_64/busybox \
  -o assets/linux-amd64/busybox
chmod +x assets/linux-amd64/busybox

# linux-arm64
curl -fSL https://busybox.net/downloads/binaries/1.36.1-aarch64/busybox \
  -o assets/linux-arm64/busybox
chmod +x assets/linux-arm64/busybox
```

### 下载 curl (static)

```bash

# amd64
curl -fSL https://github.com/stunnel/static-curl/releases/download/8.6.0-1/curl-linux-x86_64-8.6.0.tar.xz \
  -o /tmp/curl-amd64.tar.xz
tar -xJf /tmp/curl-amd64.tar.xz -C /tmp/
cp /tmp/curl internal/embed/assets/linux-amd64/curl
chmod +x internal/embed/assets/linux-amd64/curl


# arm64
curl -fSL https://github.com/stunnel/static-curl/releases/download/8.6.0-1/curl-linux-aarch64-8.6.0.tar.xz \
  -o /tmp/curl-arm64.tar.xz
tar -xJf /tmp/curl-arm64.tar.xz -C /tmp/
cp /tmp/curl internal/embed/assets/linux-arm64/curl
chmod +x internal/embed/assets/linux-arm64/curl
```

## 修改 Makefile 以支持嵌入构建

编辑 `Makefile`，更新 build 目标：

```makefile
.PHONY: build
build: download-tools ## Build opskit for current platform (with embedded assets)
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -tags=embed_assets -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/opskit

.PHONY: build-linux-amd64
build-linux-amd64: download-tools ## Build for Linux amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -tags=embed_assets -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/opskit

.PHONY: build-linux-arm64
build-linux-arm64: download-tools ## Build for Linux arm64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) -tags=embed_assets -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/opskit
```

## 验证构建是否成功

```bash
# 构建
make build-linux-amd64

# 测试
./bin/opskit-linux-amd64 list
./bin/opskit-linux-amd64 version

# 检查大小
ls -lh bin/
```

## 测试工具

```bash
# 如果系统有 jq，可以测试 fallback
./bin/opskit jq --help

# 测试管理命令
./bin/opskit list
./bin/opskit which jq
```

## 常见问题

### 1. "pattern assets/linux-amd64/*: no matching files found"

确保 assets 目录中有文件：
```bash
ls -la assets/linux-amd64/
# 应该看到: jq, curl, yq, busybox
```

### 2. 构建后工具找不到

确保使用了 `-tags=embed_assets`：
```bash
go build -tags=embed_assets ...
```

### 3. 工具执行权限问题

确保下载的二进制有执行权限：
```bash
chmod +x assets/linux-amd64/*
chmod +x assets/linux-arm64/*
```

## 使用 goreleaser (可选)

如果你有 goreleaser：

```bash
# 编辑 .goreleaser.yaml，在 builds 部分添加 tags
# builds:
#   - id: opskit
#     ...
#     tags:
#       - embed_assets

# 构建快照
goreleaser build --snapshot --clean
```
