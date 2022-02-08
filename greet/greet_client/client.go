package main

import (
	"context"
	"fmt"
	"github.com/spencermeta/proto-example/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
)

func main() {
	/*	creds, err := credentials.NewClientTLSFromFile("certs/ca.crt", "")
		if err != nil {
			log.Fatalf("Error while loading certificates: %v", err.Error())
		}*/

	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect: %v", err.Error())
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	doUnary(c)
	//doServerStreaming(c)
	//doStreaming(c)
	//doBiDirectionalStreaming(c)
	//doUnaryWithDeadline(c, 1 * time.Second)
	//doUnaryWithDeadline(c, 5 * time.Second)
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting doUnaryWithDeadlines....")
	request := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Daria",
			LastName:  "Frolova",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	response, err := c.GreetWithDeadline(ctx, request)
	if err != nil {
		statusErr, ok := status.FromError(err)

		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded.")
			} else {
				log.Fatalf("Unexpected error: %v", statusErr.Err())
			}
		} else {
			log.Fatalf("Error while calling GreetWithDeadline: %v", err.Error())
		}
	}

	fmt.Printf("Returned response: %v\n", response.GetResult())
}

func doBiDirectionalStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting doBiDirectionalStreaming....")

	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("Error while calling  GreetEveryone: %v", err.Error())
	}

	waitc := make(chan struct{})
	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Daria",
				LastName:  "Frolova",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Aleksandr",
				LastName:  "Frolov",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Banana",
				LastName:  "Frolova",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Eldar",
				LastName:  "Gadirov",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Marina",
				LastName:  "Gadirova",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Nicole",
				LastName:  "Gadirova",
			},
		},
	}

	go func() {
		for _, request := range requests {
			fmt.Printf("Sending message %v\n", request)
			err = stream.Send(request)
			if err != nil {
				fmt.Printf("Error while sending messages: %s\n", err.Error())
			}
			time.Sleep(1 * time.Second)
		}

		err = stream.CloseSend()
		if err != nil {
			fmt.Printf("Error while closing connection: %s\n", err.Error())
		}
	}()

	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				fmt.Printf("Error while receiving messages: %s\n", err.Error())
				break
			}

			fmt.Printf("Received message %v\n", response)
		}

		close(waitc)
	}()

	<-waitc
}

func doStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting doStreaming....")

	stream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("Error while calling LongGreet: %v", err.Error())
	}

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Daria",
				LastName:  "Frolova",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Aleksandr",
				LastName:  "Frolov",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Banana",
				LastName:  "Frolova",
			},
		},
	}

	for _, request := range requests {
		fmt.Printf("Sending request for %s\n", request.GetGreeting().GetFirstName())
		err = stream.Send(request)
		if err != nil {
			log.Fatalf("Error while sending request long greet: %v", err.Error())
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from long greet: %v", err.Error())
	}

	fmt.Printf("Received response: %v\n", response)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting doServerStreaming....")
	request := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Daria",
			LastName:  "Frolova",
		},
	}

	stream, err := c.GreetManyTimes(context.Background(), request)
	if err != nil {
		log.Fatalf("Could not connect: %v", err.Error())
	}

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Something went wrong: %v", err.Error())
		}

		fmt.Printf("Returned response: %v\n", msg.GetResult())
	}

}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting doUnary....")
	request := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Daria",
			LastName:  "Frolova",
		},
	}

	response, err := c.Greet(context.Background(), request)
	if err != nil {
		log.Fatalf("Could not connect: %v", err.Error())
	}

	fmt.Printf("Returned response: %v\n", response.GetResult())
}
