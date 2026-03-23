# OpsKit - 嵌入式运维工具箱

OpsKit 是一款单一二进制、零依赖、内置常用运维工具的 CLI 工具箱，专为受限环境下的故障排查设计。

## 项目介绍

在客户集群环境中，经常会遇到缺少基本运维工具（jq、curl、telnet 等）、环境没有网络连接无法在线安装工具等问题。OpsKit 完美解决了这些问题：

- **单文件分发** - 只需一个 `opskit` 二进制文件，scp 过去即可使用
- **离线使用** - 无需网络连接，所有工具已内置
- **版本锁定** - 统一工具版本，跨环境行为一致
- **权限友好** - 静态编译，可直接执行，无需安装权限
- **按需释放** - 运行时才释放工具到临时目录，不污染系统

## 功能特性

### 核心功能

- **工具代理调用** - `opskit <工具名> [参数...]` 完整透传所有参数给内嵌工具
  - stdout/stderr/退出码全部透传，与直接调用原生工具行为完全相同
  - 支持管道操作
  - 支持信号透传（Ctrl+C 正常终止）

### 工具管理命令

| 命令 | 说明 |
|------|------|
| `opskit list` | 列出所有内置工具及其版本、架构 |
| `opskit version` | 输出 opskit 自身版本及构建信息 |
| `opskit extract <工具名>` | 将指定工具释放到当前目录 |
| `opskit extract --all` | 释放所有工具到当前目录 |
| `opskit extract --dir <路径>` | 指定释放目录 |
| `opskit clean` | 清理临时释放目录 |
| `opskit which <工具名>` | 显示工具的临时释放路径 |

### 全局选项

| 选项 | 说明 |
|------|------|
| `--verbose, -v` | 输出 opskit 内部调试信息 |
| `--bin-dir <路径>` | 覆盖临时释放目录 |
| `--version` | 输出版本信息 |

## 内置工具

### 当前包含工具

| 工具 | 版本 | 描述 |
|------|------|------|
| jq | 1.7.1 | JSON 处理与过滤工具 |
| curl | 8.6.0 | HTTP 客户端工具 |
| yq | 4.40.5 | YAML/JSON 处理工具 |
| busybox | 1.36.1 | BusyBox 工具箱 |

### BusyBox 提供的工具

通过 BusyBox 可以使用以下工具：

- `telnet` - TCP 端口连通测试
- `ping` - ICMP/主机可达性测试
- `nslookup` - DNS 查询
- `netstat` - 端口与连接状态查看
- `nc` - netcat 端口监听/数据传输
- `wget` - 文件下载

## 编译步骤

### 前置要求

- Go 1.18+ (推荐 1.22+)
- 网络连接（用于下载工具二进制）
- Linux 环境（用于下载对应平台的工具）

### 编译步骤

1. **克隆仓库**

```bash
git clone https://github.com/opskit/opskit.git
cd opskit
```

2. **下载第三方工具二进制文件**

```bash
make download-tools
```

此脚本会根据 `tools.yaml` 中的配置下载对应平台的工具二进制文件到 `assets/` 目录。

3. **编译**

```bash
# 编译当前平台版本
make build

# 或直接使用 go build（需要设置 build tag）
CGO_ENABLED=0 go build -tags=embed_assets -o bin/opskit ./cmd/opskit
```

4. **编译所有平台**

```bash
# 编译 Linux amd64 和 arm64 版本
make build-all
```

5. **运行测试**

```bash
make test
```

### 使用 goreleaser 发布

```bash
# 本地快照构建
make build-snapshot

# 发布版本（需要 GITHUB_TOKEN）
make release
```

## 使用示例

### 运行内置工具

```bash
# JSON 处理
echo '{"name": "OpsKit", "version": "1.0.0"}' | opskit jq '.name'

# HTTP 请求
opskit curl https://example.com

# YAML 处理
opskit yq eval '.spec.replicas' deployment.yaml

# TCP 连通测试
opskit telnet db-service 5432

# 端口扫描
opskit nc -zv example.com 80 443

# Ping 测试
opskit ping -c 3 example.com

# DNS 查询
opskit nslookup example.com

# 文件下载
opskit wget https://example.com/file.txt
```

### 管道操作

```bash
# 分析 Kubernetes 日志
kubectl logs pod-xxx | opskit jq '.message'

# 处理 JSON API 响应
opskit curl https://api.example.com/data | opskit jq '.items[] | {name: .name}'
```

### 管理命令

```bash
# 列出所有内置工具
opskit list

# 显示版本信息
opskit version

# 提取工具到本地目录
opskit extract jq
opskit extract --dir /usr/local/bin jq

# 查看工具缓存路径
opskit which jq

# 清理缓存
opskit clean
```

### 详细模式

```bash
# 查看工具加载过程
opskit -v jq --help
```

## 项目结构

```
opskit/
├── cmd/
│   └── opskit/              # 主入口
│       └── main.go
├── internal/
│   ├── cli/                 # CLI 命令定义
│   │   ├── root.go          # 根命令
│   │   ├── list.go          # list 命令
│   │   ├── version.go       # version 命令
│   │   ├── extract.go       # extract 命令
│   │   ├── clean.go         # clean 命令
│   │   └── which.go         # which 命令
│   └── embed/               # 二进制管理与执行
│       ├── manager.go       # 二进制管理器
│       ├── executor.go      # 工具执行器
│       ├── manager_embedded.go  # 嵌入资源实现（带 build tag）
│       └── manager_noembed.go   # 非嵌入资源实现（带 build tag）
├── assets/                  # 第三方二进制（构建时下载）
│   ├── linux-amd64/
│   └── linux-arm64/
├── scripts/                 # 辅助脚本
│   └── download-tools.sh    # 下载工具脚本
├── pkg/                     # 可导出的公共包（预留）
├── tools.yaml               # 工具版本清单
├── Makefile                 # 构建脚本
├── .goreleaser.yaml         # goreleaser 配置
├── go.mod
├── go.sum
└── README.md
```

## 工作原理

### 构建阶段

利用 Go 1.16+ 的 `//go:embed` 指令，在编译阶段将各平台的原生工具二进制文件作为静态资源打包进 opskit 可执行文件。

```
jq_linux_amd64   ──┐
curl_linux_amd64  ─┤  go:embed  ──►  opskit（单一二进制）
busybox_amd64     ──┘
```

### 运行阶段

当调用 `opskit <工具名> [参数...]` 时：

1. 检测当前 GOOS/GOARCH
2. 检查工具是否已在缓存目录中且 SHA256 一致
3. 未命中缓存则从嵌入资源中释放对应二进制到缓存目录
4. 释放后设置执行权限，使用 exec.Command 调用，参数完全透传
5. stdout/stderr/退出码直接透传给用户

### 缓存目录策略

释放目录优先级：
1. `$OPSKIT_BIN_DIR/<version>/`（用户自定义）
2. `$HOME/.cache/opskit/<version>/`（推荐）
3. `/tmp/.opskit-bin-<version>/`（降级方案）

## 开发指南

### 添加新工具

1. 在 `tools.yaml` 中添加工具配置
2. 更新 `internal/embed/manager.go` 中的 `ListTools()` 方法
3. 运行 `make download-tools` 下载二进制
4. 测试新工具是否正常工作

### 代码规范

- 所有公开方法和关键逻辑必须添加清晰的中文注释
- 使用 `gofmt` 格式化代码
- 遵循 Go 官方代码规范
- 提交前运行 `make lint` 和 `make vet`

## 兼容性

### 支持的平台

| 平台 | 架构 | 最低版本 |
|------|------|----------|
| Linux | amd64 | CentOS 7+, Ubuntu 18.04+ |
| Linux | arm64 | CentOS 7+, Ubuntu 18.04+ |

### 安全特性

- 构建时验证所有工具的 SHA256，防止供应链污染
- 释放前校验 SHA256，防止临时目录中的文件被篡改
- tools.yaml 中明确记录工具来源与哈希，可审计

## 常见问题

### Q: 为什么不直接使用系统已有的工具？

A: OpsKit 的设计目标是在受限环境中使用，这些环境通常没有所需的工具，或者工具版本不兼容。

### Q: 如何减小二进制文件的大小？

A: 可以使用 UPX 压缩：`upx --best bin/opskit`。

### Q: 缓存目录会占用太多空间吗？

A: 不会，只有使用过的工具才会被释放到缓存目录。可以随时使用 `opskit clean` 清理。

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

## 致谢

感谢所有为这个项目做出贡献的开发者。

---

**OpsKit** - 一个二进制，走遍所有客户集群。
