# jayden-framework-go

本项目流程及目的：业务数据保存mysql后，发送rabbitmq，启动go协程消费队列数据，将数据推送到es，业务多条件查询时，直接查es，解决大数据查询慢的问题。
  
 ```shell script
redis->mysql->rabbitmq(kafka)->go comsume->elasticsearch
```