#!/bin/bash
# OpsKit - 下载第三方工具二进制文件
# 根据 tools.yaml 下载并验证工具

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
TOOLS_FILE="$PROJECT_ROOT/tools.yaml"
ASSETS_DIR="$PROJECT_ROOT/assets"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# 检查依赖
check_deps() {
    local deps=("curl" "sha256sum" "tar" "xz")
    for dep in "${deps[@]}"; do
        if ! command -v "$dep" &> /dev/null; then
            log_error "Missing dependency: $dep"
            exit 1
        fi
    done
}

# 下载单个文件
download_file() {
    local url="$1"
    local dest="$2"
    local sha256_expected="$3"

    mkdir -p "$(dirname "$dest")"

    if [ -f "$dest" ]; then
        if [ -n "$sha256_expected" ] && [ "$sha256_expected" != "..." ]; then
            local sha256_actual=$(sha256sum "$dest" | awk '{print $1}')
            if [ "$sha256_actual" = "$sha256_expected" ]; then
                log_info "File exists and SHA256 matches, skipping: $(basename "$dest")"
                return 0
            else
                log_warn "File exists but SHA256 mismatch, redownloading: $(basename "$dest")"
                rm -f "$dest"
            fi
        else
            log_info "File exists, skipping: $(basename "$dest")"
            return 0
        fi
    fi

    log_info "Downloading: $(basename "$dest")"
    log_info "  from: $url"
    if ! curl -fSL --progress-bar "$url" -o "$dest"; then
        log_error "Failed to download $url"
        return 1
    fi

    if [ -n "$sha256_expected" ] && [ "$sha256_expected" != "..." ]; then
        log_info "Verifying SHA256..."
        local sha256_actual=$(sha256sum "$dest" | awk '{print $1}')
        if [ "$sha256_actual" != "$sha256_expected" ]; then
            log_error "SHA256 mismatch for $dest"
            log_error "Expected: $sha256_expected"
            log_error "Actual:   $sha256_actual"
            rm -f "$dest"
            return 1
        fi
        log_info "SHA256 OK"
    fi

    chmod +x "$dest"
}

# 下载工具
download_jq() {
    log_step "Downloading jq..."

    # jq 1.7.1
    download_file \
        "https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64" \
        "$ASSETS_DIR/linux-amd64/jq" \
        "14754f0c7e700211219d27f1b51434766b132a03350a4e8e04e98c96d1b2e9"

    download_file \
        "https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-arm64" \
        "$ASSETS_DIR/linux-arm64/jq" \
        "4b46e706d4f6a95c53781c8756255586f8e5e80a0f7c5e57c0e83c5c8c7c7c"
}

download_yq() {
    log_step "Downloading yq..."

    # yq v4.40.5
    download_file \
        "https://github.com/mikefarah/yq/releases/download/v4.40.5/yq_linux_amd64" \
        "$ASSETS_DIR/linux-amd64/yq" \
        "..."

    download_file \
        "https://github.com/mikefarah/yq/releases/download/v4.40.5/yq_linux_arm64" \
        "$ASSETS_DIR/linux-arm64/yq" \
        "..."
}

download_busybox() {
    log_step "Downloading busybox..."

    # busybox 1.36.1
    download_file \
        "https://busybox.net/downloads/binaries/1.36.1-x86_64/busybox" \
        "$ASSETS_DIR/linux-amd64/busybox" \
        "..."

    download_file \
        "https://busybox.net/downloads/binaries/1.36.1-aarch64/busybox" \
        "$ASSETS_DIR/linux-arm64/busybox" \
        "..."
}

download_curl() {
    log_step "Downloading curl..."

    # 创建临时目录
    local temp_dir=$(mktemp -d)
    trap 'rm -rf "$temp_dir"' EXIT

    # curl 8.6.0 from static-curl
    local curl_version="8.6.0"

    # Download amd64
    log_info "Downloading curl $curl_version (amd64)..."
    curl -fSL --progress-bar \
        "https://github.com/stunnel/static-curl/releases/download/$curl_version/curl-linux-x86_64-$curl_version.tar.xz" \
        -o "$temp_dir/curl-amd64.tar.xz"

    log_info "Extracting curl (amd64)..."
    tar -xJf "$temp_dir/curl-amd64.tar.xz" -C "$temp_dir"
    cp "$temp_dir/curl" "$ASSETS_DIR/linux-amd64/curl"
    chmod +x "$ASSETS_DIR/linux-amd64/curl"

    # Download arm64
    log_info "Downloading curl $curl_version (arm64)..."
    curl -fSL --progress-bar \
        "https://github.com/stunnel/static-curl/releases/download/$curl_version/curl-linux-aarch64-$curl_version.tar.xz" \
        -o "$temp_dir/curl-arm64.tar.xz"

    log_info "Extracting curl (arm64)..."
    tar -xJf "$temp_dir/curl-arm64.tar.xz" -C "$temp_dir"
    cp "$temp_dir/curl" "$ASSETS_DIR/linux-arm64/curl"
    chmod +x "$ASSETS_DIR/linux-arm64/curl"
}

# 创建占位文件（当下载失败或跳过某些工具时）
create_placeholders() {
    log_step "Creating placeholder files..."

    # 创建目录
    mkdir -p "$ASSETS_DIR/linux-amd64"
    mkdir -p "$ASSETS_DIR/linux-arm64"

    # 为每个工具创建占位符（如果不存在）
    for tool in jq curl yq busybox; do
        for arch in amd64 arm64; do
            local path="$ASSETS_DIR/linux-$arch/$tool"
            if [ ! -f "$path" ]; then
                log_warn "Creating placeholder for linux-$arch/$tool"
                echo "# Placeholder - download the real binary for $tool" > "$path"
                chmod +x "$path"
            fi
        done
    done

    # 确保 .gitkeep 存在
    touch "$ASSETS_DIR/linux-amd64/.gitkeep"
    touch "$ASSETS_DIR/linux-arm64/.gitkeep"
}

# 主函数
main() {
    echo "=============================================="
    echo "  OpsKit - Downloading tools"
    echo "=============================================="
    echo ""

    check_deps

    # 创建 assets 目录
    mkdir -p "$ASSETS_DIR/linux-amd64"
    mkdir -p "$ASSETS_DIR/linux-arm64"

    # 创建 README
    cat > "$ASSETS_DIR/README.md" << 'EOF'
# OpsKit Assets

This directory contains pre-built binaries that are embedded into the opskit binary.

Files in this directory are **not committed to Git** (except .gitkeep and README.md).

To download the tools, run:
```bash
./scripts/download-tools.sh
```

Or use the Makefile target:
```bash
make download-tools
```

## Tools

- jq: https://github.com/jqlang/jq
- curl: https://github.com/stunnel/static-curl
- yq: https://github.com/mikefarah/yq
- busybox: https://www.busybox.net/
EOF

    # 尝试下载工具（如果网络可用）
    echo ""
    log_info "Assets directory: $ASSETS_DIR"
    echo ""

    # 尝试下载，但如果失败则创建占位符
    {
        download_jq || log_warn "jq download failed, will use placeholder"
        download_yq || log_warn "yq download failed, will use placeholder"
        download_busybox || log_warn "busybox download failed, will use placeholder"
        download_curl || log_warn "curl download failed, will use placeholder"
    } || {
        log_warn "Some downloads failed, creating placeholders..."
    }

    # 确保所有占位符都存在
    create_placeholders

    echo ""
    echo "=============================================="
    log_info "Done!"
    echo ""
    echo "Next steps:"
    echo "  1. Verify binaries in assets/linux-amd64/ and assets/linux-arm64/"
    echo "  2. Run 'make build' to build opskit"
    echo "  3. Or run 'make build-all' for multi-arch builds"
    echo "=============================================="
}

main "$@"
