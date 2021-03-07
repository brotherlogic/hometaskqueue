package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/brotherlogic/hometaskqueue/proto"
)

func main() {
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal("failed to load system root CA cert pool")
	}
	creds := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(creds))
	conn, err := grpc.Dial("hometaskqueue-q2ijxfqena-uw.a.run.app:443", opts...)
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
		log.Fatalf("Bad requests: %v", err)
	}

	log.Printf("Result: %v", res)
}
