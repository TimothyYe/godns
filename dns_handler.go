package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"golang.org/x/net/proxy"
	"github.com/bitly/go-simplejson"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
)

var chain=`-----BEGIN CERTIFICATE-----
MIIIPzCCByegAwIBAgIQEnH9KIzQjspNOOh1aFgB3zANBgkqhkiG9w0BAQsFADBH
MQswCQYDVQQGEwJVUzEWMBQGA1UEChMNR2VvVHJ1c3QgSW5jLjEgMB4GA1UEAxMX
R2VvVHJ1c3QgRVYgU1NMIENBIC0gRzQwHhcNMTYwNTEyMDAwMDAwWhcNMTgwNzEw
MjM1OTU5WjCB1jETMBEGCysGAQQBgjc8AgEDEwJDTjEZMBcGCysGAQQBgjc8AgEC
EwhTaGFuRG9uZzEdMBsGA1UEDxMUUHJpdmF0ZSBPcmdhbml6YXRpb24xGDAWBgNV
BAUTDzM3MDYzNTIwMDAxMzgxNDELMAkGA1UEBhMCQ04xETAPBgNVBAgMCFNoYW5E
b25nMQ8wDQYDVQQHDAZZYW50YWkxFTATBgNVBAoMDEROU1BvZCwgSW5jLjELMAkG
A1UECwwCSVQxFjAUBgNVBAMMDXd3dy5kbnNwb2QuY24wggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQCDj1tJPjOu4XILNC8rDNN4DlMn6T1IQO3Bd9/xKc1Q
Uj/GF9jA52KO37l4uX7LOo2XeNsMKy2utGU226ToWAHB9mDnXj95tvRwPvMunDRN
QcVe/NS14v8XNpXo79tT9mnj/nZVAk/mHqlVMygPUyGDgvuzePKpdOrfYOomcWyO
bsl3X/5EaEkQGuF1tx+DUQSYmhtB1FCM/uzUdd8sZav9d5ADPm4weRI7au054oWq
xdPhUngf2oTE/7lqsFRwrQ8DXK/asMx7jIAqImBLjSUOzQ6bsVZzEv+rAAHHEUNV
ovx+tvMjKR9MyMcZigAK+AxiANByR2p6t9t/BIEe8JgHAgMBAAGjggSVMIIEkTCC
AX4GA1UdEQSCAXUwggFxgg5saWJzLmRuc3BvZC5jboIRbW9uaXRvci5kbnNwb2Qu
Y26CFnNzbC5wdGxvZ2luMi5kbnNwb2QuY26CEXN0YXRpYy5kbnNwb2QuY29tggtt
LmRuc3BvZC5jboIOYXBpLmRuc3BvZC5jb22CDXd3dy5kbnNhcGkuY26CEm1vbml0
b3IuZG5zcG9kLmNvbYISdGlja2V0cy5kbnNwb2QuY29tgg5zdGF0LmRuc3BvZC5j
boIRZG9tYWlucy5kbnNwb2QuY26CDnd3dy5kbnNwb2QuY29tghF0aWNrZXRzLmRu
c3BvZC5jboIKZG5zcG9kLmNvbYIJZG5zYXBpLmNughFzdXBwb3J0LmRuc3BvZC5j
boIRc3RhdGljcy5kbnNwb2QuY26CEnN1cHBvcnQuZG5zcG9kLmNvbYIMZWMuZG5z
cG9kLmNugg5ibG9nLmRuc3BvZC5jboINd3d3LmRuc3BvZC5jboIJZG5zcG9kLmNu
MAkGA1UdEwQCMAAwDgYDVR0PAQH/BAQDAgWgMCsGA1UdHwQkMCIwIKAeoByGGmh0
dHA6Ly9nbS5zeW1jYi5jb20vZ20uY3JsMIGpBgNVHSAEgaEwgZ4wgZIGCSsGAQQB
8CIBBjCBhDA/BggrBgEFBQcCARYzaHR0cHM6Ly93d3cuZ2VvdHJ1c3QuY29tL3Jl
c291cmNlcy9yZXBvc2l0b3J5L2xlZ2FsMEEGCCsGAQUFBwICMDUMM2h0dHBzOi8v
d3d3Lmdlb3RydXN0LmNvbS9yZXNvdXJjZXMvcmVwb3NpdG9yeS9sZWdhbDAHBgVn
gQwBATAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwHwYDVR0jBBgwFoAU
3s9cULeuAh8VF6oW6A21KJ1qWvMwVwYIKwYBBQUHAQEESzBJMB8GCCsGAQUFBzAB
hhNodHRwOi8vZ20uc3ltY2QuY29tMCYGCCsGAQUFBzAChhpodHRwOi8vZ20uc3lt
Y2IuY29tL2dtLmNydDCCAX4GCisGAQQB1nkCBAIEggFuBIIBagFoAHcA3esdK3oN
T6Ygi4GtgWhwfi6OnQHVXIiNPRHEzbbsvswAAAFUpDOIuQAABAMASDBGAiEA9/1k
K0YQoivMrnAo2gP2Kz+sSm1u96VedFZazmRMcQICIQDJaWZwp/U9iadBgMcvFKRK
NLGo1+ViRw+k5O8LkVx1lwB2AKS5CZC0GFgUh7sTosxncAo8NZgE+RvfuON3zQ7I
DdwQAAABVKQxuEgAAAQDAEcwRQIgTX86Q1O7ACRJiEubi9hut+8C68+vYZR6vv4V
6dOCL1gCIQDvc/bgy9THl+2jzAf0xrXELz7gvq2CeaI+UFKWEKYqygB1AGj2mPgf
ZIK+OozuuSgdTPxxUV1nk9RE0QpnrLtPT/vEAAABVKQxuBsAAAQDAEYwRAIgV7Sf
V9iEuSz6OFvmsPTf+T6gv/6meL5FHdKhYkSDj/cCICVU99Uo9I/+zvm1BiBintsN
+dc6/N2GVk0lndEBVaL9MA0GCSqGSIb3DQEBCwUAA4IBAQDDUHx4F/jP6rP6jtpj
ADMsSb6x+UDJF/BeH0SsOEGRsFGASwlSX8uSb0+73gb2nBRziTzsGLiuuPDx8k0G
GQBpXCWXU7KVjluuX5a9C2j43swfPyWRC60rz6Dg+qVVaKocGv+r/Rqi2YcMN55N
CLdSsB9f3b0CWpJrmx9YVofaM5sNgal7cKJZ+W4o1gIsXuw97QCxj73HAmfMj/lF
fVxTgwx5l0v4RX1wAIWiPmrBDpy8okXIYUEbPJGAEhip7gD0rKZPw20uB3+raAPv
j1B+zK2wAfWkAPThX5Lwf7HxloFSiRvRQp334U0F+zcTemV+W1M49Y/37yBW8YV2
UNUv
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIEbjCCA1agAwIBAgIQboqQ68/wRIpyDQgF0IKlRDANBgkqhkiG9w0BAQsFADBY
MQswCQYDVQQGEwJVUzEWMBQGA1UEChMNR2VvVHJ1c3QgSW5jLjExMC8GA1UEAxMo
R2VvVHJ1c3QgUHJpbWFyeSBDZXJ0aWZpY2F0aW9uIEF1dGhvcml0eTAeFw0xMzEw
MzEwMDAwMDBaFw0yMzEwMzAyMzU5NTlaMEcxCzAJBgNVBAYTAlVTMRYwFAYDVQQK
Ew1HZW9UcnVzdCBJbmMuMSAwHgYDVQQDExdHZW9UcnVzdCBFViBTU0wgQ0EgLSBH
NDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANm0BfI4Zw8J53z1Yyrl
uV6oEa51cdlMhGetiV38KD0qsKXV1OYwCoTU5BjLhTfFRnHrHHtp22VpjDAFPgfh
bzzBC2HmOET8vIwvTnVX9ZaZfD6HHw+QS3DDPzlFOzpry7t7QFTRi0uhctIE6eBy
GpMRei/xq52cmFiuLOp3Xy8uh6+4a+Pi4j/WPeCWRN8RVWNSL/QmeMQPIE0KwGhw
FYY47rd2iKsYj081HtSMydt+PUTUNozBN7VZW4f56fHUxSi9HdzMlnLReqGnILW4
r/hupWB7K40f7vQr1mnNr8qAWCnoTAAgikkKbo6MqNEAEoS2xeKVosA7pGvwgtCW
XSUCAwEAAaOCAUMwggE/MBIGA1UdEwEB/wQIMAYBAf8CAQAwDgYDVR0PAQH/BAQD
AgEGMC8GCCsGAQUFBwEBBCMwITAfBggrBgEFBQcwAYYTaHR0cDovL2cyLnN5bWNi
LmNvbTBHBgNVHSAEQDA+MDwGBFUdIAAwNDAyBggrBgEFBQcCARYmaHR0cHM6Ly93
d3cuZ2VvdHJ1c3QuY29tL3Jlc291cmNlcy9jcHMwNAYDVR0fBC0wKzApoCegJYYj
aHR0cDovL2cxLnN5bWNiLmNvbS9HZW9UcnVzdFBDQS5jcmwwKQYDVR0RBCIwIKQe
MBwxGjAYBgNVBAMTEVN5bWFudGVjUEtJLTEtNTM4MB0GA1UdDgQWBBTez1xQt64C
HxUXqhboDbUonWpa8zAfBgNVHSMEGDAWgBQs1VBBlxWL8I82YVtK+2vZmckzkjAN
BgkqhkiG9w0BAQsFAAOCAQEAtI69B7mahew7Z70HYGHmhNHU7+sbuguCS5VktmZT
I723hN3ke40J2s+y9fHDv4eEvk6mqMLnEjkoNOCkVkRADJ+IoxXT6NNe4xwEYPtp
Nk9qfgwqKMHzqlgObM4dB8NKwJyNw3SxroLwGuH5Tim9Rt63Hfl929kPhMuSRcwc
sxj2oM9xbwwum9Its5mTg0SsFaqbLmfsT4hpBVZ7i7JDqTpsHBMzJRv9qMhXAvsc
4NG9O1ZEZcNj9Rvv7DDZ424uE+k5CCoMcvOazPYnKYTT70zHhBFlH8bjgQPbh8x4
97Wdlj5qf7wRhXp15kF9Dc/55YVpJY/HjQct+GkPy0FTAA==
-----END CERTIFICATE-----`

func decodePem(certInput string) tls.Certificate {
	var cert tls.Certificate
	certPEMBlock := []byte(certInput)
	var certDERBlock *pem.Block
	for {
		certDERBlock, certPEMBlock = pem.Decode(certPEMBlock)
		if certDERBlock == nil {
			break
		}
		if certDERBlock.Type == "CERTIFICATE" {
			cert.Certificate = append(cert.Certificate, certDERBlock.Bytes)
		}
	}
	return cert
}

func getCurrentIP(url string) (string, error) {
	response, err := http.Get(url)

	if err != nil {
		log.Println("Cannot get IP...")
		return "", err
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	return string(body), nil
}

func generateHeader(content url.Values) url.Values {
	header := url.Values{}
	if configuration.LoginToken != "" {
		header.Add("login_token", configuration.LoginToken)
	} else {
		header.Add("login_email", configuration.Email)
		header.Add("login_password", configuration.Password)
	}
	header.Add("format", "json")
	header.Add("lang", "en")
	header.Add("error_on_empty", "no")

	if content != nil {
		for k, _ := range content {
			header.Add(k, content.Get(k))
		}
	}

	return header
}

func apiVersion() {
	postData("/Info.Version", nil)
}

func getDomain(name string) int64 {

	var ret int64
	values := url.Values{}
	values.Add("type", "all")
	values.Add("offset", "0")
	values.Add("length", "20")

	response, err := postData("/Domain.List", values)

	if err != nil {
		log.Println("Failed to get domain list...")
		return -1
	}

	sjson, parseErr := simplejson.NewJson([]byte(response))

	if parseErr != nil {
		log.Println(parseErr)
		return -1
	}

	if sjson.Get("status").Get("code").MustString() == "1" {
		domains, _ := sjson.Get("domains").Array()

		for _, d := range domains {
			m := d.(map[string]interface{})
			if m["name"] == name {
				id := m["id"]

				switch t := id.(type) {
				case json.Number:
					ret, _ = t.Int64()
				}

				break
			}
		}
		if len(domains) == 0 {
			log.Println("domains slice is empty.")
		}
	} else {
		log.Println("get_domain:status code:", sjson.Get("status").Get("code").MustString())
	}

	return ret
}

func getSubDomain(domainID int64, name string) (string, string) {
	log.Println("debug:", domainID, name)
	var ret, ip string
	value := url.Values{}
	value.Add("domain_id", strconv.FormatInt(domainID, 10))
	value.Add("offset", "0")
	value.Add("length", "1")
	value.Add("sub_domain", name)

	response, err := postData("/Record.List", value)

	if err != nil {
		log.Println("Failed to get domain list")
		return "", ""
	}

	sjson, parseErr := simplejson.NewJson([]byte(response))

	if parseErr != nil {
		log.Println(parseErr)
		return "", ""
	}

	if sjson.Get("status").Get("code").MustString() == "1" {
		records, _ := sjson.Get("records").Array()

		for _, d := range records {
			m := d.(map[string]interface{})
			if m["name"] == name {
				ret = m["id"].(string)
				ip = m["value"].(string)
				break
			}
		}
		if len(records) == 0 {
			log.Println("records slice is empty.")
		}
	} else {
		log.Println("get_subdomain:status code:", sjson.Get("status").Get("code").MustString())
	}

	return ret, ip
}

func updateIP(domainID int64, subDomainID string, subDomainName string, ip string) {
	value := url.Values{}
	value.Add("domain_id", strconv.FormatInt(domainID, 10))
	value.Add("record_id", subDomainID)
	value.Add("sub_domain", subDomainName)
	value.Add("record_type", "A")
	value.Add("record_line", "默认")
	value.Add("value", ip)

	response, err := postData("/Record.Modify", value)

	if err != nil {
		log.Println("Failed to update record to new IP!")
		log.Println(err)
		return
	}

	sjson, parseErr := simplejson.NewJson([]byte(response))

	if parseErr != nil {
		log.Println(parseErr)
		return
	}

	if sjson.Get("status").Get("code").MustString() == "1" {
		log.Println("New IP updated!")
	}

}

func postData(url string, content url.Values) (string, error) {
	certChain := decodePem(chain)
	conf := tls.Config { }
	conf.RootCAs = x509.NewCertPool()
	for _, cert := range certChain.Certificate {
		x509Cert, err := x509.ParseCertificate(cert)
		if err != nil {
			panic(err)
		}
		conf.RootCAs.AddCert(x509Cert)
	}
	conf.BuildNameToCertificate()

	tr := http.Transport{ TLSClientConfig: &conf }
	client := &http.Client{Transport: &tr}

	if configuration.Socks5Proxy != "" {

		log.Println("use socks5 proxy:" + configuration.Socks5Proxy)

		dialer, err := proxy.SOCKS5("tcp", configuration.Socks5Proxy, nil, proxy.Direct)
		if err != nil {
			fmt.Println("can't connect to the proxy:", err)
			return "", err
		}

		httpTransport := &http.Transport{}
		client.Transport = httpTransport
		httpTransport.Dial = dialer.Dial
	}

	values := generateHeader(content)
	req, _ := http.NewRequest("POST", "https://dnsapi.cn" + url, strings.NewReader(values.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", fmt.Sprintf("GoDNS/0.1 (%s)", configuration.Email))

	response, err := client.Do(req)

	if err != nil {
		log.Println("Post failed...")
		log.Println(err)
		return "", err
	}

	defer response.Body.Close()
	resp, _ := ioutil.ReadAll(response.Body)

	return string(resp), nil
}
