package main

import (
	pb "gRPC/pb/file"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
)

const (

	// Address 监听地址
	Address string = ":8000"

	// Network 网络通信协议
	Network string = "tcp"
)

type FileService struct {
	pb.UnimplementedFileServiceServer
}

func (s *FileService) UploadFile(server pb.FileService_UploadFileServer) error {
	var res []byte
	for {
		recv, err := server.Recv()
		if recv != nil && len(recv.Content) > 0 {
			res = append(res, recv.GetContent()...)
		}
		if err == io.EOF {
			log.Println("start save file")
			create, err := os.Create("./file/upload_file/upload.txt")
			if err != nil {
				return err
			}
			_, err = create.Write(res)
			if err != nil {
				return err
			}
			create.Close()
			break
		}
		if err != nil && err != io.EOF {
			log.Println("服务端 err=", err.Error())
			err = server.SendAndClose(&pb.UploadFileResponse{
				FilePath: "",
			})
			if err != nil {
				return err
			}
		}
	}
	server.SendAndClose(&pb.UploadFileResponse{
		FilePath: "./file/download_file/upload.txt",
	})
	return nil
}

func (s *FileService) DownloadFile(request *pb.DownloadFileRequest, server pb.FileService_DownloadFileServer) error {
	uploadedFile, err := os.Open(request.GetFilePath())
	if err != nil {
		log.Fatalf("err:%s", err.Error())
	}
	defer uploadedFile.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := uploadedFile.Read(buffer)
		if err != nil && err != io.EOF {
			log.Fatalf(err.Error())
		}
		if n == 0 {
			break
		}
		err = server.Send(&pb.DownloadFileResponse{
			Content: buffer[:n],
		})
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
	return nil
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
	pb.RegisterFileServiceServer(grpcServer, &FileService{})
	//用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
	err = grpcServer.Serve(listener)

	if err != nil {
		log.Panic("grpcServer.Serve err: %v", err)
	}
}
