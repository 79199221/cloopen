# Yuntongxun International SMS SDK for GO

[容联云通讯](https://www.yuntongxun.com) 短信发送（国际版）SDK

`edited by xiaozi@ixiaozi.cn`

## Quick Start

```go
go get -u github.com/79199221/cloopen/cloopen
```

```go
package main

import (
 "github.com/cloopen/79199221/cloopen/cloopen"
 "log"
)

func main() {
  cfg := cloopen.DefaultConfig().
  // 开发者主账号
  WithAPIAccount("xxxxxxxxxxxxx").
  // 主账号令牌 TOKEN
  WithAPIToken("xxxxxxxxxxxxx").
  // 应用的APPID
  WithAppId("xxxx")
  sms := cloopen.NewJsonClient(cfg).SMS()
  // 下发包体参数
  input := &cloopen.SendRequest{
    // 手机号码
    To: "00861352*******",
    // [国内短信的]模版ID
    TemplateId: "123456",
    // [国外短信的]模版内容
    Template: "[xxx]Your verification code is: {{code}}",
    // 模版变量内容 非必填，重写了这个字段的类型由list改为了map
    Datas: map[string]string{
      "code": "123456",
    },
 }
 // 下发
 resp, err := sms.Send(input)
 if err != nil {
  log.Fatal(err)
  return
 }
 log.Printf("Response MsgId: %s \n", resp.TemplateSMS.SmsMessageSid)

}

```

## 使用说明

* 自定义配置及默认

  `WithAPIAccount(xxx)` 配置主账号   **需调用者初始化此值**

  `WithAPIToken(xxx)` 配置主账号令牌  **需调用者初始化此值**

  `WithSmsHost(xxx)` 配置ip:port    **默认 app.cloopen.com:8883**

  `WithUseSSL(true)` 配置是否使用https  **默认启用https**

  `WithHTTPClient(customHttp)` 配置自定义httpClient  **默认使用sdk封装的httpClient**

  `WithHttpConf(&HttpConf{...})` 配置sdk封装的httpClient可调整参数 **默认使用sdk封装的httpClient参数**

  **参考HttpConf默认配置：**

  ```go
  // 时间单位为毫秒
  &HttpConf{
     Timeout:             300,
     KeepAlive:           30000,
     MaxIdleConns:        100,
     IdleConnTimeout:     30000,
     TLSHandshakeTimeout: 300,
  }
  ```

* 方法调用

  `cloopen.NewJsonClient(cfg)`  json 格式包体使用此方法

  `cloopen.NewXmlClient(cfg)`    xml  格式包体使用此方法

## 源码说明

- sdk
  - config.go 接口基础配置
  - client.go  客户端定义、配置
  - fields.go 常量定义
  - sms.go 短信功能
  - util.go 工具函数
- 分支说明
  - master最新稳定发布版本
  - develop待发布版本，贡献的代码请pull request到这里:)

## 请求错误

```json
{
    "statusCode":"000000",
    "statusMsg":"成功",
    "msgId":"0c26d83092d247b99076862f3ea24d4d",
    "failList":[
        {
            "mobile": "00861352*******",
            "errorCode": "162026",
            "errorMsg": "当前地区未开通权限"
        }
    ]
} 
```