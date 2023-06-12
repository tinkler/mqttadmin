# Mqttadmin
使用MQTT的接口快速开发工具。仅需定义结构体和方法就能生成对应前端架构的接口文件。
* 目前仅实现面向模型的接口生成功能。

[English](README.md) | 简体中文

## 依赖
* gorm

## 需求
使用了Golang泛型
* Go 1.20 or hight

## 安装
```
go install github.com/tinkler/mqttadmin/cmd/gen
```
#### 或者通过拷贝github项目进行安装
```
git clone github.com/tinkler/mqttadmin
cd mqttadmin/cmd/gen
go install
```
### 使用
新建Go project,建议在GOPAT下创建。将gen.yaml放到项目下，并根据自己的需要修改。
```
gen .
```
本项目定义了一些常用的模型，如用户，角色和权限。可以在pkg/model目录下找到相关例子。
## Future
当前已经实现生成如下架构所需API接口
+ chi_route Chi RESTful框架路由
+ proto Protocol buffer定义文件
+ gsrv gRPC服务实现
+ dart Flutter模型Class和Dio接口
+ ts 原生ES2006 XMLHttpRequest
+ angular_delon Delon中台框架Angular service
+ swagger swagger api doc (计划中)

## Link
[中华覃氏族谱](https://github.com/tinkler/clancloud)

### License

The MIT License