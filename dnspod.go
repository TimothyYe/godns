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
    "github.com/bitly/go-simplejson"
    "time"
    "os"
    "errors"
)

const (
    ApiServer = "https://dnsapi.cn" //必须使用https
    AuthorEmail = "jarrysix@gmail.com"
)

var (
    loginToken string
)

// 获取接口登陆Token
func getLoginToken() string {
    if (loginToken == "") {
        loginToken = fmt.Sprintf("%d,%s", Configuration.ApiId, Configuration.ApiToken)
    }
    return loginToken
}

func get_currentIP(url string) (string, error) {
    response, err := http.Get(url)

    if err != nil {
        log.Println("Cannot get IP...")
        return "", err
    }

    defer response.Body.Close()

    body, _ := ioutil.ReadAll(response.Body)
    return string(body), nil
}

// 附加公共API参数
func appendCommParams(params url.Values) url.Values {
    if params == nil {
        panic("params not be nil!")
    }
    //header.Add("login_email", Configuration.Email)
    //header.Add("login_password", Configuration.Password)
    params.Set("format", "json")
    params.Set("lang", "en")
    params.Set("error_on_empty", "no")
    params.Set("login_token", getLoginToken())
    //仅代理用户需要使用以下参数
    //params.Set("user_id","")
    return params
}

// 获取版本
func GetApiVersion() *Version {
    ver := map[string]interface{}{}
    b, err := apiPost("/Info.Version", url.Values{})
    if err == nil {
        json.Unmarshal(b, &ver)
    }
    if err != nil {
        panic(err)
    }

    stat := ver["status"].(map[string]interface{})
    code := stat["code"]
    msg := stat["message"]
    date, _ := time.Parse("2006-01-02 15:04:05", stat["created_at"].(string))

    log.Println(code, getLoginToken())
    if code != "1" {
        log.Println(fmt.Sprintf("[ Check][ Version] - code:%s,message:%s", code, msg))
        os.Exit(0)
    }

    return &Version{
        ApiVersion: msg.(string),
        ApiDate: date,
        ClientVersion:ClientVersion,
    }
}

func get_domain(name string) int64 {
    var ret int64
    values := url.Values{}
    values.Add("type", "all")
    values.Add("offset", "0")
    values.Add("length", "20")

    response, err := apiPost("/Domain.List", values)

    if err != nil {
        log.Println("Failed to get domain list...")
        return -1
    }

    sjson, parse_err := simplejson.NewJson([]byte(response))

    if parse_err != nil {
        log.Println(parse_err)
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

func getSubdomain(domain_id int64, name string) (string, string) {
    var ret, ip string
    value := url.Values{}
    value.Add("domain_id", strconv.FormatInt(domain_id, 10))
    value.Add("offset", "0")
    value.Add("length", "1")
    value.Add("sub_domain", name)

    response, err := apiPost("/Record.List", value)

    if err != nil {
        log.Println("Failed to get domain list")
        return "", ""
    }

    sjson, parse_err := simplejson.NewJson([]byte(response))

    if parse_err != nil {
        log.Println(parse_err)
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

// 更新IP记录
func UpdateIpRecord(domainId int64, subDomainId string, subDomainName string, ip string) error {
    value := url.Values{}
    value.Add("domain_id", strconv.FormatInt(domainId, 10))
    value.Add("record_id", subDomainId)
    value.Add("sub_domain", subDomainName)
    value.Add("record_type", "A")
    value.Add("record_line", "默认")
    value.Add("value", ip)

    response, err := apiPost("/Record.Modify", value)

    if err == nil {
        sjson, _ := simplejson.NewJson([]byte(response))
        code := sjson.Get("status").Get("code").MustString()
        // 更新地址成功
        if code == "1" {
            return nil
        }
        return errors.New(fmt.Sprintf("record update fail!(code:%s)", code))
    }
    return err
}

// 提交数据到接口
func apiPost(url string, content url.Values) ([]byte, error) {
    client := http.DefaultClient
    values := appendCommParams(content)
    req, _ := http.NewRequest("POST", ApiServer + url, strings.NewReader(values.Encode()))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("User-Agent", fmt.Sprintf("GoDNS/0.1 (%s)", AuthorEmail))
    response, err := client.Do(req)
    defer response.Body.Close()

    if err != nil {
        log.Println("Post failed...", err.Error())
        return nil, err
    }
    return ioutil.ReadAll(response.Body)
}
