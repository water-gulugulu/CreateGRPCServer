### 1. 基本概念
- RPC（Remote Procedure Call）远程过程调用，简单的理解是一个节点请求另一个节点提供的服务
- 本地过程调用：如果需要将本地student对象的age+1，可以实现一个addAge()方法，将student对象传入，对年龄进行更新之后返回即可，本地方法调用的函数体通过函数指针来指定。
- 远程过程调用：上述操作的过程中，如果addAge()这个方法在服务端，执行函数的函数体在远程机器上，如何告诉机器需要调用这个方法呢？
   1. 首先客户端需要告诉服务器，需要调用的函数，这里函数和进程ID存在一个映射，客户端远程调用时，需要查一下函数，找到对应的ID，然后执行函数的代码。
   1. 客户端需要把本地参数传给远程函数，本地调用的过程中，直接压栈即可，但是在远程调用过程中不再同一个内存里，无法直接传递函数的参数，因此需要客户端把参数转换成字节流，传给服务端，然后服务端将字节流转换成自身能读取的格式，是一个序列化和反序列化的过程。
   1. 数据准备好了之后，如何进行传输？网络传输层需要把调用的ID和序列化后的参数传给服务端，然后把计算好的结果序列化传给客户端，因此TCP层即可完成上述过程，gRPC中采用的是HTTP2协议。



### 2. 准备工作
##### 2.1 介绍

   - 要完成一个RPC首先需要实现一个服务端节点，由客户端去连接节点，再进行调用节点中的方法。
   - 这里介绍使用gRPC，借助gRPC，我们可以在`.proto`文件中定义我们的服务，并以gRPC支持的任何语言来实现客户端和服务器，客户端和服务器又可以在从服务器到你自己的平板电脑的各种环境中运行-gRPC还会为你解决所有不同语言和环境之间通信的复杂性。我们还获得了使用protocol buffer的所有优点，包括有效的序列化（速度和体积两方面都比JSON更有效率），简单的IDL（接口定义语言）和轻松的接口更新。
##### 2.2 安装

   - 首先需要安装gRPC golang版本的软件包，同时官方软件包的`examples`目录里就包含了教程中示例路线图应用的代码
   - 执行 `go get google.golang.org/grpc` 将所需的包拉下来
##### 2.3 安装相关工具和插件

   - 安装protocol buffer编译器，安装编译器最简单的方式是去[https://github.com/protocolbuffers/protobuf/releases](https://github.com/protocolbuffers/protobuf/releases) 下载预编译好的protoc二进制文件，仓库中可以找到每个平台对应的编译器二进制文件。这里我们以`Mac Os`为例，从[https://github.91chifun.workers.dev/https://github.com//protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-osx-x86_64.zip](https://github.91chifun.workers.dev/https://github.com//protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-osx-x86_64.zip) 下载并解压文件，更新`PATH`系统变量，或者确保`protoc`放在了`PATH`包含的目录中了。
   - 安装protoc编译器插件，执行`go get -u github.com/golang/protobuf/protoc-gen-go` 
   - 编译器插件`protoc-gen-go`将安装在`$GOBIN`中，默认位于`$GOPATH/bin`。编译器`protoc`必须在`$PATH`中能找到它：
   - 执行`export PATH=$PATH:$GOPATH/bin`
### 3. 定义服务

   - 首先第一步是用protocol buffer定义gRPC服务还有方法的请求和响应类型，首先要创建一个后缀为`.proto` 的protoobuf文件
   - 要定义服务，你需要在`.proto`文件中定义一个`service`      
```
syntax = "proto3";

package UserRpc;

// 用户服务
service UserService{
	...
}
```

   - 然后在服务中定义rpc方法，指定rpc方法的请求以及响应类型。
   - 定义一个简单的获取用户详情rpc，客户端使用用户ID将请求发送到服务器，然后等待响应返回，就像普通的函数调用一样。
   - `rpc GetUserDetail(GetUserDetailReq) returns (GetUserDetailRes) {}` 
   - 对应的请求以及响应类型需要通过`message`来声明具体的类
```
// 获取用户详情请求
message GetUserDetailReq {
    int64 id = 1;
}
// 获取用户详情响应
message GetUserDetailRes {
    int64  id       = 1;
    string username = 2;
}
```

   - 执行`protoc --go_out=. user.proto `来生成go的pb文件

![image.png](https://cdn.nlark.com/yuque/0/2021/png/21651265/1625623303390-b4e0f033-0967-49d3-b4cd-70bae34a8f6e.png#align=left&display=inline&height=77&margin=%5Bobject%20Object%5D&name=image.png&originHeight=154&originWidth=438&size=10616&status=done&style=none&width=219)

   - 执行后目录下就会生成`user.pb.go`文件，里面就有我们定义的rpc接口了，那么接下来创建`server,client`两个文件夹，分别是服务端和客户端
   - 首先来实现一个服务端，将rpc接口代码实现
   - 首先我们会需要用到grpc包，那么是要将这个包拉下来，加入到go.mod
   - 分别执行`go get -u google.golang.org/grpc` `go get -u google.golang.org/grpc/reflection`
   - 这样就可以继续我们的下一步了
   - ![image.png](https://cdn.nlark.com/yuque/0/2021/png/21651265/1625624607459-bd54f0d1-a713-43e9-99bf-7560a14b46f0.png#align=left&display=inline&height=530&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1060&originWidth=1964&size=184152&status=done&style=none&width=982)
   - ![image.png](https://cdn.nlark.com/yuque/0/2021/png/21651265/1625624629091-4735188a-b891-4b28-a066-30f65750cf0f.png#align=left&display=inline&height=233&margin=%5Bobject%20Object%5D&name=image.png&originHeight=466&originWidth=822&size=78127&status=done&style=none&width=411)
   - 鼠标选中定义的结构体然后按住option+回车（windows是alt+回车)，选中`Implement interface` 在弹出的窗口输入`user`就会出现`UserServiceServer`这个就是服务端代码的实现了
   - ![image.png](https://cdn.nlark.com/yuque/0/2021/png/21651265/1625624717179-0d1fe1eb-a5c8-4773-8f20-4d53fd6685e7.png#align=left&display=inline&height=219&margin=%5Bobject%20Object%5D&name=image.png&originHeight=438&originWidth=1180&size=81709&status=done&style=none&width=590)![image.png](https://cdn.nlark.com/yuque/0/2021/png/21651265/1625624963693-95761ffe-a5d6-4128-93f5-87c2e1311bd2.png#align=left&display=inline&height=210&margin=%5Bobject%20Object%5D&name=image.png&originHeight=420&originWidth=2626&size=66149&status=done&style=none&width=1313)
   - 回车选中后会自动将rpc接口的服务端方法实现，接下来在里面写服务端处理的逻辑
   - ![image.png](https://cdn.nlark.com/yuque/0/2021/png/21651265/1625627296342-0f592740-aa38-4c64-8348-d2e1ba6b6c6a.png#align=left&display=inline&height=801&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1602&originWidth=2674&size=280364&status=done&style=none&width=1337)
   - 实现了服务端的一个获取用户详情Rpc接口，那么客户端要请求到这个接口。哪就需要创建一个服务器来监听某个端口的请求，将Rpc服务注册，再映射到服务中去，这样一个服务端的Rpc就完成了，接下来打开终端 进入到server目录中启动
   - ![image.png](https://cdn.nlark.com/yuque/0/2021/png/21651265/1625627405967-0cff5a19-fd32-4380-b07a-e44afd3a9c24.png#align=left&display=inline&height=172&margin=%5Bobject%20Object%5D&name=image.png&originHeight=344&originWidth=922&size=27787&status=done&style=none&width=461)
   - 这样服务端就已经完成了，那么接下来实现客户端，让客户端来连接服务端，请求接口拿到返回的数据
   - 创建客户端，首先一样的在client文件夹中创建client.go文件
   - ![image.png](https://cdn.nlark.com/yuque/0/2021/png/21651265/1625628290445-a98a39c7-ffe0-4a8b-9bb8-47c3ffe1ab15.png#align=left&display=inline&height=831&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1662&originWidth=2582&size=345323&status=done&style=none&width=1291)
   - 然后在文件中连接到服务端节点以及监听的端口，再初始化UserService服务的连接，然后通过连接去调用该服务下面的获取详情接口，接口会返回两个参数，对应的服务端的返回，一个Res，一个error 
   - 接下来再打开一个终端  来运行一下client.go去连接到服务端获取用户详情的信息
   - ![image.png](https://cdn.nlark.com/yuque/0/2021/png/21651265/1625628436520-64553652-1ba1-4726-aa09-c79d5ef8f016.png#align=left&display=inline&height=158&margin=%5Bobject%20Object%5D&name=image.png&originHeight=316&originWidth=1060&size=32172&status=done&style=none&width=530)
   - 到这里一个简单的Rpc服务端和客户端就完成了。



### 4. 结尾

- 我们的示例是一个简单的获取用户详情的应用，客户端可以传入用户ID获取用户的详细信息。
- 借助gRPC，我们可以在`.proto`文件中定义我们的服务，并以gRPC支持的任何语言来实现客户端和服务器，客户端和服务器又可以在从服务器到你自己的平板电脑的各种环境中运行-gRPC还会为你解决所有不同语言和环境之间通信的复杂性。我们还获得了使用protocol buffer的所有优点，包括有效的序列化（速度和体积两方面都比JSON更有效率），简单的IDL（接口定义语言）和轻松的接口更新。
