package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"grpc-course/greet/greetpb"
	"io"
	"log"
	"time"
)

func main()  {

	fmt.Println("Hello client")

	tls := true
	opts := grpc.WithInsecure()
	if tls {
		certFile := "ssl/ca.crt"
		creds, sslErr :=  credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatal("Error while loading CA trust certificate", sslErr)
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	conn, err := grpc.Dial("localhost:50051",  opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	c := greet_pb.NewGreetServiceClient(conn)
	doUnary(c)
	//
	//doServerStreaming(c)
	//
	//doClientStreaming(c)
	//
	//doBiDiStreaming(c)

	doUnaryWithDeadline(c, 5 * time.Second)
	doUnaryWithDeadline(c, 1 * time.Second)


}

func doUnary(c greet_pb.GreetServiceClient)  {
	req := &greet_pb.GreetRequest{
		Greeting : &greet_pb.Greeting{
			FirstName: "stephan",
			LastName:"xxx",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatal("err while calling greet", err)
	}
	fmt.Println("response is: ", res.Result)
}

func doServerStreaming(c greet_pb.GreetServiceClient)  {
	fmt.Println("Start server streaming rpc")

	req := &greet_pb.GreetManyTimesRequest{
		Greating: &greet_pb.Greeting{
			FirstName: "Stephane",
			LastName: "Maarek",
		},
	}
	stream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling stream")
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if  err != nil {
			log.Fatal(err)
		}

		log.Println(msg.GetResult())
	}

}

func doClientStreaming(c greet_pb.GreetServiceClient)  {
	fmt.Println("start Client streaming...")
	stream,  err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	requests := []*greet_pb.LongGreetRequest{
		&greet_pb.LongGreetRequest{
			Greeting: &greet_pb.Greeting{
				FirstName: "Ivan",
				LastName: "Ivanov",
			},
		},
		&greet_pb.LongGreetRequest{
			Greeting: &greet_pb.Greeting{
				FirstName: "John",
				LastName: "Ivanov",
			},
		},
		&greet_pb.LongGreetRequest{
			Greeting: &greet_pb.Greeting{
				FirstName: "Luct",
				LastName: "Ivanov",
			},
		},
		&greet_pb.LongGreetRequest{
			Greeting: &greet_pb.Greeting{
				FirstName: "Mark",
				LastName: "Ivanov",
			},
		},
		&greet_pb.LongGreetRequest{
			Greeting: &greet_pb.Greeting{
				FirstName: "UL",
				LastName: "Ivanov",
			},
		},
	}

	for _, req := range requests {
		if err := stream.Send(req); err != nil {
			fmt.Println("eroor !!!!!!!!")
			errStatus, _ := status.FromError(err)
			fmt.Println(errStatus.Code(), errStatus.Message(), errStatus.Details())
		}
		time.Sleep(time.Second * 2)
	}

	resp, err := stream.CloseAndRecv()


	if err != nil {

		errStatus, _ := status.FromError(err)
		fmt.Println(errStatus.Message())
		fmt.Println("status message: ")
		// lets print the error code which is `INVALID_ARGUMENT`
		fmt.Println(errStatus.Code())
		// Want its int version for some reason?
		// you shouldn't actullay do this, but if you need for debugging,
		// you can do `int(status_code)` which will give you `3`
		//
		// Want to take specific action based on specific error?


		log.Fatal("Error while sending long request to server ", err)
	}

	fmt.Println("LongRequest Result is ", resp.GetResult())
	}

func doBiDiStreaming(c greet_pb.GreetServiceClient)  {
	fmt.Println("Start BIDI Streaming....")

	// we create a stream by invoking the client

	// we send a bunch of messages to the client (go routine)

	// we receive a messages frin the server (go routine)

	stream, err := c.GreetEveryOne(context.Background())
	if err != nil {
		log.Fatal("error while streaming to server", err)
		return
	}

	waitc := make(chan struct{})
	//send
	go func() {
		requests := []*greet_pb.GreetEveryOneRequest{
			&greet_pb.GreetEveryOneRequest{
				Greeting: &greet_pb.Greeting{
					FirstName: "Ivan",
					LastName: "Ivanov",
				},
			},
			&greet_pb.GreetEveryOneRequest{
				Greeting: &greet_pb.Greeting{
					FirstName: "John",
					LastName: "Ivanov",
				},
			},
			&greet_pb.GreetEveryOneRequest{
				Greeting: &greet_pb.Greeting{
					FirstName: "Luct",
					LastName: "Ivanov",
				},
			},
			&greet_pb.GreetEveryOneRequest{
				Greeting: &greet_pb.Greeting{
					FirstName: "Mark",
					LastName: "Ivanov",
				},
			},
			&greet_pb.GreetEveryOneRequest{
				Greeting: &greet_pb.Greeting{
					FirstName: "UL",
					LastName: "Ivanov",
				},
			},
		}
		for _, req := range requests {
			stream.Send(req)
			time.Sleep(time.Second)
		}
		stream.CloseSend()
	}()


	// recieve
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				close(waitc)
				log.Fatal("Error while receiving from server", err)
			}
			fmt.Println("Recieved from BIDI: ", res.GetResult())
		}
		close(waitc)
	}()
	<- waitc

}

func doUnaryWithDeadline(c greet_pb.GreetServiceClient, timeout time.Duration)  {

	req := &greet_pb.GreetWithDeadlineRequest{
		Greeting : &greet_pb.Greeting{
			FirstName: "stephan",
			LastName:"xxx",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// 0r simple ctx, cancel := context.WithTimeout(context.Background(), 5* time.Second)

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
			fmt.Println("Timeout was hit! Deadline was exceeded")
			} else {
				fmt.Println("Unexpected error :", statusErr)
			}
			return
		} else {
			log.Fatal("err while calling greet", err)
		}

	}
	fmt.Println("response is: ", res.Result)

}
