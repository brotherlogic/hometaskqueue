package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"os"
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

	switch os.Args[1] {
	case "add":
		addFlags := flag.NewFlagSet("Add", flag.ExitOnError)
		var queue = addFlags.String("queue", "", "The name of the queue")
		var github = addFlags.String("github", "", "The name of the github base")
		var githubKey = addFlags.String("githubKey", "", "The github personal key")

		if err := addFlags.Parse(os.Args[2:]); err == nil {
			res, err := client.AddQueue(ctx, &pb.AddQueueRequest{
				QueueName: *queue,
				Github:    *github,
				GithubKey: *githubKey,
			})
			if err != nil {
				log.Fatalf("Bad requests: %v", err)
			}

			log.Printf("Result: %v", res)
		}

	case "task":
		taskFlags := flag.NewFlagSet("Task", flag.ExitOnError)
		var queue = taskFlags.String("queue", "", "The name of the queue")
		var taskName = taskFlags.String("name", "", "The name of the task")

		if err := taskFlags.Parse(os.Args[2:]); err == nil {
			res, err := client.AddTask(ctx, &pb.AddTaskRequest{
				Task: &pb.Task{
					Title: *taskName,
					Ttl:   10,
				},
				QueueId: *queue,
			})
			if err != nil {
				log.Fatalf("Bad add task: %v", err)
			}

			log.Printf("Result: %v", res)
		}

	case "tasks":
		taskFlags := flag.NewFlagSet("Tasks", flag.ExitOnError)
		var queue = taskFlags.String("queue", "", "The id of the queue")

		if err := taskFlags.Parse(os.Args[2:]); err == nil {
			res, err := client.GetTasks(ctx, &pb.GetTasksRequest{QueueId: *queue})
			if err != nil {
				log.Fatalf("Bad request: %v", err)
			}
			log.Printf("Result: %v", res)
		}
	}
}
