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

### 访问体验
```shell script
原来的mysql访问列表地址: 
127.0.0.1:8080/order/list?page=1&pageSize=20&order_id=asd&platform_code=EB&total_price_start=45&total_price_end=50  

elasticsearch优化后, 直接查es接口: 
127.0.0.1:8080/order/esList?page=1&pageSize=20&platform_code=AMAZON,EB&total_price_start=45&total_price_end=50

```