package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/grpc"

	firebase "firebase.google.com/go"

	pb "github.com/brotherlogic/hometaskqueue/proto"
)

var (
	localMode = flag.Bool("local", false, "Startup in local only mode")
)

type server struct {
}

func getClient(ctx context.Context, useLocal bool) *firebase.App {
	if useLocal {
		sa := option.WithCredentialsFile("firestore.json")
		app, err := firebase.NewApp(ctx, nil, sa)
		if err != nil {
			log.Fatalln(err)
		}
		return app
	}
	conf := &firebase.Config{ProjectID: "hometaskqueue"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}
	return app
}

func main() {
	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := grpc.NewServer()
	server := &server{}

	// Only init remote options if we're not running locally
	if !*localMode {

	}

	pb.RegisterHomeTaskQueueServiceServer(s, server)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Error grpc listen: %v", err)
	}
	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
