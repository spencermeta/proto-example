package main

import (
	"context"
	"fmt"
	"github.com/spencermeta/proto-example/calc/calcpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
	"net"
)

type calcServer struct{}

func (s *calcServer) Sum(ctx context.Context, request *calcpb.SumRequest) (*calcpb.SumResponse, error) {
	fmt.Printf("Reseived request with parameters: %v\n", request)
	firstNumber := request.GetNumberOne()
	secondNumber := request.GetNumberTwo()
	sum := firstNumber + secondNumber

	return &calcpb.SumResponse{
		Result: sum,
	}, nil
}

func (s *calcServer) SquareRoot(ctx context.Context, request *calcpb.SquareRootRequest) (*calcpb.SquareRootResponse, error) {
	fmt.Printf("Reseived request with parameters: %v\n", request.GetNumber())
	number := request.GetNumber()

	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received the negative number %v", number),
		)
	}

	root := math.Sqrt(number)

	return &calcpb.SquareRootResponse{
		Root: root,
	}, nil
}

func (s *calcServer) PrimeNumberDecomposition(req *calcpb.PrimeNumberDecompositionRequest, stream calcpb.CalcService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition function was invoked with %v\n", req)

	number := req.GetNumber()
	var divisor int64 = 2

	for number > 1 {
		if number%divisor == 0 {
			response := &calcpb.PrimeNumberDecompositionResponse{
				PrimeNumber: divisor,
			}

			number = number / divisor

			err := stream.Send(response)
			if err != nil {
				fmt.Printf("Message sending error: %s\n", err.Error())
			}
		} else {
			divisor++
		}
	}

	return nil
}

func (s *calcServer) Average(stream calcpb.CalcService_AverageServer) error {
	fmt.Printf("Average function was invoked with a streaming request\n")
	var numbers []float64

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			fmt.Printf("The end of messages streaming.\n")

			result := 0.0
			for _, number := range numbers {
				result += number
			}

			result = result / float64(len(numbers))

			response := &calcpb.AverageResponse{
				Average: result,
			}

			return stream.SendAndClose(response)
		}

		if err != nil {
			fmt.Printf("Message receiving error: %s\n", err.Error())

			return nil
		}

		fmt.Printf("Received number %v\n", req.GetNumber())
		numbers = append(numbers, req.GetNumber())
	}
}

func (s *calcServer) Max(stream calcpb.CalcService_MaxServer) error {
	fmt.Printf("Max function was invoked with a streaming request\n")
	max := 0.0

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			fmt.Printf("Message receiving error: %s\n", err.Error())
			return err
		}

		fmt.Printf("Received number %v\n", req.GetNumber())
		number := req.GetNumber()
		if number > max {
			max = number
			err = stream.Send(&calcpb.MaxResponse{
				Max: max,
			})

			if err != nil {
				fmt.Printf("Error while seidning message: %s\n", err.Error())
			}
		}
	}
}

func main() {
	fmt.Println("Calc server invoked!")

	listener, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen %s", err.Error())
	}

	s := grpc.NewServer()

	calcpb.RegisterCalcServiceServer(s, &calcServer{})

	err = s.Serve(listener)

	if err != nil {
		log.Fatalf("Failed to serve %s", err.Error())
	}
}
