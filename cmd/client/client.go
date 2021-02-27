package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/leonardo-santoz/go-gRPC/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not coneect to gRPC Server: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	// AddUserVerbose(client)
	// AddUser(client)
	// AddUsers(client)
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Joao",
		Email: "j@j.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Joao",
		Email: "j@j.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the msg: %v", err)
		}
		fmt.Println("Status:", stream.Status, "-", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "l1",
			Name:  "Leonardo1",
			Email: "leonardo@teste1.com",
		},
		&pb.User{
			Id:    "l2",
			Name:  "Leonardo2",
			Email: "leonardo@teste2.com",
		},
		&pb.User{
			Id:    "l3",
			Name:  "Leonardo3",
			Email: "leonardo@teste3.com",
		},
		&pb.User{
			Id:    "l4",
			Name:  "Leonardo4",
			Email: "leonardo@teste4.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %V", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %V", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "l1",
			Name:  "Leonardo1",
			Email: "leonardo@teste1.com",
		},
		&pb.User{
			Id:    "l2",
			Name:  "Leonardo2",
			Email: "leonardo@teste2.com",
		},
		&pb.User{
			Id:    "l3",
			Name:  "Leonardo3",
			Email: "leonardo@teste3.com",
		},
		&pb.User{
			Id:    "l4",
			Name:  "Leonardo4",
			Email: "leonardo@teste4.com",
		},
	}

	wait := make(chan int)

	//goroutine
	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user:", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
			}
			fmt.Printf("Receiving user %v com status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}
