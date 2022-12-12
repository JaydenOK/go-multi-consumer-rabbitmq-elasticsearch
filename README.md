# go-gin-multi-consumer-mysql-rabbitmq-elasticsearch  

## 说明
```text script
本项目流程及目的：gin框架，业务数据保存mysql后，发送rabbitmq，启动go协程消费队列数据，将数据推送到es，业务多条件查询时，
直接查es，解决大数据查询慢的问题。 gin框架，mvc架构（controller-service-model）
```

### -- 
 ```shell script
gin->redis->mysql->rabbitmq(kafka)->go consumer->elasticsearch
```

### 项目架构：  
```yaml script
app:
    config: app.yaml 项目核心配置
    constants: 常量配置
    controllers 控制器l
    events 事件
    libs 类库
    logs 运行日志
    models 模型
    routers 路由
    services 服务层
    utils 助手类库 
    tasks 消费者监听任务 

    main.go 程序入口文件
```

### 版本：
```shell script
go v1.19
gin v1.8.1
mysql v5.7
rabbitmq v3.6.1
elasticsearch v8.4.1
redis v3.2
#mongo 
```

### 启动
```shell script
系统配置如下：
app: map[httpport:8080 httpurl:127.0.0.1:8080 logfile:logs/app.log rpcport:9001 websocketport:8089 websocketurl:127.0.0.1:8089]
mysql: map[database:yb_new_hwc host:192.168.71.175 password:123456#Hsd1h port:3306 username:root]
redis: map[host:192.168.71.238 password:654321 port:6379]
mongo: map[database:hwc host:192.168.71.238 password:hwc@2018 port:27017 username:hwcuser]
rabbitmq: map[host:192.168.71.91 password:admin123. port:5672 username:admin]
elasticsearch: map[host:192.168.92.65 password: port:9200 username:]
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /order/list               --> app/controllers.(*OrderController).Lists-fm (3 handlers)
[GIN-debug] GET    /order/esList             --> app/controllers.(*OrderController).EsLists-fm (3 handlers)
[GIN-debug] POST   /order/add                --> app/controllers.(*OrderController).Add-fm (3 handlers)
[GIN-debug] POST   /order/update             --> app/controllers.(*OrderController).Update-fm (3 handlers)
[GIN-debug] POST   /order/delete             --> app/controllers.(*OrderController).Delete-fm (3 handlers)
[GIN-debug] POST   /user/register            --> app/controllers.(*UserController).Register-fm (3 handlers)
[GIN-debug] GET    /user/list                --> app/controllers.(*UserController).List-fm (3 handlers)
[GIN-debug] POST   /user/signIn              --> app/controllers.(*UserController).SignIn-fm (3 handlers)
[GIN-debug] POST   /user/signOut             --> app/controllers.(*UserController).SignOut-fm (3 handlers)
[GIN-debug] GET    /es/indexLists            --> app/controllers.(*EsController).IndexLists-fm (3 handlers)
[GIN-debug] GET    /es/indexExist            --> app/controllers.(*EsController).IndexExist-fm (3 handlers)
[GIN-debug] POST   /es/indexCreate           --> app/controllers.(*EsController).IndexCreate-fm (3 handlers)
[GIN-debug] GET    /es/indexGetMapping       --> app/controllers.(*EsController).IndexGetMapping-fm (3 handlers)
[GIN-debug] POST   /es/indexPutMapping       --> app/controllers.(*EsController).IndexPutMapping-fm (3 handlers)
[GIN-debug] POST   /es/indexReindex          --> app/controllers.(*EsController).IndexReindex-fm (3 handlers)
[GIN-debug] POST   /es/indexDelete           --> app/controllers.(*EsController).IndexDelete-fm (3 handlers)
[GIN-debug] GET    /es/indexAliasLists       --> app/controllers.(*EsController).IndexAliasLists-fm (3 handlers)
[GIN-debug] POST   /es/indexAlias            --> app/controllers.(*EsController).IndexAlias-fm (3 handlers)
[GIN-debug] GET    /consumer/startConsumer   --> app/controllers.(*ConsumerController).StartConsumer-fm (3 handlers)
[GIN-debug] GET    /consumer/stopConsumer    --> app/controllers.(*ConsumerController).StopConsumer-fm (3 handlers)
[GIN-debug] GET    /consumer/stopAll         --> app/controllers.(*ConsumerController).StopAll-fm (3 handlers)
2022/11/30 09:42:21 [200 OK] {
  "name" : "dcm_hk_getorder",
添加消费者： order_consumer
  "cluster_name" : "elasticsearch",
添加消费者： stock_consumer
  "cluster_uuid" : "zh7bym1uR8qMMJ1fRZVPgw",
  "version" : {
    "number" : "8.4.1",
    "build_flavor" : "default",
    "build_type" : "tar",
    "build_hash" : "2bd229c8e56650b42e40992322a76e7914258f0c",
    "build_date" : "2022-08-26T12:11:43.232597118Z",
    "build_snapshot" : false,
    "lucene_version" : "9.3.0",
    "minimum_wire_compatibility_version" : "7.17.0",
    "minimum_index_compatibility_version" : "7.0.0"
  },
  "tagline" : "You Know, for Search"
}
```

### 访问体验
```shell script
新增/更新/删除订单, 异步推送rabbitmq, 启动添加的消费者order_consumer消费,并推送elasticsearch: 

# 添加: 127.0.0.1:8080/order/add
{
    "order_id": "test0001",
    "platform_code": "AMAZON",
    "ship_name": "张珊",
    "total_price": 30,
    "order_status": "2",
    "middle_create_time": "2022-10-08 00:00:00"
}

# 更新: 127.0.0.1:8080/order/update
{
    "order_id": "test3333",
    "platform_code": "EB",
    "ship_name": "james",
    "ship_phone": "1888888835"
}

# 删除: 127.0.0.1:8080/order/delete
order_id: test3333



# 原来的mysql访问列表地址: 
127.0.0.1:8080/order/list?page=1&pageSize=20&order_id=asd&platform_code=EB&total_price_start=45&total_price_end=50  

# 优化后, 直接查elasticsearch接口: 
127.0.0.1:8080/order/esList?page=1&pageSize=20&platform_code=AMAZON,EB&total_price_start=45&total_price_end=50
```