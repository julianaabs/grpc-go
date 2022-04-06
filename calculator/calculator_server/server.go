package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/julianaabs/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Calculator RPC: %v", req)
	firstNumer := req.FirstNumber
	secondNumber := req.SecondNumber
	sum := firstNumer + secondNumber
	result := &calculatorpb.SumResponse{
		SumResult: sum,
	}

	return result, nil
}

func (*server) PrimeNumberDecomp(req *calculatorpb.PrimeNumberDecompRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompServer) error {
	fmt.Printf("Received PrimeNumberDecomp: %v", req)

	number := req.GetNumber()
	divisor := int64(2)

	for number < 1 {
		if number%int64(divisor) == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompResponse{
				PrimeFactor: divisor,
			})
			number = number / divisor
		} else {
			divisor++
			fmt.Printf("Divisor increased to %v", divisor)
		}
	}

	return nil
}

func (*server) ComputerAverage(stream calculatorpb.CalculatorService_ComputerAverageServer) error {
	fmt.Printf("Received ComputerAverage RPC\n")

	sum := int32(0)
	count := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(count)
			stream.SendAndClose(&calculatorpb.ComputerAverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		sum += req.GetNumber()
		count++
	}
}

func main() {
	fmt.Println("Calculator Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
