package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/mtharmer/rpc-go/rpcgo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDoStuffClient(conn)

	// Contact the server and print out its response.
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
