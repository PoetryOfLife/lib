package main

import (
	"context"
	pb "gRPC/pb/server"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (

	// Address 监听地址
	Address string = ":8000"

	// Network 网络通信协议
	Network string = "tcp"
)

// HelloService 定义我们的服务
type HelloService struct {
	pb.UnimplementedHelloServer
}

// SayHello 实现SayHello方法
func (s *HelloService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResp, error) {
	log.Println(req.Name)
	return &pb.HelloResp{Message: "hello"}, nil
}

func (s *HelloService) ServerSideHello(request *pb.ServerSideRequest, server pb.Hello_ServerSideHelloServer) error {
	log.Println(request.Name)
	for n := 0; n < 5; n++ {
		// 向流中发送消息， 默认每次send送消息最大长度为`math.MaxInt32`bytes
		err := server.Send(&pb.ServerSideResp{Message: "你好"})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *HelloService) ClientSideHello(server pb.Hello_ClientSideHelloServer) error {
	for i := 0; i < 5; i++ {
		recv, err := server.Recv()
		if err != nil {
			return err
		}
		log.Println("客户端信息：", recv)
	}
	err := server.SendAndClose(&pb.ClientSideResp{
		Message: "close",
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *HelloService) BidirectionalHello(server pb.Hello_BidirectionalHelloServer) error {
	defer func() {
		log.Println("client close!")
	}()
	for {
		recv, err := server.Recv()
		if err != nil {
			return err
		}
		log.Println(recv)
		err = server.Send(&pb.BidirectionalResp{
			Message: "server message",
		})
		if err != nil {
			return err
		}
	}
}

func main() {
	// 监听本地端口
	listener, err := net.Listen(Network, Address)
	if err != nil {
		log.Panic("net.Listen err: %v", err)
	}

	log.Println(Address + " net.Listing...")

	// 新建gRPC服务器实例
	grpcServer := grpc.NewServer()

	// 在gRPC服务器注册我们的服务
	pb.RegisterHelloServer(grpcServer, &HelloService{})
	//用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
	err = grpcServer.Serve(listener)

	if err != nil {
		log.Panic("grpcServer.Serve err: %v", err)
	}
}
