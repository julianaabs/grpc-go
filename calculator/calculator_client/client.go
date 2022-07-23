package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/julianaabs/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Calculator Client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	doUnary(c)

	doServerStreaming(c)

}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &calculatorpb.SumRequest{
		FirstNumber:  75,
		SecondNumber: 28,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from Sum: %v", res.SumResult)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do Prime decomp serve streaming RPC...")
	req := &calculatorpb.PrimeNumberDecompRequest{
		Number: 12,
	}
	stream, err := c.PrimeNumberDecomp(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Prime decomp RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetPrimeFactor())

	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do ComputerAVG client streaming RPC...")

	stream, err := c.ComputerAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream: %v", err)
	}

	numbers := []int32{3, 5, 9, 54, 23}
	for _, number := range numbers {
		stream.Send(&calculatorpb.ComputerAverageRequest{
			Number: number,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}

	fmt.Printf("The average is %v", res)
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("SquarRoot unary rpc...")

	doErrorCall(c, 12)

	doErrorCall(c, -2)

}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {
	number := int32(0)
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: 10})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("Negative number")
			}
		} else {
			log.Fatalf("Big error calling squareroot: %v", err)
		}
	}
	fmt.Println("Result of square root of %v: %v", number, res.GetNumberRoot())

}
