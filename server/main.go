package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/mtharmer/rpc-go/rpcgo"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedDoStuffServer
}

func (s *server) PrintHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello again " + in.GetName()}, nil
}

func (s *server) ProcessPerson(ctx context.Context, in *pb.PersonRequest) (*pb.PersonReply, error) {
	return &pb.PersonReply{Message: fmt.Sprintf("Hello %s, you are %d years old from %s", in.GetName(), in.GetAge(), in.GetCity()), Status: 200}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDoStuffServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
