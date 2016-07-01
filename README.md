## GoDNS

[![Build Status](https://travis-ci.org/TimothyYe/godns.svg?branch=master)](https://travis-ci.org/TimothyYe/godns)

GoDNS is a dynamic DNS (DDNS) tool, it is based on my early open source project: [DynDNS](https://github.com/TimothyYe/DynDNS). 

Now I rewrite [DynDNS](https://github.com/TimothyYe/DynDNS) by Golang and call it [GoDNS](https://github.com/TimothyYe/godns).

## Pre-condition

* GoDNS relies on [DNSPod](http://dnspod.cn) and its API. 

* To use GoDNS, you need a domain and hosted on [DNSPod](http://dnspod.cn).

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

## Run it as a daemon, manage it via Upstart

Now godns supports to run in docker.

* Pull godns image from docker hub:
```bash
docker pull timothyye/godns:1.0
```

* Run it in container and pass config parameters to it via enviroment variables:

```bash
docker run -d --name godns --restart=always \
-e EMAIL=your_dnspod_account \
-e PASSWORD=your_dnspod_password \
-e DOMAIN=your_domain \
-e SUB_DOMAIN=your_sub_domain DOCKER_IMAGE_ID                                                                                              
```



## Enjoy it!