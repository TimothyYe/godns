```text
 ██████╗  ██████╗ ██████╗ ███╗   ██╗███████╗
██╔════╝ ██╔═══██╗██╔══██╗████╗  ██║██╔════╝
██║  ███╗██║   ██║██║  ██║██╔██╗ ██║███████╗
██║   ██║██║   ██║██║  ██║██║╚██╗██║╚════██║
╚██████╔╝╚██████╔╝██████╔╝██║ ╚████║███████║
 ╚═════╝  ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝╚══════╝

Latest release: V1.1
 ```

[![Build Status](https://travis-ci.org/TimothyYe/godns.svg?branch=master)](https://travis-ci.org/TimothyYe/godns)

GoDNS is a dynamic DNS (DDNS) tool, which is forked from [TimothyYe/godns](https://github.com/TimothyYe/godns). 

In this branch the support for mips32 is added, which means it could run properly on Openwrt and LEDE.

## Pre-condition

* GoDNS relies on [DNSPod](http://dnspod.cn) and its API. 

* To use GoDNS, you need a domain hosted on [DNSPod](http://dnspod.cn).

## Build it

So far, the latest version of Golang(v1.8) has not totally supported mips32. Openwrt and LEDE devices hardly enable the FPU emulator. Therefore, we still use the third-party compiler.

### Get & build go-mips32

* Git source code from GitHub:

```bash
git clone https://github.com/gomini/go-mips32.git
```

* Go into the go-mips32 directory, set the env and then build it:

```bash
export GOOS=linux
export GOARCH=mips32
CGO_ENABLED=0 ./make.bash
```

### Get & build godns from source code

* Get source code from Github:

```bash
git clone https://github.com/hguandl/godns.git
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
docker pull timothyye/godns:1.0
```

* Run godns in container and pass config parameters to it via enviroment variables:

```bash
docker run -d --name godns --restart=always \
-e EMAIL=your_dnspod_account \
-e PASSWORD=your_dnspod_password \
-e DOMAINS="your_domain1,your_domain2" DOCKER_IMAGE_ID                                                                                              
```



## Enjoy it!
