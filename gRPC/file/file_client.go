package main

import (
	"context"
	pb "gRPC/pb/file"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
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

	grpcClient := pb.NewFileServiceClient(conn)

	stream, err := grpcClient.UploadFile(context.Background())
	if err != nil {
		log.Fatalf("err:%s", err.Error())
	}

	uploadedFile, err := os.Open("./file/test.txt")
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
		err = stream.Send(&pb.UploadFileRequest{
			FileName: "filename",
			Content:  buffer[:n],
		})
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
	stream.CloseAndRecv()

	req := &pb.DownloadFileRequest{FilePath: "./file/test.txt"}
	stream1, err := grpcClient.DownloadFile(context.Background(), req)
	if err != nil {
		log.Fatalf(err.Error())
	}

	var res []byte
	for {
		recv, err := stream1.Recv()
		if recv != nil && len(recv.Content) > 0 {
			res = append(res, recv.GetContent()...)
		}
		if err == io.EOF {
			log.Println("start save file")
			create, err := os.Create("./file/download_file/download.txt")
			if err != nil {
				log.Fatalf(err.Error())
			}
			_, err = create.Write(res)
			if err != nil {
				log.Fatalf(err.Error())
			}
			create.Close()
			break
		}
		if err != nil && err != io.EOF {
			log.Println("客户端 err=", err.Error())
		}
	}
}
