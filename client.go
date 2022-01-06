package main

import (
	"context"
	"flag"
	pb "github.com/Example-Collection/go-grpc-client/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
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

func printResponseAfterCallingListPersons(client pb.PersonServiceClient, req *pb.ListPersonRequest) {
	log.Printf("Sending request to get all persons with email: %v", req.Email)
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	stream, err := client.ListPersons(ctx, req)
	if err != nil {
		log.Fatalf("%v.ListPersons(_) = _, %v", client, err)
	}
	for {
		person, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListPersons(_) = _, %v", client, err)
		}
		log.Printf("Response: Person(name: %v, message: %v)", person.Name, person.Age)
	}
}

func main() {
	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewPersonServiceClient(conn)

	// Unary RPC
	log.Println("Unary RPC!!")
	personRequest := &pb.PersonRequest{
		Name:     "Sangwoo",
		Age:      25,
		Email:    "robbyra@gmail.com",
		Password: "sangwooPassword",
	}
	printResponseAfterCallingGetPersonInformation(client, personRequest)

	// Server Streaming RPC
	log.Println("Server Streaming RPC!!")
	listPersonRequest := &pb.ListPersonRequest{Email: "robbyra@gmail.com"}
	printResponseAfterCallingListPersons(client, listPersonRequest)
}
