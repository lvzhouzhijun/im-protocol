# 使用Mage生成Protocol Buffers

此仓库仅生成Go和TS的pb文件。

## 先决条件

- Go 1.18或更高版本
- 在你的环境变量PATH中安装Protocol Buffer编译器（protoc）。
- 将mage bin安装到您的PATH环境变量中。

> 通知：
通常，`protoc`的不同版本不会产生显著影响，因为它们之间是相互兼容的。

## 安装 Go

- 请查看[Go安装文档](https://go.dev/doc/install)以安装Go。
- 将`go/bin`添加到PATH环境变量中。

## 安装Protocol Buffer编译器（protoc）

请查看[Protocol Buffer编译器安装文档](https://grpc.io/docs/protoc-installation/)。

- 下载与您的操作系统和架构相对应的最新版Protocol Buffers编译器压缩文件。
- 将文件解压到 `$HOME/.local` 目录下或您选择的目录中。
- 更新PATH环境变量，使其包含`protoc`可执行文件。

## 安装mage

使用 Go Install 并确保 Go 版本 >= 1.18：

```shell
执行命令：go install github.com/magefile/mage@latest
```

<details>
<summary>在 Go 版本 < 1.18 的情况下：</summary>

您可以使用`bootstrap_install_mage.bat`或`bootstrap_install_mage.sh`来快速安装mage。

</details>

## 编译您的Protocol Buffers

### Go 生成代码

- 执行`mage InstallDepend`来安装Go的依赖项。

- 执行`mage GenGo`以生成Go代码。

- 您还可以查看[Go语言使用文档](https://grpc.io/docs/languages/go/quickstart/#prerequisites)以获取更多信息。

### TypeScript 生成代码

- 在工作目录中执行 `npm install ts-proto`。
- 执行 `mage GenTypeScript` 以生成 TypeScript 代码。

## 修改Protocol Buffers

### 编写Protocol Buffers

在我们的示例中，我们有一个简单的`hello/hello.proto`文件：

```proto
syntax = "proto3";

// define a request message
message HelloRequest {
  string name = 1;
  UserInfo user = 2;
}

// define a response message
message HelloResponse {
  string message = 1;
}

// define a parameter message
message UserInfo {
  string name = 1;
  int32 age = 2;
}

// define a service
service HelloService {
  // define a rpc method
  rpc SayHello (HelloRequest) returns (HelloResponse);
}
```

编写你的方法请求和响应消息。例如`HelloRequest`和`HelloResponse`。
- 编写你的服务方法。比如`SayHello`。
- 你也可以定义参数消息，例如“UserInfo”。
- 将模块名称添加到`magefile.go`中的`protoModules`变量中。在此示例中，您需要在`protoModules`中添加`"hello"`。我们建议使用模块名称作为目录名称和文件名。
- 执行相应的语言命令以生成protobuf代码。更多信息请查看[编译您的Protocol Buffers](#compiling-your-protocol-buffers)。

## 更多信息：

- [Protocol Buffers文档](https://protobuf.dev/)
- [gRPC 文档](https://grpc.io/docs)