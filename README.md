## GoDNS

GoDNS is a dynamic DNS (DDNS) tool
## Pre-condition

* GoDNS relies on [DNSPod](http://dnspod.cn) and its API. 

* To use GoDNS, you need a domain and hosted on [DNSPod](http://dnspod.cn).

## Build it

### Get & build it from source code

* Get source code from Github:

```bash
git clone https://github.com/jsix/godns.git
```
* Go into the godns directory, get related library and then build it:

```bash
cd godns
go get
go build
```

## Run it
 Configure api-id ,api-token and domain/sub-domains of DNSPod account.

```bash
./godns
```
You can use your config file like this:
```bash
./godns -conf=your.conf
```
* Enjoy it!
