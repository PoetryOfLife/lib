package main

import (
	"context"
	pb "gRPC/pb/server"
	"google.golang.org/grpc"
	"io"
	"log"
	"strconv"
)

const (
	// ServerAddress 连接地址
	ServerAddress string = ":8000"
)

func main() {

	// 连接服务器
	conn, err := grpc.Dial(ServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("net.Connect err: %v", err)
	}

	defer conn.Close()

	// 建立gRPC连接
	grpcClient := pb.NewHelloClient(conn)

	// 创建发送结构体
	req := pb.HelloRequest{
		Name: "grpc",
	}

	// 调用我们的服务(SayHello方法)
	// 同时传入了一个 context.Context ，在有需要时可以让我们改变RPC的行为，比如超时/取消一个正在运行的RPC
	res, err := grpcClient.SayHello(context.Background(), &req)
	if err != nil {
		log.Fatalf("Call SayHello err: %v", err)

	}
	log.Println(res)

	req1 := pb.ServerSideRequest{
		Name: "我来打开你啦",
	}
	stream, err := grpcClient.ServerSideHello(context.Background(), &req1)
	if err != nil {
		log.Fatalf("Call ServerSideHello err: %v", err)
	}
	for i := 0; i < 5; i++ {
		res1, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Conversations get stream err: %v", err)
		}
		// 打印返回值
		log.Println(res1.Message)
	}
	// 打印返回值

	res2, err := grpcClient.ClientSideHello(context.Background())
	if err != nil {
		log.Fatalf("Call ClientSideHello err: %v", err)
	}
	for i := 0; i < 5; i++ {
		err = res2.Send(&pb.ClientSideRequest{Name: "client"})
		if err != nil {
			return
		}
	}
	log.Println(res2.CloseAndRecv())

	stream1, err := grpcClient.BidirectionalHello(context.Background())
	if err != nil {
		log.Fatalf("get BidirectionalHello stream err: %v", err)
	}

	for i := 0; i < 5; i++ {
		err = stream1.Send(&pb.BidirectionalRequest{Name: "direction" + strconv.Itoa(i)})
		if err != nil {
			log.Fatalf("stream request err: %v", err)
		}
		res3, err := stream1.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Conversations get stream err: %v", err)
		}
		log.Println(res3.Message)
	}
}
