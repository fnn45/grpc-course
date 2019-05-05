package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	calculatorpb "grpc-course/calculator/calculator_pb"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) SquareRoot(ctx context.Context,req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Received a  negative number",
			)
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil


}

func (*server) ComputeAvarage(stream calculatorpb.CalculatorService_ComputeAvarageServer) error {
	fmt.Println("Received Avarage calculate")
	sum, count := 0., 0.
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.ComputeAvarageResponse{
				Avarage: sum / count,
			})
		}
		if err != nil {
			log.Fatal("Error while reading client stream ....", err)
		}
		sum += float64(req.GetNumber())
		count++
	}
	return nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest,stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Received PrimeNumberDecomposition RPC: %v\n", req)
	number := req.GetNumber()
	divisor := int64(2)
	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				Primefactor: divisor,
			})
			number = number / divisor
		} else {
			divisor++
			fmt.Println("divisor has increased to :", divisor)
		}

	}
	return nil
}

func (*server) Sum(ctx context.Context,req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	firstnum := req.FirstNumber
	secnum := req.SecondNumber
	return &calculatorpb.SumResponse{
		SumResult: firstnum + secnum,
	}, nil
}

func main()  {
	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatal("Failed to listen %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve", err)
	}
}
