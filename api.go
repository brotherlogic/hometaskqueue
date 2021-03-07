package main

import (
	"context"
	"fmt"

	pb "github.com/brotherlogic/hometaskqueue/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) AddQueue(ctx context.Context, req *pb.AddQueueRequest) (*pb.AddQueueResponse, error) {
	return nil, status.Errorf(codes.InvalidArgument, "Bad request")
}
func (s *server) GetTasks(ctx context.Context, req *pb.GetTasksRequest) (*pb.GetTasksResponse, error) {
	return nil, fmt.Errorf("Bad request")
}
func (s *server) AddTask(ctx context.Context, req *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
	return nil, fmt.Errorf("Bad request")
}
