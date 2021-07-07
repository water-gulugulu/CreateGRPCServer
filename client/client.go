// @File  : client.go
// @Author: JunLong.Liao&此处不应有BUG!
// @Date  : 2021/7/7
// @slogan: 又是不想写代码的一天，神兽保佑，代码无BUG！

package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	UserRpc "rpcService/pb"
)

func main() {
	// 连接到Rpc节点
	conn, err := grpc.Dial(":1234", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Rpc连接失败！")
		return
	}
	// 初始化rpc服务连接
	client := UserRpc.NewUserServiceClient(conn)
	// 调用Rpc服务的用户详情接口
	res, err := client.GetUserDetail(context.Background(), &UserRpc.GetUserDetailReq{
		Id: 1,
	})
	if err != nil {
		fmt.Printf("用户详情Rpc请求失败，错误：%s\n", err.Error())
		return
	}

	fmt.Printf("用户详情Rpc请求成功，返回ID：%v，username：%v\n", res.Id, res.Username)
}
