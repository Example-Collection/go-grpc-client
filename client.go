package main

import (
	"context"
	"flag"
	pb "github.com/Example-Collection/go-grpc-client/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var (
	serverAddr = flag.String("addr", "localhost:8081", "The server address in the format of host:port")
)

func printResponseAfterCallingGetPersonInformation(client pb.PersonServiceClient, req *pb.PersonRequest) {
	log.Printf("Sending request to save person (name: %v, age: %d, email: %v)", req.Name, req.Age, req.Email)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	response, err := client.GetPersonInformation(ctx, req)
	if err != nil {
		log.Fatalf("%v.GetPersonInformation(_) = _, %v", client, err)
	}
	log.Printf("Response: Person(name: %v, message: %v", response.Name, response.Message)
}

func main() {
	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewPersonServiceClient(conn)

	// Unary RPC
	personRequest := &pb.PersonRequest{
		Name:     "Sangwoo",
		Age:      25,
		Email:    "robbyra@gmail.com",
		Password: "sangwooPassword",
	}
	printResponseAfterCallingGetPersonInformation(client, personRequest)
}
