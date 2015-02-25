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
git clone https://github.com/abotoo/godns.git
```
* Go into the godns directory, get related library and then build it:

```bash
cd godns
go get
go build
```
* Then you get GoDNS.

## Run it

* Get [config_sample.json](https://github.com/abotoo/godns/blob/master/config_sample.json) from Github.
* Rename it to **config.json**.
* Configure your domain/sub-domain info, username and password of DNSPod account.
* Configure log file path, max size of log file, max count of log file.
* Configure user id, group id for safty.
* Save it in the same directory of GoDNS, or use -c=your_conf_path command.
* The last step, run godns:

```bash
./godns
```
* Enjoy it!
