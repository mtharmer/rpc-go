package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	pb "github.com/mtharmer/rpc-go/rpcgo"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

const (
	defaultAddress = "localhost:50051"
	defaultName    = "world"
)

var (
	addr = flag.String("addr", defaultAddress, "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
	age  = flag.Int("age", 0, "Age of the person")
	city = flag.String("city", "", "City of the person")
)

var secretToken = ""

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

	certFile := os.Getenv("CERT_FILE")
	hostname := os.Getenv("RPC_HOSTNAME")

	perRPC := oauth.TokenSource{TokenSource: oauth2.StaticTokenSource(fetchToken())}
	creds, _ := credentials.NewClientTLSFromFile(certFile, hostname)

	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(perRPC),
		grpc.WithTransportCredentials(creds),
	}

	conn, err := grpc.Dial(*addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDoStuffClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.PrintHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())

	if *city != "" && *age > 0 {
		resp, err := c.ProcessPerson(ctx, &pb.PersonRequest{Name: *name, Age: int32(*age), City: *city})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", resp.GetMessage())
	}
}

func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: secretToken,
	}
}
