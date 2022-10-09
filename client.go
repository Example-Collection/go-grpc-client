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

func savePersons(client pb.PersonServiceClient) {
	requests := []*pb.PersonRequest{
		{
			Email:    "email1@test.com",
			Age:      1,
			Name:     "name1",
			Password: "password1",
		},
		{
			Email:    "email2@test.com",
			Age:      2,
			Name:     "name2",
			Password: "password2",
		},
		{
			Email:    "email3@test.com",
			Age:      3,
			Name:     "name3",
			Password: "password3",
		},
	}

	stream, err := client.SavePersons(context.Background())
	if err != nil {
		log.Fatalf("%v.SavePersons(_) = _, %v", client, err)
	}
	for _, req := range requests {
		time.Sleep(time.Second)
		if err := stream.Send(req); err != nil {
			log.Fatalf("%v.Send(%v) = %v", stream, req, err)
		}
	}
	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Printf("SavePersons() response: %v\n", response)
}

func printResponseInAskAndGetPersons(client pb.PersonServiceClient) {
	requests := []*pb.PersonRequest{
		{
			Email:    "email1@test.com",
			Age:      1,
			Name:     "name1",
			Password: "password1",
		},
		{
			Email:    "email2@test.com",
			Age:      2,
			Name:     "name2",
			Password: "password2",
		},
		{
			Email:    "email3@test.com",
			Age:      3,
			Name:     "name3",
			Password: "password3",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	stream, err := client.AskAndGetPersons(ctx)
	if err != nil {
		log.Fatalf("%v.AskAndGetPersons(_) = _, %v", client, err)
	}

	waitChannel := make(chan bool)
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(waitChannel)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive PersonResponse : %v", err)
			}
			log.Printf("PersonResponse(email: %v, name: %v, age: %d) arrived.\n", res.Email, res.Name, res.Age)
		}
	}()

	for _, req := range requests {
		time.Sleep(time.Second)
		log.Printf("PersonRequest(email: %v, name: %v, age: %d) sent.", req.Email, req.Name, req.Age)
		if err := stream.Send(req); err != nil {
			log.Fatalf("Failed to send PersonRequest: %v", err)
		}
	}

	stream.CloseSend()
	<-waitChannel
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

	// Client Streaming RPC
	savePersons(client)

	// Bidirectional Streaming RPC
	printResponseInAskAndGetPersons(client)
}
