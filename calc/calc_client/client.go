package main

import (
	"context"
	"fmt"
	"github.com/aerostatka/proto-example/calc/calcpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("You need to pass arguments!")
	}

	var numbers []float64

	for i := 1; i < len(os.Args); i++ {
		number, err := strconv.ParseFloat(os.Args[i], 10)
		if err != nil {
			log.Fatalf("Could not parse argument: %v", err.Error())
		}

		numbers = append(numbers, number)
	}

	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Could not connect: %v", err.Error())
	}

	defer cc.Close()

	c := calcpb.NewCalcServiceClient(cc)

	//doSum(c, numberOne, numberTwo)
	//doNumberDecomposition(c, number)
	//doAverage(c, numbers)
	//doMax(c, numbers)
	doErrorUnary(c, numbers)
}

func doErrorUnary(c calcpb.CalcServiceClient, numbers []float64) {
	fmt.Println("Starting doErrorUnary....")

	for _, number := range numbers {
		request := &calcpb.SquareRootRequest{
			Number: number,
		}

		response, err := c.SquareRoot(context.Background(), request)
		if err != nil {
			responseError, ok := status.FromError(err)
			if ok {
				fmt.Printf("Error with code %s and message %s\n", responseError.Code(), responseError.Message())
			} else {
				log.Fatalf("Could not connect: %v", err.Error())
			}

		}

		fmt.Printf("Returned response: %v\n", response.GetRoot())
	}
}

func doMax(c calcpb.CalcServiceClient, numbers []float64) {
	fmt.Println("Starting doMax....")

	stream, err := c.Max(context.Background())
	if err != nil {
		log.Fatalf("Could not connect: %v", err.Error())
	}

	waitc := make(chan struct{})

	go func() {
		for _, number := range numbers {
			fmt.Printf("Sending number %v\n", number)

			err = stream.Send(&calcpb.MaxRequest{
				Number: number,
			})
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

			fmt.Printf("Received a new maximum: %v\n", response.GetMax())
		}

		close(waitc)
	}()

	<-waitc
}

func doAverage(c calcpb.CalcServiceClient, numbers []float64) {
	fmt.Println("Starting doAverage....")

	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("Could not connect: %v", err.Error())
	}

	for _, number := range numbers {
		fmt.Printf("Sending number: %v\n", number)
		req := &calcpb.AverageRequest{
			Number: number,
		}

		err = stream.Send(req)
		if err != nil {
			fmt.Printf("Error while sending request: %s\n", err.Error())
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error in receiving average response: %v", err.Error())
	}

	fmt.Printf("Response received. Average is %v\n", response.GetAverage())
}

func doNumberDecomposition(c calcpb.CalcServiceClient, number int64) {
	fmt.Println("Starting doNumberDecomposition....")
	request := &calcpb.PrimeNumberDecompositionRequest{
		Number: number,
	}

	stream, err := c.PrimeNumberDecomposition(context.Background(), request)
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

		fmt.Printf("Returned response: %v\n", msg.GetPrimeNumber())
	}
}

func doSum(c calcpb.CalcServiceClient, nOne int64, nTwo int64) {
	fmt.Println("Starting doSum....")
	request := &calcpb.SumRequest{
		NumberOne: nOne,
		NumberTwo: nTwo,
	}

	response, err := c.Sum(context.Background(), request)
	if err != nil {
		log.Fatalf("Could not connect: %v", err.Error())
	}

	fmt.Printf("Returned response: %v\n", response.GetResult())
}
