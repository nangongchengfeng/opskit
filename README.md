# OpsKit

OpsKit 是一款单一二进制、零依赖、内置常用运维工具的 CLI 工具箱，专为受限环境下的故障排查设计。

## 特性

- **单文件分发** - 只需一个 `opskit` 二进制文件
- **内置工具** - 包含 jq、curl、yq、busybox 等常用工具
- **离线使用** - 无需网络连接，开箱即用
- **版本锁定** - 统一工具版本，行为可预测
- **跨平台** - 支持 Linux (amd64/arm64)

## 快速开始

### 下载

从 [Releases](../../releases) 页面下载对应平台的二进制文件。

### 安装

```bash
# 下载后直接使用
chmod +x opskit
./opskit --help

# 或安装到系统路径
sudo mv opskit /usr/local/bin/
```

## 使用示例

### 运行内置工具

```bash
# JSON 处理
opskit jq '.name' data.json

# HTTP 请求
opskit curl https://example.com

# YAML 处理
opskit yq eval '.spec.replicas' deployment.yaml

# TCP 连通测试
opskit telnet db-service 5432

# 端口扫描
opskit nc -zv example.com 80 443
```

### 管理命令

```bash
# 列出所有内置工具
opskit list

# 显示版本信息
opskit version

# 提取工具到本地目录
opskit extract jq
opskit extract --all --dir /usr/local/bin/

# 查看工具缓存路径
opskit which jq

# 清理缓存
opskit clean
```

## 内置工具

| 工具 | 版本 | 描述 |
|------|------|------|
| jq | 1.7.1 | JSON 处理工具 |
| curl | 8.6.0 | HTTP 客户端 |
| yq | 4.40.5 | YAML/JSON 处理器 |
| busybox | 1.36.1 | 提供 telnet/ping/nslookup/netstat/nc/wget |

## 构建

### 前置要求

- Go 1.22+
- 网络连接（用于下载工具二进制）

### 构建步骤

```bash
# 克隆仓库
git clone https://github.com/opskit/opskit.git
cd opskit

# 下载第三方工具
make download-tools

# 构建
make build

# 运行测试
make test
```

### 发布构建

```bash
# 本地快照构建
make build-all

# 发布版本（需要 GITHUB_TOKEN）
make release
```

## 项目结构

```
opskit/
├── cmd/
│   └── opskit/          # 主入口
├── internal/
│   ├── cli/             # CLI 命令
│   └── embed/           # 二进制管理
├── assets/              # 第三方二进制（构建时下载）
├── scripts/             # 辅助脚本
├── tools.yaml           # 工具版本清单
└── Makefile
```

## 工作原理

OpsKit 使用 Go 1.16+ 的 `embed` 特性，将预编译的工具二进制打包到单个可执行文件中。运行时按需释放到缓存目录并执行。

```
构建时: jq/curl/yq/busybox ──go:embed──► opskit (单一二进制)

运行时: opskit jq ... ──► 释放到 ~/.cache/opskit/ ──► exec 执行
```

## 开发路线图

- [x] MVP 核心框架
- [x] jq/curl/yq/busybox 集成
- [x] 管理命令 (list/extract/clean/which)
- [ ] 实际二进制嵌入
- [ ] 更多工具 (grpcurl/tcpdump/openssl)
- [ ] Windows/macOS 支持
- [ ] UPX 压缩

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

---

**OpsKit** - 一个二进制，走遍所有客户集群。
