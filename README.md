# gokit Golang工具集

> [TOP]


## common 工具

## TCP socket封装

1、需要再次优化 事件问题，把所有回调改为事件处理机制，并且事件带有事件值

2、PING,PONG 双向互相检测问题，如果尤一方收到PING消息后就应该重置消息PING的时间

3、自动重连机制必须添加，添加自动重连处理。注意自动重连后的首次操作，

4、结构优化，start方法直接返回就是一定成功了。不需要单独等待
Start(onStart fun(Channel)...) error

## Logs 日志框架


## Config 配置工具

os.Expend,os.ExpendEnv 两个方法使用到string字段上


## 定义一些常用的异常

- ErrNotFound
- ErrInvalid

