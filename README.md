```text
 ██████╗  ██████╗ ██████╗ ███╗   ██╗███████╗
██╔════╝ ██╔═══██╗██╔══██╗████╗  ██║██╔════╝
██║  ███╗██║   ██║██║  ██║██╔██╗ ██║███████╗
██║   ██║██║   ██║██║  ██║██║╚██╗██║╚════██║
╚██████╔╝╚██████╔╝██████╔╝██║ ╚████║███████║
 ╚═════╝  ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝╚══════╝
 ```

[![Release][7]][8] [![MIT licensed][9]][10] [![Build Status][1]][2] [![Docker][3]][4] [![Go Report Card][11]][12] [![Cover.Run][15]][16] [![GoDoc][13]][14]

[1]: https://travis-ci.org/TimothyYe/godns.svg?branch=master
[2]: https://travis-ci.org/TimothyYe/godns
[3]: https://images.microbadger.com/badges/image/timothyye/godns.svg
[4]: https://microbadger.com/images/timothyye/godns
[7]: http://github-release-version.herokuapp.com/github/timothyye/godns/release.svg?style=flat
[8]: https://github.com/TimothyYe/godns/releases
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

## Supported DNS Provider
* DNSPod ([https://www.dnspod.cn/](https://www.dnspod.cn/))
* HE.net (Hurricane Electric) ([https://dns.he.net/](https://dns.he.net/))

## Supported Platforms
* Linux
* MacOS
* ARM Linux (Raspberry Pi, etc...)
* Windows

## MIPS32 platform

For MIPS32 platform, please checkout the [mips32](https://github.com/TimothyYe/godns/tree/mips32) branch, this branch is contributed by [hguandl](https://github.com/hguandl), in this branch, the support for mips32 is added, which means it could run properly on Openwrt and LEDE.

## Pre-condition

* Register and own a domain.

* Domain's nameserver points to [DNSPod](https://www.dnspod.cn/) or [HE.net](https://dns.he.net/).

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
  -d    Run it as docker mode
  -h    Show help
```

## Config it

* Get [config_sample.json](https://github.com/timothyye/godns/blob/master/config_sample.json) from Github.
* Rename it to **config.json**.
* Configure your provider, domain/sub-domain info, username and password, etc.
* Configure the SMTP options if you want, a mail notification will sent to your mailbox once the IP is changed.
* Save it in the same directory of GoDNS, or use -c=your_conf_path command.

### Config example for DNSPod

For DNSPod, you need to provide email & password,  and config all the domains & subdomains.

```json
{
  "provider": "DNSPod",
  "email": "example@gmail.com",
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
  "ip_url": "http://members.3322.org/dyndns/getip",
  "log_path": "./godns.log",
  "socks5_proxy": ""
}
```
### Config example for HE.net

For HE, email is not needed, just fill DDNS key to password, and config all the domains & subdomains.

```json
{
  "provider": "HE",
  "email": "",
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
  "ip_url": "http://members.3322.org/dyndns/getip",
  "log_path":"/users/timothy/workspace/src/godns/godns.log",
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

### Email notification support

Update config file and provide your SMTP options, a notification mail will be sent to your mailbox once the IP is changed and updated.  

```json
  "notify": {
    "enabled": true,
    "smtp_server": "smtp.example.com",
    "smtp_username": "user",
    "smtp_password": "password",
    "smtp_port": 25,
    "send_to": "my_mail@example.com"
  }
```

Notification mail example:  

<img src="https://github.com/TimothyYe/godns/blob/master/snapshots/mail.png?raw=true" />  

### SOCKS5 proxy support

You can also use SOCKS5 proxy, just fill SOCKS5 address to the ```socks5_proxy``` item:

```json
"socks5_proxy": "127.0.0.1:7070"
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

## Enjoy it!
