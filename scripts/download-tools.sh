#!/bin/bash
# OpsKit - 下载第三方工具二进制文件
# 根据 tools.yaml 下载并验证工具

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
TOOLS_FILE="$PROJECT_ROOT/tools.yaml"
ASSETS_DIR="$PROJECT_ROOT/assets"

# 全局临时目录，避免unbound variable错误
TEMP_DIR=""
trap 'if [ -n "$TEMP_DIR" ] && [ -d "$TEMP_DIR" ]; then rm -rf "$TEMP_DIR"; fi' EXIT

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 国内镜像加速（如果需要可以取消注释）
# GITHUB_MIRROR="https://mirror.ghproxy.com"
GITHUB_MIRROR=""

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
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

    # 使用国内镜像加速
    if [[ "$url" == https://github.com/* ]]; then
        url="${GITHUB_MIRROR}/${url}"
    fi

    log_info "Downloading: $(basename "$dest")"
    log_info "  from: $url"
    if ! curl -fSL --progress-bar "$url" -o "$dest"; then
        log_error "Failed to download $url"
    fi

    if [ -n "$sha256_expected" ] && [ "$sha256_expected" != "..." ]; then
        log_info "Verifying SHA256..."
        local sha256_actual=$(sha256sum "$dest" | awk '{print $1}')
        if [ "$sha256_actual" != "$sha256_expected" ]; then
            log_error "SHA256 mismatch for $dest\nExpected: $sha256_expected\nActual:   $sha256_actual"
        fi
        log_info "SHA256 OK"
    fi

    chmod +x "$dest"
}

# 下载工具
download_jq() {
    log_step "Downloading jq..."

    # jq 1.7.1 - 正确的SHA256
    download_file \
        "https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64" \
        "$ASSETS_DIR/linux-amd64/jq" \
        "5942c9b0934e510ee61eb3e30273f1b3fe2590df93933a93d7c58b81d19c8ff5"

    download_file \
        "https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-arm64" \
        "$ASSETS_DIR/linux-arm64/jq" \
        "4dd2d8a0661df0b22f1bb9a1f9830f06b6f3b8f7d91211a1ef5d7c4f06a8b4a5"
}

download_yq() {
    log_step "Downloading yq..."

    # yq v4.40.5 - 添加正确的SHA256
    download_file \
        "https://github.com/mikefarah/yq/releases/download/v4.40.5/yq_linux_amd64" \
        "$ASSETS_DIR/linux-amd64/yq" \
        "80c5b96695631b5196447249854987ff1a07c5e875b760f1c7b104b570d173f5"

    download_file \
        "https://github.com/mikefarah/yq/releases/download/v4.40.5/yq_linux_arm64" \
        "$ASSETS_DIR/linux-arm64/yq" \
        "18602f890b4f54773c56512384aa8e324d8d74c5d48978c5653e8e58e0386111"
}

download_busybox() {
    log_step "Downloading busybox..."

    # busybox 1.36.1 - 官方源
    download_file \
        "https://busybox.net/downloads/binaries/1.36.1-x86_64/busybox" \
        "$ASSETS_DIR/linux-amd64/busybox" \
        "523800dd278f988a73933977f8e1779e503bd59578e55f8d4d77f93b85c62470"

    download_file \
        "https://busybox.net/downloads/binaries/1.36.1-aarch64/busybox" \
        "$ASSETS_DIR/linux-arm64/busybox" \
        "f2a527c335955c73656e16b6f505a343935c29d66f6a27f75f4a079e9c551333"
}

download_curl() {
    log_step "Downloading curl..."

    # 使用全局临时目录
    TEMP_DIR=$(mktemp -d)

    # curl 8.8.0 from static-curl
    local curl_version="8.8.0"

    # Download amd64
    log_info "Downloading curl $curl_version (amd64)..."
    curl -fSL --progress-bar \
        "${GITHUB_MIRROR}/https://github.com/stunnel/static-curl/releases/download/$curl_version/curl-linux-x86_64-$curl_version.tar.xz" \
        -o "$TEMP_DIR/curl-amd64.tar.xz"

    log_info "Extracting curl (amd64)..."
    tar -xJf "$TEMP_DIR/curl-amd64.tar.xz" -C "$TEMP_DIR"
    cp "$TEMP_DIR/curl" "$ASSETS_DIR/linux-amd64/curl"
    chmod +x "$ASSETS_DIR/linux-amd64/curl"

    # Download arm64
    log_info "Downloading curl $curl_version (arm64)..."
    curl -fSL --progress-bar \
        "${GITHUB_MIRROR}/https://github.com/stunnel/static-curl/releases/download/$curl_version/curl-linux-aarch64-$curl_version.tar.xz" \
        -o "$TEMP_DIR/curl-arm64.tar.xz"

    log_info "Extracting curl (arm64)..."
    tar -xJf "$TEMP_DIR/curl-arm64.tar.xz" -C "$TEMP_DIR"
    cp "$TEMP_DIR/curl" "$ASSETS_DIR/linux-arm64/curl"
    chmod +x "$ASSETS_DIR/linux-arm64/curl"
}

# 验证所有文件都存在
validate_downloads() {
    log_step "Validating downloads..."
    local required_files=(
        "$ASSETS_DIR/linux-amd64/jq"
        "$ASSETS_DIR/linux-amd64/yq"
        "$ASSETS_DIR/linux-amd64/busybox"
        "$ASSETS_DIR/linux-amd64/curl"
        "$ASSETS_DIR/linux-arm64/jq"
        "$ASSETS_DIR/linux-arm64/yq"
        "$ASSETS_DIR/linux-arm64/busybox"
        "$ASSETS_DIR/linux-arm64/curl"
    )

    for file in "${required_files[@]}"; do
        if [ ! -f "$file" ] || [ ! -x "$file" ]; then
            log_error "Missing or non-executable file: $file"
        fi
        # 检查不是占位符
        if head -n1 "$file" | grep -q "Placeholder"; then
            log_error "Placeholder file found: $file, download failed"
        fi
        log_info "✓ $(basename "$file"): OK"
    done
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

    echo ""
    log_info "Assets directory: $ASSETS_DIR"
    log_info "Using GitHub mirror: $GITHUB_MIRROR"
    echo ""

    # 下载所有工具，失败直接报错
    download_jq
    download_yq
    download_busybox
    download_curl

    # 验证所有下载
    validate_downloads

    # 确保 .gitkeep 存在
    touch "$ASSETS_DIR/linux-amd64/.gitkeep"
    touch "$ASSETS_DIR/linux-arm64/.gitkeep"

    echo ""
    echo "=============================================="
    log_info "All tools downloaded successfully!"
    echo ""
    echo "Next steps:"
    echo "  1. Run 'make build' to build opskit with embedded assets"
    echo "  2. Or run 'make build-all' for multi-arch builds"
    echo "=============================================="
}

main "$@"
