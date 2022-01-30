# socketgo

socket with golang

本仓库大部分的代码来自 [wentby](https://github.com/secondtonone1/wentby) 。

最开始是直接goimport的方式使用，
因本人只需要原始socket方面的代码, 其它大部分代码用不到, 所以并没有fork完整项目。
只是简单的goimport的方式使用 `netmodel` 和 `protocol` 包中的内容。

后因为个人需要修改了部分内容，之后每次要使用相关代码都要copy文件很不方便。

所以提交上来, 方便使用。

# 使用方法

1. 项目goimport的方式加载项目中。
2. 参考example中的样式定义 `自己的msg`、`dispatcher`、`protocol` 即可。

# 样例的运行方法
Example的消息结构定义依赖 `google protobuf`, 因此需要提前安装好 `protoc` 和 `gofaster插件`。 然后按如下步骤执行:

1. 进入 `example/proto`目录。
2. 编译消息结构 `protoc --gofast_out=. *.proto`
3. 分别进入 client 和 server 目录，`go build`
4. 执行编译出来的 `server` 和 `client` 即可。

# 相关链接

* [wentby](https://github.com/secondtonone1/wentby)