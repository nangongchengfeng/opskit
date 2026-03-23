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
