# Changelog

## v1.4.2-rc.1 (2026-11-13)

- 全局标点符号统一为英文半角
- 集群DNS支持所有网站共用一个CNAME, 开启后不再生成网站独立CNAME记录
- 允许通过CNAME访问网站默认关闭, 共用CNAME时强制关闭
- 移除节点数量上限和授权配额检查
- 移除平台用户登录的商业版限制
- 健康检查自动上下线移除商业版分支, 统一使用节点IP级状态
- DNS服务商描述移除商业版措辞
- 升级全部依赖链并升级 Go 工具链到 1.27.4, govulncheck 符号级扫描可达漏洞 22 -> 0
- google.golang.org/grpc v1.63.2 -> v1.85.0, 修复 :path 缺少前导斜杠导致的授权绕过 (GO-2026-4762) 及 HTTP/2 DATA 帧堆内存耗尽 (GO-2026-6348)
- github.com/go-acme/lego/v4 v4.17.3 -> v4.35.2, 修复 ACME API 不强制 HTTPS (GO-2025-3847)
- github.com/go-jose/go-jose/v4 v4.0.1 -> v4.1.9, 修复 JWE 解密 panic (GO-2026-4945)
- golang.org/x/crypto v0.23.0 -> v0.59.0, golang.org/x/net v0.25.0 -> v0.62.0, golang.org/x/text v0.15.0 -> v0.44.0, andybalholm/brotli v1.0.5 -> v1.3.0
- 新增 internal/dnsclients/dnsla 桩包: 该包在开源版中被裁剪但 provider_dns_la.go 仍有引用, 仅为通过编译, DNS.LA 功能不可用
- 将模块路径, 源码导入和文档链接中的仓库地址统一迁移到 git.ovh.im/yeqing
- cmd/edge-api 启动时读取 EDGE_DB_HOST, EDGE_DB_PORT, EDGE_DB_USER, EDGE_DB_PASSWORD, EDGE_DB_NAME, 自动生成 configs/db.yaml 和 configs/.db.yaml
- 管理员创建时自动创建同用户名的关联普通用户, 老库升级为每个管理员补建关联用户
- 管理员修改用户名/密码/全名时同步到关联用户
- 管理员创建网站默认绑定其关联用户, 未指定时自动关联
- 初始化不再创建随机USER_xxx用户, 改为按管理员创建关联用户
- 移除用户手机号/联系电话字段及手机验证码相关模型和服务
- edgeUsers 去除 mobile/tel/verifiedMobile/mobileIsVerified 列, edgeAdmins 新增 userId 列, 老库自动升级
