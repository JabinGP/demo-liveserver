# demo-liveserver

简单实现一个Live Server，监听一个文件，当文件更改时，自动刷新浏览器。

## 项目启动

```cmd
git clone https://github.com/JabinGP/demo-liveserver.git
cd demo-liveserver
go run main.go
```

如果出现端口冲突问题，可以在main.go中修改端口

```go
app.Run(iris.Addr(":1454"), iris.WithoutServerError(iris.ErrServerClosed))
```

启动后访问[http://localhost:1454/test/index.html](http://localhost:1454/test/index.html)可以浏览示例html文件

修改/test/index.html并保存可以触发浏览器的刷新！

## 实现原理

有空了更新~~