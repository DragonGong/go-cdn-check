# 注意事项
``在调用前需要您进行操作的事项。``\
``以下项均在common\Client.go中``
## 一. 访问域名
设置您将要访问的域名（HOST）
```
        // {endPoint} --> 请求域名
        // 默认为 constant.HttpRequestDomain
        EndPoint: "{endPoint}"
        config.EndPoint = "{endPoint}"
        // ps: EndPoint: "open.chinanetcenter.com"
        // ps: config.EndPoint = "open.chinanetcenter.com"
```

## 二. 请求信息
设置 {objName}Request对象 各属性字段值
```
        {objName}Request.set{ParamName}(value);
```

## 三. 自定义请求信息
```
// 调用
// 可传自定义JSON字符串 替换`${objName}.String()`
response := auth.Invoke(config, ${objName}.String())
```

## 四. 账号信息
设置您的 AK,SK 信息
```
        // {accessKey} 您的AK
        config.AccessKey = "{accessKey}"
        // {secretKey} 您的SK
        config.SecretKey = "{secretKey}"
```

