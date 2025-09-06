# 多提供商配置指南

GoDNS 现已支持同时配置多个 DNS 提供商，允许您从单个 GoDNS 实例管理不同 DNS 服务的域名。

## 功能特性

- **多提供商支持**：在一个配置文件中配置多个 DNS 提供商（Cloudflare、DNSPod、DigitalOcean 等）
- **每域名指定提供商**：为不同域名指定不同的提供商
- **完全向后兼容**：现有的单提供商配置无需修改即可继续使用
- **混合配置**：结合传统全局提供商与新的按域名提供商设置
- **提供商特定凭据**：每个提供商可以有自己的身份验证凭据

## 配置选项

### 选项 1：多提供商配置

配置多个提供商并按域名分配：

```json
{
  "providers": {
    "Cloudflare": {
      "email": "user@example.com",
      "password": "your-cloudflare-api-token"
    },
    "DNSPod": {
      "login_token": "your-dnspod-token"
    },
    "DigitalOcean": {
      "login_token": "your-digitalocean-api-token"
    }
  },
  "domains": [
    {
      "domain_name": "example.com",
      "sub_domains": ["www", "api"],
      "provider": "Cloudflare"
    },
    {
      "domain_name": "mysite.net",
      "sub_domains": ["mail", "ftp"], 
      "provider": "DNSPod"
    }
  ]
}
```

### 选项 2：传统单提供商（向后兼容）

现有配置无需修改即可继续使用：

```json
{
  "provider": "DNSPod",
  "login_token": "your-dnspod-token",
  "domains": [
    {
      "domain_name": "example.com",
      "sub_domains": ["www", "test"]
    }
  ]
}
```

### 选项 3：混合配置

结合全局提供商与特定的按域名提供商：

```json
{
  "provider": "DNSPod", 
  "login_token": "your-dnspod-token",
  "providers": {
    "Cloudflare": {
      "email": "user@example.com",
      "password": "your-cloudflare-api-token"
    }
  },
  "domains": [
    {
      "domain_name": "oldsite.com",
      "sub_domains": ["www", "mail"]
    },
    {
      "domain_name": "newsite.com", 
      "sub_domains": ["www", "api"],
      "provider": "Cloudflare"
    }
  ]
}
```

## 提供商配置字段

`providers` 部分中的每个提供商支持以下通用字段：

- `email`：用于身份验证的电子邮件地址（Cloudflare 等）
- `password`：API 令牌或密码
- `password_file`：包含密码/令牌的文件路径
- `login_token`：基于令牌的身份验证登录令牌（DNSPod 等）
- `login_token_file`：包含登录令牌的文件路径
- `app_key`：应用程序密钥（提供商特定）
- `app_secret`：应用程序密钥（提供商特定）
- `consumer_key`：消费者密钥（提供商特定）

## 域名配置

域名支持可选的 `provider` 字段：

```json
{
  "domain_name": "example.com",
  "sub_domains": ["www", "api", "@"],
  "provider": "Cloudflare"
}
```

如果省略 `provider`，域名将使用全局 `provider` 设置。

## 支持的提供商

多提供商模式支持所有现有提供商。请在配置中使用这些**精确**的提供商名称：

| 提供商名称 | 配置值 | 身份验证方法 |
|-----------|-------|-------------|
| Cloudflare | `"Cloudflare"` | `email` + `password` 或 `login_token` |
| DNSPod | `"DNSPod"` | `password` 或 `login_token` |
| DigitalOcean | `"DigitalOcean"` | `login_token` |
| AliDNS | `"AliDNS"` | `email` + `password` |
| Google Domains | `"Google"` | `email` + `password` |
| Hurricane Electric | `"HE"` | `password` |
| Dreamhost | `"Dreamhost"` | `login_token` |
| Duck DNS | `"DuckDNS"` | `login_token` |
| NoIP | `"NoIP"` | `email` + `password` |
| Scaleway | `"Scaleway"` | `login_token` |
| DynV6 | `"Dynv6"` | `login_token` |
| Linode | `"Linode"` | `login_token` |
| Strato | `"Strato"` | `password` |
| LoopiaSE | `"LoopiaSE"` | `password` |
| Infomaniak | `"Infomaniak"` | `password` |
| Hetzner | `"Hetzner"` | `login_token` |
| OVH | `"OVH"` | `app_key` + `app_secret` + `consumer_key` |
| Dynu | `"Dynu"` | `password` |
| IONOS | `"IONOS"` | `login_token` |
| TransIP | `"TransIP"` | `email` + `login_token` |

**重要提示**：提供商名称区分大小写。请使用"配置值"列中的确切值。

## 迁移指南

### 从单提供商到多提供商

1. **保持现有配置正常工作**：当前设置无需任何更改
2. **逐步添加新提供商**：
   ```json
   {
     "provider": "DNSPod",        // 保持现有设置
     "login_token": "old-token",
     "providers": {               // 添加新提供商
       "Cloudflare": {
         "email": "user@example.com",
         "password": "cf-token"
       }
     },
     "domains": [
       {
         "domain_name": "old-domain.com",
         "sub_domains": ["www"]    // 使用 DNSPod（全局提供商）
       },
       {
         "domain_name": "new-domain.com", 
         "sub_domains": ["www"],
         "provider": "Cloudflare"  // 使用 Cloudflare
       }
     ]
   }
   ```

3. **完成迁移**：一旦所有域名都指定了提供商，就可以移除全局提供商

## 配置示例

查看示例配置文件：

- `configs/config_multi_sample.json` - 完整的多提供商设置
- `configs/config_multi_sample.yaml` - 多提供商设置的 YAML 版本
- `configs/config_legacy_compatible.json` - 显示向后兼容性
- `configs/config_mixed_sample.json` - 混合传统 + 新提供商配置

## 日志记录和通知

多提供商配置在日志消息和通知中包含提供商信息：

```
INFO [2024-01-01T12:00:00Z] 已初始化提供商: cloudflare
INFO [2024-01-01T12:00:00Z] 已初始化提供商: dnspod  
INFO [2024-01-01T12:00:00Z] [ www, api ] of example.com (通过 cloudflare)
```

## 优势

1. **整合**：从一个 GoDNS 实例管理多个 DNS 提供商
2. **灵活性**：为每个域名使用最佳提供商
3. **冗余性**：将域名分布在提供商之间以提高弹性
4. **迁移**：逐步在提供商之间移动域名
5. **成本优化**：根据定价/功能使用不同提供商
6. **合规性**：满足地理分布的要求

## 故障排除

### 找不到提供商错误
```
ERROR provider 'cloudflare' not found for domain example.com
```
**解决方案**：
1. 确保提供商已在 `providers` 部分中配置
2. 检查您是否使用了正确的区分大小写的提供商名称（例如，`"Cloudflare"` 而不是 `"cloudflare"`）

### 身份验证失败
```  
ERROR failed to create provider Cloudflare: authentication failed
```
**解决方案**：验证提供商配置部分中的凭据，并确保您使用的是该提供商的正确身份验证方法。

### 混合配置问题
如果域名没有指定 `provider` 字段，它将使用全局 `provider`。确保：
1. 使用混合配置时设置了全局 `provider`
2. 为全局提供商提供了所有必需的凭据