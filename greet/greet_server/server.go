package main

import (
	"context"
	"fmt"
	"github.com/spencermeta/proto-example/greet/greetpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()

	r := "Hello, " + firstName + "!"
	response := &greetpb.GreetResponse{
		Result: r,
	}

	return response, nil
}

func (s *server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Printf("GreetWithDeadline function was invoked with %v\n", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("The client has cancelled the request")

			return nil, status.Error(codes.Canceled, "The client has cancelled the request")
		}

		time.Sleep(1 * time.Second)
	}
	firstName := req.GetGreeting().GetFirstName()

	r := "Hello, " + firstName + "!"
	response := &greetpb.GreetWithDeadlineResponse{
		Result: r,
	}

	return response, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()

	for i := 0; i < 10; i++ {
		result := "Hello, " + firstName + " for the " + strconv.Itoa(i) + " time"
		response := &greetpb.GreetManyTimesResponse{
			Result: result,
		}

		stream.Send(response)
		time.Sleep(1 * time.Second)
	}

	return nil
}

func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked with a streaming request.\n")
	result := "Hello, "
	for {
		req, err := stream.Recv()
		fmt.Printf("LongGreet function was invoked with %v\n", req)

		if err == io.EOF {
			fmt.Println("We have finished the stream.")
			response := &greetpb.LongGreetResponse{
				Result: result,
			}
			return stream.SendAndClose(response)
		}

		if err != nil {
			log.Fatalf("Error while reading the stream: %v\n", err.Error())

			return err
		}

		firstName := req.GetGreeting().GetFirstName()
		result += firstName + "! "
	}

	return nil
}

func (s *server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone function was invoked with a streaming request.\n")

	for {
		req, err := stream.Recv()
		fmt.Printf("GreetEveryone function was invoked with %v\n", req)

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading the stream: %v\n", err.Error())

			return err
		}

		result := "Hello, " + req.GetGreeting().GetFirstName()
		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})

		if err != nil {
			log.Fatalf("Error while sending data to client: %v\n", err.Error())
			return err
		}
	}

	return nil
}

func main() {
	fmt.Println("Greet server invoked!")

	/*	creds, err := credentials.NewServerTLSFromFile("certs/server.crt", "certs/server.pem")
		if err != nil {
			log.Fatalf("Failed loading certificates: %s", err.Error())
		}*/
	listener, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen %s", err.Error())
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	reflection.Register(s)

	err = s.Serve(listener)

	if err != nil {
		log.Fatalf("Failed to serve %s", err.Error())
	}
}
