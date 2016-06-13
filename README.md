# ssdbproxy

ssdb proxy route service

## 功能描述

SSDB路由代理服务：采用SSDB协议，调用像直连SSDB Server一样简单，无需其他改动。根据name或key的前缀进行匹配，将请求路由到不同的SSDB中。若把name和key看做是库和表，则相当于分库分表服务。

## 配置文件

采用JSON格式的配置

### 配置示例

```json
{
    "host":"0.0.0.0",
    "port":8888,
    "reloadSecond":60,
    "args":{
        "getClientTimeout":5,
        "maxPoolSize":50,
        "minPoolSize":5,
        "acquireIncrement":5,
        "maxIdleTime":120,
        "maxWaitSize":1200,
        "healthSecond":300
    },
    "master": {
        "id":"Default",
        "route":"",
        "host":"192.168.1.20",
        "port":8888,
        "password":"",
        "open":true,
        "ownArgs":false
    },
    "nodes":[
        {
            "id":"Test1",
            "route":"t1",
            "host":"192.168.1.20",
            "port":8876,
            "password":"",
            "open":true,
            "ownArgs":false
        },
        {
            "id":"Test2",
            "route":"t2",
            "host":"192.168.1.20",
            "port":8877,
            "password":"",
            "open":true,
            "ownArgs":true,
            "args":{
                "getClientTimeout":5,
                "maxPoolSize":50,
                "minPoolSize":5,
                "acquireIncrement":5,
                "maxIdleTime":120,
                "maxWaitSize":1200,
                "healthSecond":300
            }
        }
    ]
}
```

### 配置说明

| Fields | 描述 |
| ---- | ---- |
| host | 监听的IP或主机名 |
| port | 监听的端口 |
| reloadSecond | 重载配置文件时间间隔（单位：秒） |
| args | SSDB连接池默认参数 |
| master | SSDB主节点（匹配不到其他节点时，路由到主节点） |
| nodes | 路由节点列表 |

### 节点配置说明

| Fields | 描述 |
| ---- | ---- |
| id | 节点名称 |
| route | 路由规则（name或key的前缀） |
| host | SSDB 的IP或主机名 |
| port | SSDB 的端口 |
| password | SSDB 认证密码 |
| open | 节点是否开启 true/false （不开启相当于没有此节点） |
| ownArgs | 是否拥有自定义连接池参数 true/false （若没有则使用默认连接池参数） |
| args | 连接池参数 （当ownArgs为true时配置此项） |

### 连接池参数说明

| Fields | 描述 |
| ---- | ---- |
| getClientTimeout | 获取连接超时时间，单位为秒。 |
| maxPoolSize | 最大连接池个数。 |
| minPoolSize | 最小连接池数。 |
| acquireIncrement | 当连接池中的连接耗尽的时候一次同时获取的连接数。 |
| maxIdleTime | 最大空闲时间，指定秒内未使用则连接被丢弃。若为0则永不丢弃。 |
| maxWaitSize | 最大等待数目，当连接池满后，新建连接将等待池中连接释放后才可以继续，本值限制最大等待的数量，超过本值后将抛出异常。 |
| healthSecond | 健康检查时间隔，单位为秒。 |
