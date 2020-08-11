```text
 ██████╗  ██████╗ ██████╗ ███╗   ██╗███████╗
██╔════╝ ██╔═══██╗██╔══██╗████╗  ██║██╔════╝
██║  ███╗██║   ██║██║  ██║██╔██╗ ██║███████╗
██║   ██║██║   ██║██║  ██║██║╚██╗██║╚════██║
╚██████╔╝╚██████╔╝██████╔╝██║ ╚████║███████║
 ╚═════╝  ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝╚══════╝
 ```

[![MIT licensed][9]][10] [![LICENSE](https://img.shields.io/badge/license-NPL%20(The%20996%20Prohibited%20License)-blue.svg)](https://github.com/996icu/996.ICU/blob/master/LICENSE) [![Docker][3]][4] [![Go Report Card][11]][12] [![Cover.Run][15]][16] [![GoDoc][13]][14]

[3]: https://images.microbadger.com/badges/image/timothyye/godns.svg
[4]: https://microbadger.com/images/timothyye/godns
[9]: https://img.shields.io/badge/license-Apache-blue.svg
[10]: LICENSE
[11]: https://goreportcard.com/badge/github.com/timothyye/godns
[12]: https://goreportcard.com/report/github.com/timothyye/godns
[13]: https://godoc.org/github.com/TimothyYe/godns?status.svg
[14]: https://godoc.org/github.com/TimothyYe/godns
[15]: https://img.shields.io/badge/cover.run-88.2%25-green.svg
[16]: https://cover.run/go/github.com/timothyye/godns

GoDNS is a dynamic DNS (DDNS) client tool, it is based on my early open source project: [DynDNS](https://github.com/TimothyYe/DynDNS). 

Now I rewrite [DynDNS](https://github.com/TimothyYe/DynDNS) by Golang and call it [GoDNS](https://github.com/TimothyYe/godns).

## Supported DNS Providers

* Cloudflare ([https://cloudflare.com](https://cloudflare.com))
* Google Domains ([https://domains.google](https://domains.google))
* DNSPod ([https://www.dnspod.cn/](https://www.dnspod.cn/))
* HE.net (Hurricane Electric) ([https://dns.he.net/](https://dns.he.net/))
* AliDNS ([https://help.aliyun.com/product/29697.html](https://help.aliyun.com/product/29697.html))
* DuckDNS ([https://www.duckdns.org](https://www.duckdns.org))
* Dreamhost ([https://www.dreamhost.com](https://www.dreamhost.com))
* No-IP ([https://www.noip.com/](https://www.noip.com))

## Supported Platforms

* Linux
* MacOS
* ARM Linux (Raspberry Pi, etc...)
* Windows
* MIPS32 platform

## MIPS32 platform

To compile binaries for MIPS (mips or mipsle):

```bash
GOOS=linux GOARCH=mips/mipsle GOMIPS=softfloat go build -a
```

And the binary can run well on routers.  

## Pre-condition

* Register and own a domain.

* Domain's nameserver points to [DNSPod](https://www.dnspod.cn/) or [HE.net](https://dns.he.net/) or [Cloudflare](https://www.cloudflare.com/) or [Google Domains](https://domains.google) or [AliDNS](https://dc.console.aliyun.com).

* Or just register an account from [DuckDNS](https://www.duckdns.org/).

## Get it

### Build it from source code

* Get source code from Github:

```bash
git clone https://github.com/timothyye/godns.git
```

* Go into the godns directory, get related library and then build it:

```bash
cd cmd/godns
go get -v
go build
```

### Download from releases

Download compiled binaries from [releases](https://github.com/TimothyYe/godns/releases)

## Get help

```bash
$ ./godns -h
Usage of ./godns:
  -c string
        Specify a config file (default "./config.json")
  -h    Show help
```

## Config it

* Get [config_sample.json](https://github.com/timothyye/godns/blob/master/config_sample.json) from Github.
* Rename it to **config.json**.
* Configure your provider, domain/subdomain info, username and password, etc.
* Configure the SMTP options if you want, a mail notification will sent to your mailbox once the IP is changed.
* Save it in the same directory of GoDNS, or use -c=your_conf_path command.

## Config fields

* provider: The providers that GoDNS supports, available values are: `Cloudflare`, `Google`, `DNSPod`, `AliDNS`, `HE`, `DuckDNS`, `Dreamhost`.
* email: Email or account name of your DNS provider.
* password: Password of your account.
* login_token: API token of your account.
* domains: Domains list, with your sub domains.
* ip_url: A site helps you to get your public IPv4 IP address.
* ipv6_url: A site helps you to get your public IPv6 address.
* ip_type: To configure GoDNS under IPv4 mode or IPv6 mode, available values are: `IPv4`, `IPv6`.
* interval: The interval `seconds` that GoDNS check your public IP.
* socks5_proxy: Socks5 proxy server.
* resolver: The address of the public DNS server. For example, to run GoDNS in `IPv4` mode, you can set resolver as `8.8.8.8`, to GoDNS in `IPv6` mode, you can set resolver as `2001:4860:4860::8888`.

## IPv6 support

Supported provider(s):
* Cloudflare
* DNSPod
* Dreamhost
* DuckDNS
* Google Domains
* HE.net
* No-IP

To enable the `IPv6` support of GoDNS, there are 2 solutions you can choose:
* Get IPv6 address online
* Get IPv6 address from your network interface

To get IPv6 address online:  
* Set the `ip_type` as `IPv6`, and make sure the `ipv6_url` is configured.
* Add one `AAAA` record to your provider.

For example:

```json
{
  "domains": [
    {
      "domain_name": "example.com",
      "sub_domains": [
        "ipv6"
      ]
    }
  ],
  "resolver": "2001:4860:4860::8888",
  "ipv6_url": "https://api-ipv6.ip.sb/ip",
  "ip_type": "IPv6"
}
```

To get IPv6 address from your interface:

It is quite simple, just left "ip_url" & "ipv6_url" as empty. Please note that an IPv6 address configured for your interface is the prerequisite for this feature. 

### Config example for Cloudflare

For Cloudflare, you need to provide the email & Global API Key as password (or to use the API token) and config all the domains & subdomains.

* Using email & Global API Key

```json
{
  "provider": "Cloudflare",
  "email": "you@example.com",
  "password": "Global API Key",
  "domains": [{
      "domain_name": "example.com",
      "sub_domains": ["www","test"]
    },{
      "domain_name": "example2.com",
      "sub_domains": ["www","test"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_url": "https://myip.biturl.top",
  "interval": 300,
  "socks5_proxy": ""
}
```

* Using the API Token

```json
{
  "provider": "Cloudflare",
  "login_token": "API Token",
  "domains": [{
      "domain_name": "example.com",
      "sub_domains": ["www","test"]
    },{
      "domain_name": "example2.com",
      "sub_domains": ["www","test"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_url": "https://myip.biturl.top",
  "interval": 300,
  "socks5_proxy": ""
}
```

### Config example for DNSPod

For DNSPod, you need to provide your API Token(you can create it [here](https://www.dnspod.cn/console/user/security)), and config all the domains & subdomains.

```json
{
  "provider": "DNSPod",
  "login_token": "your_id,your_token",
  "domains": [{
      "domain_name": "example.com",
      "sub_domains": ["www","test"]
    },{
      "domain_name": "example2.com",
      "sub_domains": ["www","test"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_url": "https://myip.biturl.top",
  "ip_type": "IPV4",
  "interval": 300,
  "socks5_proxy": ""
}
```

### Config example for Dreamhost

For Dreamhost, you need to provide your API Token(you can create it [here](https://panel.dreamhost.com/?tree=home.api)), and config all the domains & subdomains.

```json
{
  "provider": "Dreamhost",
  "login_token": "your_api_key",
  "domains": [{
      "domain_name": "example.com",
      "sub_domains": ["www","test"]
    },{
      "domain_name": "example2.com",
      "sub_domains": ["www","test"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_url": "https://myip.biturl.top",
  "ip_type": "IPV4",
  "interval": 300,
  "resolver": "ns1.dreamhost.com",
  "socks5_proxy": ""
}
```

### Config example for Google Domains

For Google Domains, you need to provide email & password, and config all the domains & subdomains.

```json
{
  "provider": "Google",
  "email": "Your_Username",
  "password": "Your_Password",
  "domains": [{
      "domain_name": "example.com",
      "sub_domains": ["www","test"]
    },{
      "domain_name": "example2.com",
      "sub_domains": ["www","test"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_url": "https://myip.biturl.top",
  "interval": 300,
  "socks5_proxy": ""
}
```

### Config example for AliDNS

For AliDNS, you need to provide `AccessKeyID` & `AccessKeySecret` as `email` & `password`,  and config all the domains & subdomains.

```json
{
  "provider": "AliDNS",
  "email": "AccessKeyID",
  "password": "AccessKeySecret",
  "login_token": "",
  "domains": [{
      "domain_name": "example.com",
      "sub_domains": ["www","test"]
    },{
      "domain_name": "example2.com",
      "sub_domains": ["www","test"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_url": "https://myip.biturl.top",
  "interval": 300,
  "socks5_proxy": ""
}
```

### Config example for DuckDNS

For DuckDNS, only need to provide the `token`, config 1 default domain & subdomains.

```json
{
  "provider": "DuckDNS",
  "password": "",
  "login_token": "3aaaaaaaa-f411-4198-a5dc-8381cac61b87",
  "domains": [
    {
      "domain_name": "www.duckdns.org",
      "sub_domains": [
        "myname"
      ]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_url": "https://myip.biturl.top",
  "interval": 300,
  "socks5_proxy": ""
}
```

### Config example for No-IP 

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
  "ip_url": "https://myip.biturl.top",
  "resolver": "8.8.8.8",
  "interval": 300,
  "socks5_proxy": ""
}
```

### Config example for HE.net

For HE, email is not needed, just fill DDNS key to password, and config all the domains & subdomains.

```json
{
  "provider": "HE",
  "password": "YourPassword",
  "login_token": "",
  "domains": [{
      "domain_name": "example.com",
      "sub_domains": ["www","test"]
    },{
      "domain_name": "example2.com",
      "sub_domains": ["www","test"]
    }
  ],
  "resolver": "8.8.8.8",
  "ip_url": "https://myip.biturl.top",
  "interval": 300,
  "socks5_proxy": ""
}
```

### HE.net DDNS configuration

Add a new "A record", make sure that "Enable entry for dynamic dns" is checked:

<img src="https://github.com/TimothyYe/godns/blob/master/snapshots/he1.png?raw=true" width="640" />

Fill your own DDNS key or generate a random DDNS key for this new created "A record":

<img src="https://github.com/TimothyYe/godns/blob/master/snapshots/he2.png?raw=true" width="640" />

Remember the DDNS key and fill it as password to the config.json.

__NOTICE__: If you have multiple domains or subdomains, make sure their DDNS key are the same.

### Get an IP address from the interface

For some reasons if you want to get an IP directly from the interface, say `eth0` for Linux or `Local Area Connection` for Windows, update config file like this:

```json
  "ip_url": "",
  "ip_interface": "eth0",
```

If you set both `ip_url` and `ip_interface`, it first tries to get an IP address online, and if not succeed, gets
an IP address from the interface as a fallback.

### Email notification support

Update config file and provide your SMTP options, a notification mail will be sent to your mailbox once the IP is changed and updated.  

```json
  "notify": {
    "mail": {
      "enabled": true,
      "smtp_server": "smtp.example.com",
      "smtp_username": "user",
      "smtp_password": "password",
      "smtp_port": 25,
      "send_to": "my_mail@example.com"
    }
  }
```

Notification mail example:  

<img src="https://github.com/TimothyYe/godns/blob/master/snapshots/mail.png?raw=true" />  

### Telegram notification support

Update config file and provide your Telegram options, a notification message will be sent to your telegram channel once the IP is changed and updated.  

```json
  "notify": {
    "telegram": {
      "enabled": true,
      "bot_api_key": "11111:aaaa-bbbb",
      "chat_id": "-123456",
      "message_template": "Domain *{{ .Domain }}* is updated to %0A{{ .CurrentIP }}",
      "use_proxy": false
    },
  }
```
Markdown is supported in message template, and use `%0A` for newline.

### Slack notification support

Update config file and provide your Slack options, a notification message will be sent to your slack channel once the IP is changed and updated.  

```json
  "notify": {
    "slack": {
      "enabled": true,
      "bot_api_token": "xoxb-xxx",
      "channel": "your_channel",
      "message_template": "Domain *{{ .Domain }}* is updated to \n{{ .CurrentIP }}",
      "use_proxy": false
    },
  }
```
Markdown is supported in message template, and use `\n` for newline.


### SOCKS5 proxy support

You can also use SOCKS5 proxy, just fill SOCKS5 address to the ```socks5_proxy``` item:

```json
"socks5_proxy": "127.0.0.1:7070"
"use_proxy": true
```

Now all the queries will go through the specified SOCKS5 proxy.

## Run it as a daemon manually

```bash
nohup ./godns &
```

## Run it as a daemon, manage it via Upstart

* Install `upstart` first
* Copy `./upstart/godns.conf` to `/etc/init`
* Start it as a system service:

```bash
sudo start godns
```

## Run it as a daemon, manage it via Systemd

* Modify `./systemd/godns.service` and config it.
* Copy `./systemd/godns.service` to `/lib/systemd/system`
* Start it as a systemd service:

```bash
sudo systemctl enable godns
sudo systemctl start godns
```

## Run it with docker

Now godns supports to run in docker.

* Get [config_sample.json](https://github.com/timothyye/godns/blob/master/config_sample.json) from Github.
* Rename it to **config.json**.
* Run GoDNS with docker:

```bash
docker run -d --name godns --restart=always \
-v /path/to/config.json:/usr/local/godns/config.json timothyye/godns:latest
```

## Run it as a Windows service

After creating your config.json file:

* Get the latest [NSSM](https://nssm.cc/download) for create a windows service.
* Open the command prompt (as administrator) in the folder where you downloaded NSSM (e.g. C:\Downloads\nssm\ **win64**) and run:

```
nssm install YOURSERVICENAME
```

You will have an interface to configure your service, it is very simple in the "Application" tab just indicate where your `godns.exe` file is. Optionally you can also define a description on the "Details" tab and define a log file on the "I/O" tab.

* Finish using the "Install service" button.
* Done. Now whenever windows start your service will be loaded.
* To uninstall:

```
nssm remove YOURSERVICENAME
```

## Special Thanks
<img src="https://i.imgur.com/xhe5RLZ.jpg" width="80px" align="right" />

Thanks JetBrains for sponsoring this project with [free open source license](https://www.jetbrains.com/community/opensource/). 

> I like GoLand, it is an amazing and productive tool.
