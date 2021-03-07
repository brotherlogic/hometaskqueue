package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/brotherlogic/hometaskqueue/proto"
)

func main() {
	conn, err := grpc.Dial("https://hometaskqueue-q2ijxfqena-uw.a.run.app", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial Error: %v", err)
	}

	client := pb.NewHomeTaskQueueServiceClient(conn)
	if err != nil {
		log.Fatalf("Client err: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	res, err := client.AddQueue(ctx, &pb.AddQueueRequest{})
	if err != nil {
		log.Fatalf("Bad dial: %v", err)
	}

	log.Printf("Result: %v", res)
}
