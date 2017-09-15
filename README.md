```text
 ██████╗  ██████╗ ██████╗ ███╗   ██╗███████╗
██╔════╝ ██╔═══██╗██╔══██╗████╗  ██║██╔════╝
██║  ███╗██║   ██║██║  ██║██╔██╗ ██║███████╗
██║   ██║██║   ██║██║  ██║██║╚██╗██║╚════██║
╚██████╔╝╚██████╔╝██████╔╝██║ ╚████║███████║
 ╚═════╝  ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝╚══════╝
 ```

[![Release][7]][8] [![MIT licensed][9]][10] [![Build Status][1]][2] [![Downloads][5]][6] [![Docker][3]][4] [![Go Report Card][11]][12]

[1]: https://travis-ci.org/TimothyYe/godns.svg?branch=master
[2]: https://travis-ci.org/TimothyYe/godns
[3]: https://images.microbadger.com/badges/image/timothyye/godns.svg
[4]: https://microbadger.com/images/timothyye/godns
[5]: https://img.shields.io/badge/downloads-1.95MB-brightgreen.svg
[6]: https://github.com/TimothyYe/godns/releases
[7]: https://img.shields.io/badge/release-v1.2-brightgreen.svg
[8]: https://github.com/TimothyYe/godns/releases
[9]: https://img.shields.io/badge/license-Apache-blue.svg
[10]: LICENSE
[11]: https://goreportcard.com/badge/github.com/timothyye/godns
[12]: https://goreportcard.com/report/github.com/timothyye/godns

GoDNS is a dynamic DNS (DDNS) tool, it is based on my early open source project: [DynDNS](https://github.com/TimothyYe/DynDNS). 

Now I rewrite [DynDNS](https://github.com/TimothyYe/DynDNS) by Golang and call it [GoDNS](https://github.com/TimothyYe/godns).

## MIPS32 platform

For MIPS32 platform, please checkout the [mips32](https://github.com/TimothyYe/godns/tree/mips32) branch, this branch is contributed by [hguandl](https://github.com/hguandl), in this branch, the support for mips32 is added, which means it could run properly on Openwrt and LEDE.

## Pre-condition

* GoDNS relies on [DNSPod](http://dnspod.cn) and its API. 

* To use GoDNS, you need a domain hosted on [DNSPod](http://dnspod.cn).

## Build it

### Get & build it from source code

* Get source code from Github:

```bash
git clone https://github.com/timothyye/godns.git
```
* Go into the godns directory, get related library and then build it:

```bash
cd godns
go get
go build
```

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
* Configure your domain/sub-domain info, username and password of DNSPod account.
* Configure log file path, max size of log file, max count of log file.
* Configure user id, group id for safety.
* Save it in the same directory of GoDNS, or use -c=your_conf_path command.

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

## Run it in docker

Now godns supports to run in docker.

* Pull godns image from docker hub:
```bash
docker pull timothyye/godns:1.2
```

* Run godns in container and pass config parameters to it via enviroment variables:

```bash
docker run -d --name godns --restart=always \
-v /path/to/config.json:/usr/local/godns/config.json timothyye/godns:1.2
```

## Enjoy it!
