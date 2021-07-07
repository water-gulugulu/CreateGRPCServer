// @File  : server.go.go
// @Author: JunLong.Liao&此处不应有BUG!
// @Date  : 2021/7/7
// @slogan: 又是不想写代码的一天，神兽保佑，代码无BUG！

package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	UserRpc "rpcService/pb"
)

type UserServer struct {
}

// 获取用户详情
func (u UserServer) GetUserDetail(ctx context.Context, req *UserRpc.GetUserDetailReq) (*UserRpc.GetUserDetailRes, error) {

	// 处理逻辑
	return &UserRpc.GetUserDetailRes{
		Id:       req.Id,
		Username: "测试用户",
	}, nil
}

func main() {
	// 监听1234端口
	listen, err := net.Listen("tcp", "0.0.0.0:1234")

	if err != nil {
		log.Fatal(err)
		return
	}
	// 初始服务对象
	server := grpc.NewServer()
	// 将定义的用户Rpc注册
	UserRpc.RegisterUserServiceServer(server, new(UserServer))
	// 注册服务
	reflection.Register(server)
	// 注册监听的端口服务
	log.Println("服务启动了！")
	if err := server.Serve(listen); err != nil {
		fmt.Printf("服务启动失败: %s", err)
		return
	}
}
