# Changelog

## v1.4.2-rc.1 (2026-11-13)

- 升级全部依赖链并升级 Go 工具链到 1.27.4, govulncheck 符号级扫描可达漏洞 22 -> 0
- google.golang.org/grpc v1.63.2 -> v1.85.0, 修复 :path 缺少前导斜杠导致的授权绕过 (GO-2026-4762) 及 HTTP/2 DATA 帧堆内存耗尽 (GO-2026-6348)
- github.com/go-acme/lego/v4 v4.17.3 -> v4.35.2, 修复 ACME API 不强制 HTTPS (GO-2025-3847)
- github.com/go-jose/go-jose/v4 v4.0.1 -> v4.1.9, 修复 JWE 解密 panic (GO-2026-4945)
- golang.org/x/crypto v0.23.0 -> v0.59.0, golang.org/x/net v0.25.0 -> v0.62.0, golang.org/x/text v0.15.0 -> v0.44.0, andybalholm/brotli v1.0.5 -> v1.3.0
- 新增 internal/dnsclients/dnsla 桩包: 该包在开源版中被裁剪但 provider_dns_la.go 仍有引用, 仅为通过编译, DNS.LA 功能不可用
- 将模块路径, 源码导入和文档链接中的仓库地址统一迁移到 git.ovh.im/yeqing
- cmd/edge-api 启动时读取 EDGE_DB_HOST, EDGE_DB_PORT, EDGE_DB_USER, EDGE_DB_PASSWORD, EDGE_DB_NAME, 自动生成 configs/db.yaml 和 configs/.db.yaml
