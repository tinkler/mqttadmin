## Mqttadmin
A development tool for mqtt
In the future, it will be a tool for mqtt management and monitoring
But now it's just a tool for generating service router and mqtt client code

English | [简体中文](README-zh_CN.md)

## Dependencies
* gorm
* [RabbitMQ](github.com/streadway/amqp)

## Requirements
* Go 1.20 or higher

## install
```
go install github.com/tinkler/mqttadmin/cmd/gen
```
#### or
```
git clone github.com/tinkler/mqttadmin
cd mqttadmin/cmd/gen
go install
```
### Use
Create a new project under your GOPATH using ```go mod``` init.
Place a [gen.yaml](gen.yaml) file in the project root and make the neccessary change you need.
Then run
```
gen .
```
This project has modified some common models, Such as user, role and permission, which can be found in the pkg/model directory. The generated route files was located in pkg/route directory.
## Future
Here is a list of supported generation options
+ chi_route Chi framework route
+ proto protocol buffer
+ gsrv gRPC server
+ dart Flutter Dio
+ ts XMLHttpRequest
+ angular_delon Delon service
+ swagger swagger api doc (In processing)
## Link
[ClanCloud a Saas project base on mqttadmin](https://github.com/tinkler/clancloud)
### License
The MIT License