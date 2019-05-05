package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-course/calculator/calculator_pb"
	"io"
	"log"
)

func main()  {
	fmt.Println("Hello client")
	conn, err := grpc.Dial("localhost:50052",  grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	c := calculatorpb.NewCalculatorServiceClient(conn)
	doUnary(c)

	doServerStreaming(c)

	doClientStreaming(c)

	doErrorUnary(c)

}

func doUnary(c calculatorpb.CalculatorServiceClient)  {
	req := &calculatorpb.SumRequest{
		FirstNumber: 40,
		SecondNumber: 5,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatal("err while calling greet", err)
	}
	fmt.Println("response is: ", res.SumResult)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient)  {
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 12,
	}
	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatal("Error while calling PrimeDecomposition RPC", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("error while calling primedecomposition")
		}
		fmt.Println(res.Primefactor)
	}

}

func doClientStreaming(c calculatorpb.CalculatorServiceClient)  {
	fmt.Println("Starting compute avarage streaming....")
	stream, err := c.ComputeAvarage(context.Background())
	if err != nil {
		log.Fatal("error while computer avarage client streaming", err)
	}
	requests := []*calculatorpb.ComputeAvarageRequest{
		&calculatorpb.ComputeAvarageRequest{
			Number: 14,
		},
		&calculatorpb.ComputeAvarageRequest{
			Number: 32,
		},
		&calculatorpb.ComputeAvarageRequest{
			Number: 1,
		},
		&calculatorpb.ComputeAvarageRequest{
			Number: 56,
		},
		&calculatorpb.ComputeAvarageRequest{
			Number: 12,
		},
		&calculatorpb.ComputeAvarageRequest{
			Number: 19,
		},
	}

	for _, req := range requests {
		stream.Send(req)
	}

	result, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("Error while closing stream (compute avarage)", err)
	}
	fmt.Println("Compute Avarage Result is :", result.GetAvarage())

}

func doErrorUnary(c calculatorpb.CalculatorServiceClient)  {

	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{
		Number: 10,
	})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from grpc
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number")
			}

		} else {
			log.Fatal("Big error calling SquareRoot", err)
		}
	}
	fmt.Println("Result of squre is : ", res.GetNumberRoot())

}