```text
 ██████╗  ██████╗ ██████╗ ███╗   ██╗███████╗
██╔════╝ ██╔═══██╗██╔══██╗████╗  ██║██╔════╝
██║  ███╗██║   ██║██║  ██║██╔██╗ ██║███████╗
██║   ██║██║   ██║██║  ██║██║╚██╗██║╚════██║
╚██████╔╝╚██████╔╝██████╔╝██║ ╚████║███████║
 ╚═════╝  ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝╚══════╝
```

[![Apache licensed][9]][10] [![Docker][3]][4] [![Go Report Card][11]][12] [![Cover.Run][15]][16] [![GoDoc][13]][14]

[3]: https://img.shields.io/docker/image-size/timothyye/godns/latest
[4]: https://hub.docker.com/r/timothyye/godns
[9]: https://img.shields.io/badge/license-Apache-blue.svg
[10]: LICENSE
[11]: https://goreportcard.com/badge/github.com/timothyye/godns
[12]: https://goreportcard.com/report/github.com/timothyye/godns
[13]: https://godoc.org/github.com/TimothyYe/godns?status.svg
[14]: https://godoc.org/github.com/TimothyYe/godns
[15]: https://img.shields.io/badge/cover.run-88.2%25-green.svg
[16]: https://cover.run/go/github.com/timothyye/godns

[GoDNS](https://github.com/TimothyYe/godns) is a dynamic DNS (DDNS) client tool. It is a rewrite in [Go](https://golang.org) of my early [DynDNS](https://github.com/TimothyYe/DynDNS) open source project.

---

- [Supported DNS Providers](#supported-dns-providers)
- [Supported Platforms](#supported-platforms)
- [Pre-conditions](#pre-conditions)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
  - [Overview](#overview)
  - [Configuration file format](#configuration-file-format)
  - [Configuration properties](#configuration-properties)
  - [Update root domain](#update-root-domain)
  - [Configuration examples](#configuration-examples)
    - [Cloudflare](#cloudflare)
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
    - [OVH](#ovh)
    - [Dynu](#dynu)
  - [Notifications](#notifications)
    - [Email](#email)
    - [Telegram](#telegram)
    - [Slack](#slack)
    - [Discord](#discord)
    - [Pushover](#pushover)
  - [Webhook](#webhook)
    - [HTTP GET Request](#webhook-with-http-get-reqeust)
    - [HTTP POST Request](#webhook-with-http-post-request)
  - [Miscellaneous topics](#miscellaneous-topics)
    - [IPv6 support](#ipv6-support)
    - [Network interface IP address](#network-interface-ip-address)
    - [SOCKS5 proxy support](#socks5-proxy-support)
    - [Display debug info](#display-debug-info)
    - [Recommended APIs](#recommended-apis)
- [Running GoDNS](#running-godns)
  - [Manually](#manually)
  - [As a manual daemon](#as-a-manual-daemon)
  - [As a managed daemon (with upstart)](#as-a-managed-daemon-with-upstart)
  - [As a managed daemon (with systemd)](#as-a-managed-daemon-with-systemd)
  - [As a Docker container](#as-a-docker-container)
  - [As a Windows service](#as-a-windows-service)
- [Special Thanks](#special-thanks)

---

## Supported DNS Providers

| Provider                              |    IPv4 support    |    IPv6 support    |    Root Domain     |     Subdomains     |
| ------------------------------------- | :----------------: | :----------------: | :----------------: | :----------------: |
| [Cloudflare][cloudflare]              | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
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

[cloudflare]: https://cloudflare.com
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

Tip: You can follow this [issue](https://github.com/TimothyYe/godns/issues/76) to view the current status of DDNS for root domains.

## Supported Platforms

- Linux
- MacOS
- ARM Linux (Raspberry Pi, etc.)
- Windows
- MIPS32 platform

  To compile binaries for MIPS (mips or mipsle), run:

  ```bash
  GOOS=linux GOARCH=mips/mipsle GOMIPS=softfloat go build -a
  ```

  The binary can run on routers as well.

## Pre-conditions

To use GoDNS, it is assumed:

- You registered (now own) a domain
- Domain was delegated to a supported [DNS provider](#supported-dns-providers) (i.e. it has nameserver `NS` records pointing at a supported provider)

Alternatively, you can sign in to [DuckDNS](https://www.duckdns.org) (with a social account) and get a subdomain on the duckdns.org domain for free.

## Installation

Build GoDNS by running (from the root of the repository):

```bash
cd cmd/godns        # go to the GoDNS directory
go mod download     # get dependencies
go build            # build
```

You can also download a compiled binary from the [releases](https://github.com/TimothyYe/godns/releases).

## Usage

Print usage/help by running:

```bash
$ ./godns -h
Usage of ./godns:
  -c string
        Specify a config file (default "./config.json")
  -h    Show help
```

## Configuration

### Overview

- Make a copy of [config_sample.json](configs/config_sample.json) and name it as `config.json`, or make a copy of [config_sample.yaml](configs/config_sample.yaml) and name it as `config.yaml`.
- Configure your provider, domain/subdomain info, credentials, etc.
- Configure a notification medium (e.g. SMTP to receive emails) to get notified when your IP address changes
- Place the file in the same directory of GoDNS or use the `-c=path/to/your/file.json` option

### Configuration file format

GoDNS supports 2 different configuration file formats:

- JSON
- YAML

By default, GoDNS uses `JSON` config file. However, you can specify to use the `YAML` format via: `./godns -c /path/to/config.yaml`

### Configuration properties

- `provider` — One of the [supported provider to use](#supported-dns-providers): `Cloudflare`, `Google`, `DNSPod`, `AliDNS`, `HE`, `DuckDNS` or `Dreamhost`.
- `email` — Email or account name of the DNS provider.
- `password` — Password of the DNS provider.
- `login_token` — API token of the DNS provider.
- `domains` — Domains list, with your sub domains.
- `ip_urls` — A URL array for fetching one's public IPv4 address.
- `ipv6_urls` — A URL array for fetching one's public IPv6 address.
- `ip_type` — Switch deciding if IPv4 or IPv6 should be used (when [supported](#supported-dns-providers)). Available values: `IPv4` or `IPv6`.
- `interval` — How often (in seconds) the public IP should be updated.
- `socks5_proxy` — Socks5 proxy server.
- `resolver` — Address of a public DNS server to use. For instance to use [Google's public DNS](https://developers.google.com/speed/public-dns/docs/using), you can set `8.8.8.8` when using GoDNS in IPv4 mode or `2001:4860:4860::8888` in IPv6 mode.

### Update root domain

By simply putting `@` into `sub_domains`, for example:

```json
"domains": [{
      "domain_name": "example.com",
      "sub_domains": ["@"]
    }]
```

### Configuration examples

#### Cloudflare

For Cloudflare, you need to provide the email & Global API Key as password (or to use the API token) and config all the domains & subdomains.

By setting the option `proxied = true`, the record is receiving the performance and security benefits of Cloudflare. This option is only available for Cloudflare.

<details>
<summary>Using email & Global API Key</summary>

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
<summary>Using the API Token</summary>

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

For DNSPod, you need to provide your API Token(you can create it [here](https://www.dnspod.cn/console/user/security)), and config all the domains & subdomains.

<details>
<summary>Example</summary>

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

#### Dreamhost

For Dreamhost, you need to provide your API Token(you can create it [here](https://panel.dreamhost.com/?tree=home.api)), and config all the domains & subdomains.

<details>
<summary>Example</summary>

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

For Dynv6, only need to provide the `token`, config 1 default domain & subdomains.

<details>
<summary>Example</summary>

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

For Google Domains, you need to provide email & password, and config all the domains & subdomains.

<details>
<summary>Example</summary>

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

For AliDNS, you need to provide `AccessKeyID` & `AccessKeySecret` as `email` & `password`, and config all the domains & subdomains.

<details>
<summary>Example</summary>

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

For DuckDNS, the only thing needed is to provide the `token`, config 1 default domain & subdomains.

<details>
<summary>Example</summary>

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
<summary>Example</summary>

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

For HE, email is not needed, just fill the DDNS key as password, and config all the domains & subdomains.

<details>
<summary>Example</summary>

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
<summary>Provider configuration</summary>

Add a new "A record" and make sure that "Enable entry for dynamic dns" is checked:

<img src="assets/snapshots/he1.png" width="640" />

Fill in your own DDNS key or generate a random DDNS key for this new created "A record":

<img src="assets/snapshots/he2.png" width="640" />

Remember the DDNS key and set it in the `password` property in the configuration file.

**NOTICE**: If you have multiple domains or subdomains, make sure their DDNS key are the same.

</details>

#### Scaleway

For Scaleway, you need to provide an API Secret Key as the `login_token` ([How to generate an API key](https://www.scaleway.com/en/docs/generate-api-keys/)), and configure the domains and subdomains. `domain_name` should equal a DNS zone, or the root domain in Scaleway. TTL for the DNS records will be set to the `interval` value. Make sure `A` or `AAAA` records exist for the relevant sub domains, these can be set up in the [Scaleway console](https://www.scaleway.com/en/docs/scaleway-dns/#-Managing-Records).

<details>
<summary>Example</summary>
 
```json
{
  "provider": "Scaleway",
  "login_token": "API Secret Key",
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

#### Linode

To authenticate with the Linode API you will need to provide a Personal Access Token with "Read/Write" access on the "Domain" scope. Linode has [a help page about creating access tokens](https://www.linode.com/docs/guides/getting-started-with-the-linode-api/). Pass this token into the `login_token` field of the config file.

The `domain_name` field of the config file must be the name of an existing Domain managed by Linode. Linode has [a help page about adding domains](https://www.linode.com/docs/guides/dns-manager/). The GoDNS Linode handler will not create domains automatically, but will create subdomains automatically.

The GoDNS Linode handler currently uses a fixed TTL of 30 seconds for Linode DNS records.

<details>
<summary>Example</summary>
 
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

For Strato, you need to provide email & password, and config all the domains & subdomains.  
More Info: [German](https://www.strato.de/faq/hosting/so-einfach-richten-sie-dyndns-fuer-ihre-domains-ein/) [English](https://www.strato-hosting.co.uk/faq/hosting/this-is-how-easy-it-is-to-set-up-dyndns-for-your-domains/)

<details>
<summary>Example</summary>

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

For LoopiaSE, you need to provide username & password, and config all the domains & subdomains.
More info: [Swedish](https://support.loopia.se/wiki/om-dyndns-stodet/)

<details>
<summary>Example</summary>

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

For Infomaniak, you need to provide username & password, and config all the domains & subdomains.
More info: [English](https://faq.infomaniak.com/2376)

<details>
<summary>Example</summary>

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

For Hetzner, you have to create an access token. This can be done in the DNS-Console.  
(Person Icon in the top left corner --> API Tokens)  
Notice: If a domain has multiple Records **only the first** Record will be updated.  
Make shure there is just one record.

<details>
<summary>Example</summary>

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

For OVH, you need to provide a Comsumerkey, an Appsecret, an Appkey and configure all the domains & subdomains.
The neeeded values can be obtaines by visting [this site](https://www.ovh.com/auth/api/createToken)  
Rights should be '\*' on GET, POST and PUT  
More info: [help.ovhcloud.com](https://help.ovhcloud.com/csm/en-gb-api-getting-started-ovhcloud-api?id=kb_article_view&sysparm_article=KB0042784)

<details>
<summary>Example</summary>

```json
{
  "provider": "OVH",
  "comsumer_key": "e389ac80cc8da9c7451bc7b8f171bf4f",
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

For Dynu, you need to configure the `password`, config 1 default domain & subdomain.

<details>
<summary>Example</summary>

```json
{
  "provider": "Dynu",
  "password": "Your_Password",
  "domains": [
    {
      "domain_name": "your_domain.com",
      "sub_domains": [
        "your_subdomain"
      ]
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

### Notifications

GoDNS can send a notification each time the IP changes.

#### Email

Emails are sent over [SMTP](https://en.wikipedia.org/wiki/Simple_Mail_Transfer_Protocol). Update your configuration with the following snippet:

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

Each time the IP changes, you will receive an email like that:

<img src="https://github.com/TimothyYe/godns/blob/master/assets/snapshots/mail.png?raw=true" />

#### Telegram

To receive a [Telegram](https://telegram.org/) message each time the IP changes, update your configuration with the following snippet:

```json
  "notify": {
    "telegram": {
      "enabled": true,
      "bot_api_key": "11111:aaaa-bbbb",
      "chat_id": "-123456",
      "message_template": "Domain *{{ .Domain }}* updated to %0A{{ .CurrentIP }}",
      "use_proxy": false
    },
  }
```

The `message_template` property supports [markdown](https://www.markdownguide.org). New lines needs to be escaped with `%0A`.

#### Slack

To receive a [Slack](https://slack.com) message each time the IP changes, update your configuration with the following snippet:

```json
  "notify": {
    "slack": {
      "enabled": true,
      "bot_api_token": "xoxb-xxx",
      "channel": "your_channel",
      "message_template": "Domain *{{ .Domain }}* updated to \n{{ .CurrentIP }}",
      "use_proxy": false
    },
  }
```

The `message_template` property supports [markdown](https://www.markdownguide.org). New lines needs to be escaped with `\n`.

#### Discord

To receive a [Discord](https://discord.gg) message each time the IP changes, update your configuration with the following snippit:

```json
  "notify": {
    "discord": {
          "enabled": true,
          "bot_api_token": "discord_bot_token",
          "channel": "your_channel",
          "message_template": "(Optional) Domain *{{ .Domain }}* updated to \n{{ .CurrentIP }}",
        }
  }
```

#### Pushover

To receive a [Pushover](https://pushover.net/) message each time the IP changes, update your configuration with the following snippet:

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

The `message_template` property supports [html](https://pushover.net/api#html) if the `html` parameter is `1`. If it is left empty a default message will be used.
If the `device` and `title` parameters are left empty, Pushover will choose defaults [see](https://pushover.net/api#messages). More details on the priority parameter
can be found on the Pushover [API description](https://pushover.net/api#priority).

### Webhook

Webhook is another feature that GoDNS provides to deliver notifications to the other applications while the IP is changed. GoDNS delivers a notification to the target URL via an HTTP `GET` or `POST` request.

The configuration section `webhook` is used for customizing the webhook request. In general, there are 2 fields used for the webhook request:

> - `url`: The target URL for sending webhook request.
> - `request_body`: The content for sending `POST` request, if this field is empty, a HTTP GET request will be sent instead of the HTTP POST request.

Available variables:

> - `Domain`: The current domain.
> - `IP`: The new IP address.
> - `IPType`: The type of the IP: `IPV4` or `IPV6`.

#### Webhook with HTTP GET reqeust

```json
"webhook": {
  "enabled": true,
  "url": "http://localhost:5000/api/v1/send?domain={{.Domain}}&ip={{.CurrentIP}}&ip_type={{.IPType}}",
  "request_body": ""
}
```

For this example, a webhook with query string parameters will be sent to the target URL:

```
http://localhost:5000/api/v1/send?domain=ddns.example.com&ip=192.168.1.1&ip_type=IPV4
```

#### Webhook with HTTP POST request

```json
"webhook": {
  "enabled": true,
  "url": "http://localhost:5000/api/v1/send",
  "request_body": "{ \"domain\": \"{{.Domain}}\", \"ip\": \"{{.CurrentIP}}\", \"ip_type\": \"{{.IPType}}\" }"
}
```

For this example, a webhook will be triggered when the IP changes, the target URL `http://localhost:5000/api/v1/send` will receive an `HTTP POST` request with request body:

```json
{ "domain": "ddns.example.com", "ip": "192.168.1.1", "ip_type": "IPV4" }
```

### Miscellaneous topics

#### IPv6 support

Most of the [providers](#supported-dns-providers) support IPv6.

To enable the `IPv6` support of GoDNS, there are two solutions to choose from:

1. Use an online service to lookup the external IPv6

   For that:

   - Set the `ip_type` as `IPv6`, and make sure the `ipv6_urls` is configured
   - Create an `AAAA` record instead of an `A` record in your DNS provider

   <details>
   <summary>Configuration example</summary>

   ```json
   {
     "domains": [
       {
         "domain_name": "example.com",
         "sub_domains": ["ipv6"]
       }
     ],
     "resolver": "2001:4860:4860::8888",
     "ipv6_urls": ["https://api-ipv6.ip.sb/ip"],
     "ip_type": "IPv6"
   }
   ```

   </details>

2. Let GoDNS find the IPv6 of the network interface of the machine it is running on (more on that [later](#network-interface-ip-address)).

   For this to happen, just leave `ip_urls` and `ipv6_urls` empty.

   Note that the network interface must be configured with an IPv6 for this to work.

#### Network interface IP address

For some reasons if you want to get the IP address associated to a network interface (instead of performing an online lookup), you can specify it in the configuration file this way:

```json
  "ip_urls": [""],
  "ip_interface": "interface-name",
```

With `interface-name` replaced by the name of the network interface, e.g. `eth0` on Linux or `Local Area Connection` on Windows.

Note: If `ip_urls` is also specified, it will be used to perform an online lookup first and the network interface IP will be used as a fallback in case of failure.

#### SOCKS5 proxy support

You can make all remote calls go through a [SOCKS5 proxy](https://en.wikipedia.org/wiki/SOCKS#SOCKS5) by specifying it in the configuration file this way:

```json
"socks5_proxy": "127.0.0.1:7070"
"use_proxy": true
```

#### Display debug info

To display debug info, set `debug_info` as `true` to enable this feature. By default, the debug info is disabled.

```json
  "debug_info": true,
```

#### Multiple API URLs

GoDNS supports to fetch the public IP from multiple URLs via a simple round-robin algorithm. If the first URL fails, it will try the next one until it succeeds. Here is an example of the configuration:

```json
  "ip_urls": [
  "https://api.ipify.org",
  "https://myip.biturl.top",
  "https://api-ipv4.ip.sb/ip"
  ],
```

#### Recommended APIs

- https://api.ipify.org
- https://myip.biturl.top
- https://ip4.seeip.org
- https://ipecho.net/plain
- https://api-ipv4.ip.sb/ip

## Running GoDNS

There are a few ways to run GoDNS.

### Manually

Note: make sure to set the `run_once` parameter in your config file so the program will quit after the first run (the default is `false`).

Can be added to `cron` or attached to other events on your system.

```json
{
  "...": "...",
  "run_once": true
}
```

Then run

```bash
./godns
```

### As a manual daemon

```bash
nohup ./godns &
```

Note: when the program stops, it will not be restarted.

### As a managed daemon (with upstart)

1. Install `upstart` first (if not available already)
2. Copy `./config/upstart/godns.conf` to `/etc/init` (and tweak it to your needs)
3. Start the service:

   ```bash
   sudo start godns
   ```

### As a managed daemon (with systemd)

1. Install `systemd` first (it not available already)
2. Copy `./config/systemd/godns.service` to `/lib/systemd/system` (and tweak it to your needs)
3. Start the service:

   ```bash
   sudo systemctl enable godns
   sudo systemctl start godns
   ```

### As a Docker container

Available docker registries:

- https://hub.docker.com/r/timothyye/godns
- https://github.com/TimothyYe/godns/pkgs/container/godns

Visit https://hub.docker.com/r/timothyye/godns to fetch the latest docker image.  
With `/path/to/config.json` your local configuration file, run:

```bash
docker run \
-d --name godns --restart=always \
-v /path/to/config.json:/config.json \
timothyye/godns:latest
```

To run it with a `YAML` config file:

```bash
docker run \
-d --name godns \
-e CONFIG=/config.yaml \
--restart=always \
-v /path/to/config.yaml:/config.yaml \
timothyye/godns:latest
```

### As a Windows service

1. Download the latest version of [NSSM](https://nssm.cc/download)

2. In an administrative prompt, from the folder where NSSM was downloaded, e.g. `C:\Downloads\nssm\` **win64**, run:

   ```
   nssm install YOURSERVICENAME
   ```

3. Follow the interface to configure the service. In the "Application" tab just indicate where the `godns.exe` file is. Optionally you can also define a description on the "Details" tab and define a log file on the "I/O" tab. Finish by clicking on the "Install service" button.

4. The service will now start along Windows.

Note: you can uninstall the service by running:

```
nssm remove YOURSERVICENAME
```

## Special Thanks

<img src="https://i.imgur.com/xhe5RLZ.jpg" width="80px" align="right" />

Thanks JetBrains for sponsoring this project with [free open source license](https://www.jetbrains.com/community/opensource/).

> I like GoLand, it is an amazing and productive tool.
