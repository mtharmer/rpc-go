package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/joho/godotenv"
	pb "github.com/mtharmer/rpc-go/rpcgo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

var secretToken = ""

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

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretToken = os.Getenv("RPC_SECRET_KEY")
	if secretToken == "" {
		log.Fatal("Secret token is empty")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	certFile := os.Getenv("CERT_FILE")
	keyFile := os.Getenv("KEY_FILE")

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load cert pair: %s", err)
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(ensureValidToken),
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}
	s := grpc.NewServer(opts...)

	pb.RegisterDoStuffServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	return token == secretToken
}

func ensureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	if !valid(md["authorization"]) {
		return nil, errInvalidToken
	}
	return handler(ctx, req)
}
