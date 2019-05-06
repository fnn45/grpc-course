package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-course/blog/blogpb"
	"log"
)

func main()  {

	fmt.Println("Blog client")

	opts := grpc.WithInsecure()
	conn, err := grpc.Dial("localhost:50051",  opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	c := blogpb.NewBlogServiceClient(conn)

	fmt.Println("Creating the blog")

	blog := &blogpb.Blog{
		AuthorId: "123qwertty",
		Title: "My first blog",
		Content: "Content of my first blog",
	}

	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog} )
	if err != nil {
		log.Fatal("Unxepected error", err)
	}

	fmt.Println("Blog has been created: %v", createBlogRes )

	// read Blog

	blogId := createBlogRes.GetBlog().GetId()

	fmt.Println("Reading the blog")
	readRes, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{ BlogId: blogId})
	if err != nil {
		fmt.Println("Error happend while reading %v", err)
	}

	fmt.Printf("Retriving blog : %#v", readRes.Blog)
}
