```text
 ██████╗  ██████╗ ██████╗ ███╗   ██╗███████╗
██╔════╝ ██╔═══██╗██╔══██╗████╗  ██║██╔════╝
██║  ███╗██║   ██║██║  ██║██╔██╗ ██║███████╗
██║   ██║██║   ██║██║  ██║██║╚██╗██║╚════██║
╚██████╔╝╚██████╔╝██████╔╝██║ ╚████║███████║
 ╚═════╝  ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝╚══════╝
```

[![Apache licensed][9]][10] [![Docker][3]][4] [![Go Report Card][11]][12] [![GoDoc][13]][14]

[3]: https://img.shields.io/docker/image-size/timothyye/godns/latest
[4]: https://hub.docker.com/r/timothyye/godns
[9]: https://img.shields.io/badge/license-Apache-blue.svg
[10]: LICENSE
[11]: https://goreportcard.com/badge/github.com/timothyye/godns
[12]: https://goreportcard.com/report/github.com/timothyye/godns
[13]: https://godoc.org/github.com/TimothyYe/godns?status.svg
[14]: https://godoc.org/github.com/TimothyYe/godns

[GoDNS](https://github.com/TimothyYe/godns) 是一个动态 DNS (DDNS) 客户端工具。它是用 [Go](https://golang.org) 重写的我早期的 [DynDNS](https://github.com/TimothyYe/DynDNS) 开源项目。

<img src="https://github.com/TimothyYe/godns/blob/master/assets/snapshots/web-panel.jpg?raw=true" />

- [支持的 DNS 提供商](#支持的-dns-提供商)
- [支持的平台](#支持的平台)
- [前提条件](#前提条件)
- [安装](#安装)
- [使用方法](#使用方法)
- [配置](#配置)
  - [概述](#概述)
  - [配置文件格式](#配置文件格式)
  - [动态加载配置](#动态加载配置)
  - [配置属性](#配置属性)
  - [更新根域名](#更新根域名)
  - [配置示例](#配置示例)
    - [Cloudflare](#cloudflare)
    - [DigitalOcean](#digitalocean)
    - [DNSPod](#dnspod)
    - [Dreamhost](#dreamhost)
    - [Dynv6](#dynv6)
    - [Google Domains](#google-domains)
    - [AliDNS](#alidns)
    - [DuckDNS](#duckdns)
    - [No-IP](#no-ip)
    - [HE.net](#henet)
    - [Scaleway](#scaleway)
    - [Linode](#linode)
    - [Strato](#strato)
    - [LoopiaSE](#loopiase)
    - [Infomaniak](#infomaniak)
    - [Hetzner](#hetzner)
    - [OVH](#ovh)
    - [Dynu](#dynu)
    - [IONOS](#ionos)
    - [TransIP](#transip)
  - [通知](#通知)
    - [电子邮件](#电子邮件)
    - [Telegram](#telegram)
    - [Slack](#slack)
    - [Discord](#discord)
    - [Pushover](#pushover)
    - [Bark](#bark)
  - [Webhook](#webhook)
    - [使用 HTTP GET 请求的 Webhook](#使用-http-get-请求的-webhook)
    - [使用 HTTP POST 请求的 Webhook](#使用-http-post-请求的-webhook)
  - [杂项主题](#杂项主题)
    - [IPv6 支持](#ipv6-支持)
    - [网络接口 IP 地址](#网络接口-ip-地址)
    - [SOCKS5 代理支持](#socks5-代理支持)
    - [显示调试信息](#显示调试信息)
    - [从 RouterOS 获取 IP](#从-routeros-获取-ip)
    - [多个 API URL](#多个-api-url)
    - [推荐的 API](#推荐的-api)
- [Web 面板](#web-面板)
- [运行 GoDNS](#运行-godns)
  - [手动运行](#手动运行)
  - [作为手动守护进程](#作为手动守护进程)
  - [作为托管守护进程（使用 upstart）](#作为托管守护进程使用-upstart)
  - [作为托管守护进程（使用 systemd）](#作为托管守护进程使用-systemd)
  - [作为托管守护进程（使用 procd）](#作为托管守护进程使用-procd)
  - [作为 Docker 容器](#作为-docker-容器)
  - [作为 Windows 服务](#作为-windows-服务)
- [贡献](#贡献)
  - [设置前端开发环境](#设置前端开发环境)
  - [构建前端](#构建前端)
  - [运行前端](#运行前端)
- [特别感谢](#特别感谢)

---

## 支持的 DNS 提供商

| 提供商                                |     IPv4 支持      |     IPv6 支持      |       根域名       |       子域名       |
| ------------------------------------- | :----------------: | :----------------: | :----------------: | :----------------: |
| [Cloudflare][cloudflare]              | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| [DigitalOcean][digitalocean]          | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| [Google Domains][google.domains]      | :white_check_mark: | :white_check_mark: |        :x:         | :white_check_mark: |
| [DNSPod][dnspod]                      | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| [Dynv6][dynv6]                        | :white_check_mark: | :white_check_mark: |        :x:         | :white_check_mark: |
| [HE.net (Hurricane Electric)][he.net] | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| [AliDNS][alidns]                      | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| [DuckDNS][duckdns]                    | :white_check_mark: | :white_check_mark: |        :x:         | :white_check_mark: |
| [Dreamhost][dreamhost]                | :white_check_mark: | :white_check_mark: |        :x:         | :white_check_mark: |
| [No-IP][no-ip]                        | :white_check_mark: | :white_check_mark: |        :x:         | :white_check_mark: |
| [Scaleway][scaleway]                  | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| [Linode][linode]                      | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| [Strato][strato]                      | :white_check_mark: | :white_check_mark: |        :x:         | :white_check_mark: |
| [LoopiaSE][loopiase]                  | :white_check_mark: | :white_check_mark: |        :x:         | :white_check_mark: |
| [Infomaniak][infomaniak]              | :white_check_mark: | :white_check_mark: |        :x:         | :white_check_mark: |
| [Hetzner][hetzner]                    | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| [OVH][ovh]                            | :white_check_mark: | :white_check_mark: |        :x:         | :white_check_mark: |
| [Dynu][dynu]                          | :white_check_mark: | :white_check_mark: |        :x:         | :white_check_mark: |
| [IONOS][ionos]                        | :white_check_mark: | :white_check_mark: |        :x:         | :white_check_mark: |
| [TransIP][transip]                    | :white_check_mark: | :white_check_mark: |        :x:         | :white_check_mark: |

[cloudflare]: https://cloudflare.com
[digitalocean]: https://digitalocean.com
[google.domains]: https://domains.google
[dnspod]: https://www.dnspod.cn
[dynv6]: https://dynv6.com
[he.net]: https://dns.he.net
[alidns]: https://help.aliyun.com/product/29697.html
[duckdns]: https://www.duckdns.org
[dreamhost]: https://www.dreamhost.com
[no-ip]: https://www.noip.com
[scaleway]: https://www.scaleway.com/
[linode]: https://www.linode.com
[strato]: https://strato.de
[loopiase]: https://www.loopia.se/
[infomaniak]: https://www.infomaniak.com/
[hetzner]: https://hetzner.com/
[ovh]: https://www.ovh.com
[dynu]: https://www.dynu.com/
[ionos]: https://www.ionos.com/
[transip]: https://www.transip.net/

提示：您可以关注此 [问题](https://github.com/TimothyYe/godns/issues/76) 查看根域名 DDNS 的当前状态。

## 支持的平台

- Linux
- MacOS
- ARM Linux（如 Raspberry Pi 等）
- Windows
- MIPS32 平台

  要为 MIPS（mips 或 mipsle）编译二进制文件，请运行：

  ```bash
  GOOS=linux GOARCH=mips/mipsle GOMIPS=softfloat go build -a
  ```

  该二进制文件也可以在路由器上运行。

## 前提条件

要使用 GoDNS，假设：

- 您已注册（现在拥有）一个域名
- 域名已委托给受支持的 [DNS 提供商](#支持的-dns-提供商)（即它有指向受支持提供商的 nameserver `NS` 记录）

或者，您可以登录 [DuckDNS](https://www.duckdns.org)（使用社交账户），免费获取 duckdns.org 域名下的子域名。

## 安装

通过运行以下命令构建 GoDNS（从仓库根目录）：

```bash
cd cmd/godns        # 进入 GoDNS 目录
go mod download     # 获取依赖项
go build            # 构建
```

您还可以从 [releases](https://github.com/TimothyYe/godns/releases) 下载已编译的二进制文件。

## 使用方法

通过运行以下命令打印使用/帮助信息：

```bash
$ ./godns -h
Usage of ./godns:
  -c string
        指定配置文件（默认 "./config.json"）
  -h    显示帮助
```

## 配置

### 概述

- 复制 [config_sample.json](configs/config_sample.json) 并命名为 `config.json`，或者复制 [config_sample.yaml](configs/config_sample.yaml) 并命名为 `config.yaml`。
- 配置您的提供商、域名/子域名信息、凭据等。
- 配置通知媒介（例如 SMTP 接收电子邮件），以便在您的 IP 地址更改时收到通知
- 将文件放置在 GoDNS 的同一目录中，或者使用 `-c=path/to/your/file.json` 选项

### 配置文件格式

GoDNS 支持两种不同的配置文件格式：

- JSON
- YAML

默认情况下，GoDNS 使用 `JSON` 配置文件。但是，您可以通过 `./godns -c /path/to/config.yaml` 指定使用 `YAML` 格式。

### 动态加载配置

GoDNS 支持动态加载配置。如果您修改了配置文件，GoDNS 将自动重新加载配置并应用更改。

### 配置属性

- `provider` — 使用的一个 [支持的提供商](#支持的-dns-提供商)：`Cloudflare`、`Google`、`DNSPod`、`AliDNS`、`HE`、`DuckDNS` 或 `Dreamhost`。
- `email` — DNS 提供商的电子邮件或账户名。
- `password` — DNS 提供商的密码。
- `login_token` — DNS 提供商的 API 令牌。
- `domains` — 域名列表，包含您的子域名。
- `ip_urls` — 用于获取公共 IPv4 地址的 URL 数组。
- `ipv6_urls` — 用于获取公共 IPv6 地址的 URL 数组。
- `ip_type` — 决定使用 IPv4 还是 IPv6 的开关（当 [支持](#支持的-dns-提供商) 时）。可用值：`IPv4` 或 `IPv6`。
- `interval` — 公共 IP 更新的频率（以秒为单位）。
- `socks5_proxy` — Socks5 代理服务器。
- `resolver` — 要使用的公共 DNS 服务器地址。例如，要使用 [Google 的公共 DNS](https://developers.google.com/speed/public-dns/docs/using)，您可以在使用 GoDNS 的 IPv4 模式时设置 `8.8.8.8`，或在 IPv6 模式时设置 `2001:4860:4860::8888`。
- `skip_ssl_verify` - 跳过对 https 请求的 SSL 证书验证。

### 更新根域名

只需将 `@` 放入 `sub_domains`，例如：

```json
"domains": [{
      "domain_name": "example.com",
      "sub_domains": ["@"]
    }]
```

### 配置示例

#### Cloudflare

对于 Cloudflare，您需要提供电子邮件和全局 API 密钥作为密码（或使用 API 令牌），并配置所有域名和子域名。

通过设置选项 `proxied = true`，记录将获得 Cloudflare 的性能和安全优势。此选项仅适用于 Cloudflare。

<details>
<summary>使用电子邮件和全局 API 密钥</summary>

```json
{
  "provider": "Cloudflare",
  "email": "you@example.com",
  "password": "Global API Key",
  "domains": [
    {
      "domain_name": "example.com",
      "sub_domains": ["www", "test"]
    },
    {
      "domain_name": "example2.com",
      "sub_domains": ["www", "test"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_urls": ["https://api.ip.sb/ip"],
  "ip_type": "IPv4",
  "interval": 300,
  "socks5_proxy": "",
  "proxied": false
}
```

</details>

<details>
<summary>使用 API 令牌</summary>

```json
{
  "provider": "Cloudflare",
  "login_token": "API Token",
  "domains": [
    {
      "domain_name": "example.com",
      "sub_domains": ["www", "test"]
    },
    {
      "domain_name": "example2.com",
      "sub_domains": ["www", "test"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_urls": ["https://api.ip.sb/ip"],
  "ip_type": "IPv4",
  "interval": 300,
  "socks5_proxy": ""
}
```

</details>

#### DNSPod

对于 DNSPod，您需要提供您的 API 令牌（您可以在[这里](https://www.dnspod.cn/console/user/security)创建），并配置所有域名和子域名。

<details>
<summary>示例</summary>

```json
{
  "provider": "DNSPod",
  "login_token": "your_id,your_token",
  "domains": [
    {
      "domain_name": "example.com",
      "sub_domains": ["www", "test"]
    },
    {
      "domain_name": "example2.com",
      "sub_domains": ["www", "test"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_urls": ["https://api.ip.sb/ip"],
  "ip_type": "IPv4",
  "interval": 300,
  "socks5_proxy": ""
}
```

</details>

#### DigitalOcean

对于 DigitalOcean，您需要提供一个具有 `domain` 范围的 API 令牌（您可以在[这里](https://cloud.digitalocean.com/account/api/tokens/new)创建），并配置所有域名和子域名。

<details>
<summary>示例</summary>

```json
{
  "provider": "DigitalOcean",
  "login_token": "dop_v1_00112233445566778899aabbccddeeff",
  "domains": [
    {
      "domain_name": "example.com",
      "sub_domains": ["@", "www"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_urls": ["https://api.ip.sb/ip"],
  "ip_type": "IPv4",
  "interval": 300
}
```

</details>

#### Dreamhost

对于 Dreamhost，您需要提供您的 API 令牌（您可以在[这里](https://panel.dreamhost.com/?tree=home.api)创建），并配置所有域名和子域名。

<details>
<summary>示例</summary>

```json
{
  "provider": "Dreamhost",
  "login_token": "your_api_key",
  "domains": [
    {
      "domain_name": "example.com",
      "sub_domains": ["www", "test"]
    },
    {
      "domain_name": "example2.com",
      "sub_domains": ["www", "test"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_urls": ["https://api.ip.sb/ip"],
  "ip_type": "IPv4",
  "interval": 300,
  "resolver": "ns1.dreamhost.com",
  "socks5_proxy": ""
}
```

</details>

#### Dynv6

对于 Dynv6，只需提供 `token`，配置 1 个默认域名和子域名。

<details>
<summary>示例</summary>

```json
{
  "provider": "Dynv6",
  "password": "",
  "login_token": "1234567ABCDEFGabcdefg123456789",
  "domains": [
    {
      "domain_name": "dynv6.net",
      "sub_domains": ["myname"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_urls": ["https://api.ip.sb/ip"],
  "ip_type": "IPv4",
  "interval": 300,
  "socks5_proxy": ""
}
```

</details>

#### Google Domains

对于 Google Domains，您需要提供电子邮件和密码，并配置所有域名和子域名。

<details>
<summary>示例</summary>

```json
{
  "provider": "Google",
  "email": "Your_Username",
  "password": "Your_Password",
  "domains": [
    {
      "domain_name": "example.com",
      "sub_domains": ["www", "test"]
    },
    {
      "domain_name": "example2.com",
      "sub_domains": ["www", "test"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_urls": ["https://api.ip.sb/ip"],
  "ip_type": "IPv4",
  "interval": 300,
  "socks5_proxy": ""
}
```

</details>

#### AliDNS

对于 AliDNS，您需要提供 `AccessKeyID` 和 `AccessKeySecret` 作为 `email` 和 `password`，并配置所有域名和子域名。

<details>
<summary>示例</summary>

```json
{
  "provider": "AliDNS",
  "email": "AccessKeyID",
  "password": "AccessKeySecret",
  "login_token": "",
  "domains": [
    {
      "domain_name": "example.com",
      "sub_domains": ["www", "test"]
    },
    {
      "domain_name": "example2.com",
      "sub_domains": ["www", "test"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_urls": ["https://api.ip.sb/ip"],
  "ip_type": "IPv4",
  "interval": 300,
  "socks5_proxy": ""
}
```

</details>

#### DuckDNS

对于 DuckDNS，只需提供 `token`，配置 1 个默认域名和子域名。

<details>
<summary>示例</summary>

```json
{
  "provider": "DuckDNS",
  "password": "",
  "login_token": "3aaaaaaaa-f411-4198-a5dc-8381cac61b87",
  "domains": [
    {
      "domain_name": "www.duckdns.org",
      "sub_domains": ["myname"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_urls": ["https://api.ip.sb/ip"],
  "ip_type": "IPv4",
  "interval": 300,
  "socks5_proxy": ""
}
```

</details>

#### No-IP

<details>
<summary>示例</summary>

```json
{
  "provider": "NoIP",
  "email": "mail@example.com",
  "password": "YourPassword",
  "domains": [
    {
      "domain_name": "ddns.net",
      "sub_domains": ["timothyye6"]
    }
  ],
  "ip_type": "IPv4",
  "ip_urls": ["https://api.ip.sb/ip"],
  "resolver": "8.8.8.8",
  "interval": 300,
  "socks5_proxy": ""
}
```

</details>

#### HE.net

对于 HE，不需要电子邮件，只需填写 DDNS 密钥作为密码，并配置所有域名和子域名。

<details>
<summary>示例</summary>

```json
{
  "provider": "HE",
  "password": "Your DDNS Key",
  "login_token": "",
  "domains": [
    {
      "domain_name": "example.com",
      "sub_domains": ["www", "test"]
    },
    {
      "domain_name": "example2.com",
      "sub_domains": ["www", "test"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_urls": ["https://api.ip.sb/ip"],
  "ip_type": "IPv4",
  "interval": 300,
  "socks5_proxy": ""
}
```

</details>

<details>
<summary>提供商配置</summary>

添加一个新的 "A 记录" 并确保勾选 "启用动态 DNS 条目"：

<img src="assets/snapshots/he1.png" width="640" />

填写您自己的 DDNS 密钥或为这个新创建的 "A 记录" 生成一个随机 DDNS 密钥：

<img src="assets/snapshots/he2.png" width="640" />

记住 DDNS 密钥并在配置文件中设置到 `password` 属性中。

**注意**：如果您有多个域名或子域名，请确保它们的 DDNS 密钥相同。

</details>

#### Scaleway

对于 Scaleway，您需要提供 API 密钥作为 `login_token`（[如何生成 API 密钥](https://www.scaleway.com/en/docs/generate-api-keys/)），并配置域名和子域名。`domain_name` 应等于 Scaleway 中的 DNS 区域或根域名。DNS 记录的 TTL 将设置为 `interval` 值。确保相关子域名的 `A` 或 `AAAA` 记录存在，这些可以在 [Scaleway 控制台](https://www.scaleway.com/en/docs/scaleway-dns/#-Managing-Records) 中设置。

<details>
<summary>示例</summary>

```json
{
  "provider": "Scaleway",
  "login_token": "API Secret Key",
  "domains": [
    {
      "domain_name": "example.com",
      "sub_domains": ["www", "@"]
    },
    {
      "domain_name": "samplednszone.example.com",
      "sub_domains": ["www", "test"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_urls": ["https://api.ip.sb/ip"],
  "ip_type": "IPv4",
  "interval": 300
}
```

</details>

#### Linode

要与 Linode API 进行身份验证，您需要提供一个具有“读/写”访问权限的个人访问令牌，范围为“Domain”。Linode 有一个[关于创建访问令牌的帮助页面](https://www.linode.com/docs/guides/getting-started-with-the-linode-api/)。将此令牌传入配置文件中的 `login_token` 字段。

配置文件中的 `domain_name` 字段必须是 Linode 管理的现有域名的名称。Linode 有一个[关于添加域名的帮助页面](https://www.linode.com/docs/guides/dns-manager/)。GoDNS Linode 处理程序不会自动创建域名，但会自动创建子域名。

GoDNS Linode 处理程序目前对 Linode DNS 记录使用固定的 30 秒 TTL。

<details>
<summary>示例</summary>

```json
{
  "provider": "Linode",
  "login_token": ${PERSONAL_ACCESS_TOKEN},
  "domains": [{
      "domain_name": "example.com",
      "sub_domains": ["www","@"]
    },{
      "domain_name": "samplednszone.example.com",
      "sub_domains": ["www","test"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_urls": ["https://api.ip.sb/ip"],
  "ip_type": "IPv4",
  "interval": 300
}
```

</details>

#### Strato

对于 Strato，您需要提供电子邮件和密码，并配置所有域名和子域名。
更多信息：[德语](https://www.strato.de/faq/hosting/so-einfach-richten-sie-dyndns-fuer-ihre-domains-ein/) [英语](https://www.strato-hosting.co.uk/faq/hosting/this-is-how-easy-it-is-to-set-up-dyndns-for-your-domains/)

<details>
<summary>示例</summary>

```json
{
  "provider": "strato",
  "password": "Your_Password",
  "domains": [
    {
      "domain_name": "example.com",
      "sub_domains": ["www", "test"]
    },
    {
      "domain_name": "example2.com",
      "sub_domains": ["www", "test"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_urls": ["https://api.ip.sb/ip"],
  "ip_type": "IPv4",
  "interval": 300,
  "socks5_proxy": ""
}
```

</details>

#### LoopiaSE

对于 LoopiaSE，您需要提供用户名和密码，并配置所有域名和子域名。
更多信息：[瑞典语](https://support.loopia.se/wiki/om-dyndns-stodet/)

<details>
<summary>示例</summary>

```json
{
  "provider": "LoopiaSE",
  "email": "Your_Username",
  "password": "Your_Password",
  "domains": [
    {
      "domain_name": "example.com",
      "sub_domains": ["www", "test"]
    },
    {
      "domain_name": "example2.com",
      "sub_domains": ["www", "test"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_urls": ["https://api.ip.sb/ip"],
  "ip_type": "IPv4",
  "interval": 300,
  "socks5_proxy": ""
}
```

</details>

#### Infomaniak

对于 Infomaniak，您需要提供用户名和密码，并配置所有域名和子域名。
更多信息：[英语](https://faq.infomaniak.com/2376)

<details>
<summary>示例</summary>

```json
{
  "provider": "Infomaniak",
  "email": "Your_Username",
  "password": "Your_Password",
  "domains": [
    {
      "domain_name": "example.com",
      "sub_domains": ["www", "test"]
    },
    {
      "domain_name": "example2.com",
      "sub_domains": ["www", "test"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_urls": ["https://api.ip.sb/ip"],
  "ip_type": "IPv4",
  "interval": 300,
  "socks5_proxy": ""
}
```

</details>

#### Hetzner

对于 Hetzner，您必须创建一个访问令牌。这可以在 DNS 控制台中完成。
（左上角的个人图标 --> API 令牌）
注意：如果一个域名有多个记录，**只有第一个**记录会被更新。
确保只有一个记录。

<details>
<summary>示例</summary>

```json
{
  "provider": "hetzner",
  "login_token": "<token>",
  "domains": [
    {
      "domain_name": "example.com",
      "sub_domains": ["www", "test"]
    },
    {
      "domain_name": "example2.com",
      "sub_domains": ["www", "test"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_urls": ["https://api.ip.sb/ip"],
  "ip_type": "IPv4"
}
```

</details>

#### OVH

对于 OVH，您需要提供 Consumerkey、Appsecret 和 Appkey，并配置所有域名和子域名。
所需的值可以通过访问[此网站](https://www.ovh.com/auth/api/createToken)获取
权限应在 GET、POST 和 PUT 上设置为 '\*'
更多信息：[help.ovhcloud.com](https://help.ovhcloud.com/csm/en-gb-api-getting-started-ovhcloud-api?id=kb_article_view&sysparm_article=KB0042784)

<details>
<summary>示例</summary>

```json
{
  "provider": "OVH",
  "consumer_key": "e389ac80cc8da9c7451bc7b8f171bf4f",
  "app_secret": "d1ffee354d3643d70deaab48a09131fd",
  "app_key": "cd338839d6472064",
  "domains": [
    {
      "domain_name": "example.com",
      "sub_domains": ["www", "test"]
    },
    {
      "domain_name": "example2.com",
      "sub_domains": ["www", "test"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_urls": ["https://api.ip.sb/ip"],
  "ip_type": "IPv4",
  "interval": 300,
  "socks5_proxy": ""
}
```

</details>

#### Dynu

对于 Dynu，您需要配置 `password`，配置 1 个默认域名和子域名。

<details>
<summary>示例</summary>

```json
{
  "provider": "Dynu",
  "password": "Your_Password",
  "domains": [
    {
      "domain_name": "your_domain.com",
      "sub_domains": ["your_subdomain"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_urls": ["https://api.ip.sb/ip"],
  "ip_type": "IPv4",
  "interval": 300,
  "socks5_proxy": ""
}
```

</details>

#### IONOS

这是针对 IONOS 托管服务的，**不是** IONOS 云。
您需要[注册 IONOS API 访问托管服务](https://my.ionos.com/shop/product/ionos-api)，然后创建一个[API 密钥](https://developer.hosting.ionos.com/keys)。
您可以在[IONOS API 文档](https://developer.hosting.ionos.com/docs/getstarted)中找到完整指南。
**注意**：GoDNS 使用的 API 密钥必须遵循上述文档中描述的 `publicprefix.secret` 形式。

<details>
<summary>示例</summary>

```yaml
provider: IONOS
login_token: publicprefix.secret
domains:
  - domain_name: example.com
    sub_domains:
      - somesubdomain
      - anothersubdomain
resolver: 1.1.1.1
ip_urls:
  - https://api.ipify.org
ip_type: IPv4
interval: 300
socks5_proxy: ""
```

</details>

#### TransIP

对于 TransIP，您需要提供您的 API 私钥作为 `login_token`，用户名作为 `email`，并配置所有域名和子域名。

<details>
<summary>示例</summary>

```yaml
provider: TransIP
email: account_name
login_token: api_key
domains:
  - domain_name: example.com
    sub_domains:
      - "@"
      - somesubdomain
      - anothersubdomain
resolver: 1.1.1.1
ip_urls:
  - https://api.ipify.org
ip_type: IPv4
interval: 300
socks5_proxy: ""
```

<details>

### 通知

GoDNS 可以在 IP 更改时发送通知。

#### 电子邮件

电子邮件通过 [SMTP](https://en.wikipedia.org/wiki/Simple_Mail_Transfer_Protocol) 发送。使用以下片段更新您的配置：

```json
  "notify": {
    "mail": {
      "enabled": true,
      "smtp_server": "smtp.example.com",
      "smtp_username": "user",
      "smtp_password": "password",
      "smtp_port": 25,
      "send_from": "my_mail@example.com"
      "send_to": "my_mail@example.com"
    }
  }
```

每次 IP 更改时，您将收到如下电子邮件：

<img src="https://github.com/TimothyYe/godns/blob/master/assets/snapshots/mail.png?raw=true" />

#### Telegram

要在 IP 更改时接收 [Telegram](https://telegram.org/) 消息，使用以下片段更新您的配置：

```json
  "notify": {
    "telegram": {
      "enabled": true,
      "bot_api_key": "11111:aaaa-bbbb",
      "chat_id": "-123456",
      "message_template": "域名 *{{ .Domain }}* 已更新为 %0A{{ .CurrentIP }}",
      "use_proxy": false
    },
  }
```

`message_template` 属性支持 [markdown](https://www.markdownguide.org)。新行需要用 `%0A` 转义。

#### Slack

要在 IP 更改时接收 [Slack](https://slack.com) 消息，使用以下片段更新您的配置：

```json
  "notify": {
    "slack": {
      "enabled": true,
      "bot_api_token": "xoxb-xxx",
      "channel": "your_channel",
      "message_template": "域名 *{{ .Domain }}* 已更新为 \n{{ .CurrentIP }}",
      "use_proxy": false
    },
  }
```

`message_template` 属性支持 [markdown](https://www.markdownguide.org)。新行需要用 `\n` 转义。

#### Discord

要在 IP 更改时接收 [Discord](https://discord.gg) 消息，使用以下片段更新您的配置：

```json
  "notify": {
    "discord": {
          "enabled": true,
          "bot_api_token": "discord_bot_token",
          "channel": "your_channel",
          "message_template": "(可选) 域名 *{{ .Domain }}* 已更新为 \n{{ .CurrentIP }}",
        }
  }
```

#### Pushover

要在 IP 更改时接收 [Pushover](https://pushover.net/) 消息，使用以下片段更新您的配置：

```json
  "notify": {
    "pushover": {
      "enabled": true,
      "token": "abcdefghijklmnopqrstuvwxyz1234",
      "user": "abcdefghijklmnopqrstuvwxyz1234",
      "message_template": "",
      "device": "",
      "title": "",
      "priority": 0,
      "html": 1
    }
  }
```

如果 `html` 参数为 `1`，`message_template` 属性支持 [html](https://pushover.net/api#html)。如果留空，将使用默认消息。
如果 `device` 和 `title` 参数留空，Pushover 将选择默认值[参见](https://pushover.net/api#messages)。有关优先级参数的更多详细信息
可以在 Pushover [API 描述](https://pushover.net/api#priority) 中找到。

#### Bark

要在 IP 更改时接收 [Bark](https://bark.day.app/) 消息，使用以下片段更新您的配置：

```json
  "notify": {
    "bark": {
      "enabled": true,
      "server": "https://api.day.app",
      "device_keys": "",
      "params": "{ \"isArchive\": 1, \"action\": \"none\" }"
    }
  }
```

`server` Bark 服务器地址，可使用官方默认服务器 `https://api.day.app`，也可设置为自建服务器地址。
`device_keys` 设备 key，支持多个（英文逗号分隔），多个时，用于批量推送
`params` Bark 请求参数，可参考 [Bark API](https://bark.day.app/#/tutorial?id=%e8%af%b7%e6%b1%82%e5%8f%82%e6%95%b0)
更多内容请参阅 [Bark 官方文档](https://bark.day.app/)

### Webhook

Webhook 是 GoDNS 提供的另一个功能，用于在 IP 更改时向其他应用程序发送通知。GoDNS 通过 HTTP `GET` 或 `POST` 请求向目标 URL 发送通知。

配置部分 `webhook` 用于自定义 webhook 请求。通常，有 2 个字段用于 webhook 请求：

> - `url`：发送 webhook 请求的目标 URL。
> - `request_body`：发送 `POST` 请求的内容，如果此字段为空，则发送 HTTP GET 请求而不是 HTTP POST 请求。

可用变量：

> - `Domain`：当前域名。
> - `IP`：新 IP
