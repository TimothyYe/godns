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
git clone https://github.com/TimothyYe/godns.git
```
* Go into the godns directory, get related library and then build it:

```bash
cd godns
go get
go build
```
* Then you get GoDNS.

### Build godns from the 3rd party 

* Visit this URL provided by [GoBuild](http://gobuild.io/download/github.com/TimothyYe/godns).
* Select the platform you need.
* Input the build type and branch name.
* Build and download GoDNS.

## Run it

* Get [config_sample.json](https://github.com/TimothyYe/godns/blob/master/config_sample.json) from Github.
* Rename it to **config.json**.
* Configure your domain/sub-domain info, username and password of DNSPod account.
* Save it in the same directory of GoDNS.
* The last step, run godns:
```bash
nohup ./godns &
```
* Enjoy it!
