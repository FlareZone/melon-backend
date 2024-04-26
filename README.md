# melon-backend
melon backend api

## 环境准备

## start
```shell
cp config.test.yaml.example ${your config path}/config.test.yaml
# edit your config.test.yaml file
go run main.go api -c ${your config path direction}/config.test.yaml
```

## 测试助记词

```cgo
check antique innocent spice much neglect split lottery trouble twelve report tennis

eth_address: 0xcBea7440fB18e72578c34AFA885f9Cc16be2f3d5
```


## 部署流程
1. 打开zeabur控制台
2. 使用prebuild（https://zeabur.com/docs/zh-CN/deploy/customize-prebuilt） 创建mysql和redis服务，将相关参数填写到config.yaml中。
3. 将配置文件config.yaml，通过（https://zeabur.com/docs/zh-CN/deploy/config-edit） 放到`/etc/melon/config.prod.yaml`路径下。
4. 点击重新自动部署即可。

## 本地测试
启动命令`go run main.go api -e dev`

配置文件config.dev.yaml示例
```yaml
app_name: melon
app_url: http://localhost:8080
log_level: INFO
database:
  melon:
    logging: true
    dsn: root:J2fOuZmhq816g7cKvPH53LACs4D9SX0M@tcp(hkg1.clusters.zeabur.com:30642)/zeabur?parseTime=true
redis:
  melon:
    addr: hkg1.clusters.zeabur.com:30494
    password:
    DB: 0
jwt:
  secret: f7&vh@VHuU0XpXZLm6wT2*HlYpPG#j0G3Tzsi!Qwsk#fY9+gF43g8Eu
  issuer: melon
#  （https://developers.google.com/identity/protocols/oauth2/web-server?hl=zh-cn）
oauth_v2: 
  google:
    client_id: 1030441591409-i0eiesff2uj64mhe3bl66338ofcv8sar.apps.googleusercontent.com
    client_secret: 
    redirect_url: http://localhost:8080/auth/google/callback
eip712:
  chain_id: 1
  version: 1
  name: melon
  verifying_Contract: 0xdAC17F958D2ee523a2206206994597C13D831ec7
oss:
  aliyun:
    endpoint: oss-cn-beijing.aliyuncs.com
    accessKeyId: LTAI5t9bPq1vCoj2qyUX1K1a
    accessKeySecret:
    bucketName: melon-save-test
    selfDomain:
mail:
  google:
#（邮箱应用专用密码，而不是邮箱密码。）（https://support.google.com/accounts/answer/6103523?sjid=7624534582211851258-AP）
    password:                      
    sender: lantianlaoli@gmail.com
smart_contract:
  proposal_logic_contract_address:
  flare_token_contract_address:
#    前往infura等节点提供商获取，（https://app.infura.io/）
  jsonrpc:
```