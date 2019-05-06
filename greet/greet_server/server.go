package main

import (
	"bytes"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"grpc-course/greet/greetpb"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct{}
var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)


func (s *server) GreetWithDeadline(ctx context.Context,req *greet_pb.GreetWithDeadlineRequest) (*greet_pb.GreetWithDeadlineResponse, error) {
	for i := 0 ; i < 3; i++ {
		time.Sleep(time.Second)
		if ctx.Err() == context.Canceled {
			// the client cancel request
			fmt.Println("client canceled the request!")
			return nil, status.Error(codes.Canceled, "Client cancel request")
		}
	}
	firstname := req.GetGreeting().FirstName
	res := &greet_pb.GreetWithDeadlineResponse{
		Result: "Hello " + firstname,
	}
	return res, nil
}

func (s *server) GreetEveryOne(stream greet_pb.GreetService_GreetEveryOneServer) error {
	fmt.Println("Start bidirectional communication...")
	buf := bytes.NewBufferString("Hello : ")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatal("Error while reading client stream", err)
			return err
		}
		firstname := req.GetGreeting().GetFirstName()
		fmt.Println("sending to server", firstname)
		buf.WriteString(firstname)
		buf.WriteString("! ,")
		err = stream.Send(&greet_pb.GreetEveryOneResponse{
			Result: buf.String(),
		})
		if err != nil {
			log.Fatal("Error while sending response to client", err)
			return err
		}
	}
}

func (s *server) LongGreet(stream greet_pb.GreetService_LongGreetServer) error {
	result := bytes.NewBufferString("Hello :")

	count := 0

for{
	count++
	if count == 3 {
		break
	}
	fmt.Println("count is :", count)
	req, err := stream.Recv()
	if err == io.EOF {
		rs := result.String()
		result.Reset()
	return stream.SendAndClose(&greet_pb.LongGreetResponse{
		Result: rs,
 	})
	}
	if err != nil {
		log.Fatal("Error while reading client stream: ", err)
	}
	firstname := req.GetGreeting().GetFirstName()
	result.WriteString(firstname)
}
return nil
}

func (s *server) GreetManyTimes(req *greet_pb.GreetManyTimesRequest,stream greet_pb.GreetService_GreetManyTimesServer) error {
	firstname := req.GetGreating().FirstName
	lastname := req.GetGreating().LastName
	for i := 0; i < 10; i++ {
		res := &greet_pb.GreetManyTimesResponse{
			Result: fmt.Sprintf("Hello %s %s %v", firstname, lastname, i),
		}
		stream.Send(res)
		time.Sleep(time.Second)
	}
return nil

}

func (s *server) Greet(ctx context.Context,req *greet_pb.GreetRequest) (*greet_pb.GreetResponse, error) {
	firstName :=  req.GetGreeting().FirstName
	result := "Hello " + firstName
	res := &greet_pb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func main()  {
	lis, err := net.Listen("tcp", "0.0.0.0:50051") // default port for grpc
	if err != nil {
		log.Fatal("Failed to listen %v", err)
	}
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(ensureValidToken),
	}
	tls := true
	if tls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatal("Failed loadig certificates",  sslErr)
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)
	greet_pb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve", err)
	}
}

func ensureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	fmt.Println("Interceptor triggered!!!!!")
	fmt.Println("Metadata is: ", md )
	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	//if !valid(md["authorization"]) {
	//	return nil, errInvalidToken
	//}
	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}
