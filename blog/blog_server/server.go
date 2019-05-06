package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
	"grpc-course/blog/blogpb"
	"log"
	"net"
	"os"
	"os/signal"
	"reflect"
	"time"
)

var collection *mongo.Collection

type server struct {}

type blogItem struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string `bson:"author_id"`
	Content string `bson:"content"`
	Title string `bson:"title"`

}

func (server) CreateBlog(ctx context.Context,req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blog := req.GetBlog()
	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Title: blog.GetTitle(),
		Content: blog.GetContent(),
	}
	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err %v", err))
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	fmt.Println(reflect.TypeOf(oid))
	if !ok {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot convert to string"))
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id: oid.Hex(),
			AuthorId:blog.GetAuthorId(),
			Title: blog.GetTitle(),
			Content: blog.GetContent(),
		},
	}, nil
}

func (server) ReadBlog(ctx context.Context,req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	fmt.Println("Read blog request")
	blogID := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Cannot parse id"))
	}
	data := &blogItem{}
	filter := bson.D{{"_id", oid}}
	res := collection.FindOne(context.Background(), filter)
	if err = res.Decode(data); err != nil {
		return nil , status.Errorf(codes.NotFound, fmt.Sprintf("Cannot find blog with specified ID: %v", err))
	}
	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			Id : data.ID.Hex(),
			AuthorId: data.AuthorID,
			Content: data.Content,
			Title: data.Title,
		},
	}, nil
}

func main()  {

	// connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil { log.Fatal(err) }
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil { log.Fatal(err) }

	collection = client.Database("mydb").Collection("blog")

	lis, err := net.Listen("tcp", "0.0.0.0:50051") // default port for grpc
	if err != nil {
		log.Fatal("Failed to listen %v", err)
	}
	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting blog server...")
		log.Fatal(fmt.Errorf("failed while start blog server %v", s.Serve(lis)))

	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until signal is revieved

	<- ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing listener")
	lis.Close()
	fmt.Println("Closing MongoDB connection..")
	client.Disconnect(ctx)
	fmt.Println("End of program")
	}
	

