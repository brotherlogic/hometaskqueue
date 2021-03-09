package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/brotherlogic/hometaskqueue/proto"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func makeID() string {
	return stringWithCharset(32, charset)
}

//Storage is the basic storage system
type Storage struct {
	client *firestore.Client
}

//NewStorage builds out a new storage
func NewStorage() *Storage {
	app := getClient(context.Background(), false)
	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	return &Storage{client: client}
}

//GetQueue gets a queue from the db
func (st *Storage) GetQueue(ctx context.Context, queueName string) (*pb.Queue, error) {
	ud, err := st.client.Collection("queues").Doc(queueName).Get(ctx)
	code := status.Convert(err)
	if code.Code() == codes.OK {
		queue := &pb.Queue{}
		ud.DataTo(queue)
		return queue, nil
	}

	return nil, err
}

//GetQueueID gets a queue from the db
func (st *Storage) GetQueueID(ctx context.Context, queueID string) (*pb.Queue, error) {
	ud, err := st.client.Collection("queues").Where("Id", "==", queueID).Documents(ctx).GetAll()
	code := status.Convert(err)
	if code.Code() == codes.OK {
		if len(ud) != 1 {
			return nil, status.Errorf(codes.Internal, "Too many queues returned (%v)", len(ud))
		}

		queue := &pb.Queue{}
		ud[0].DataTo(queue)
		return queue, nil
	}

	return nil, err
}

//SaveQueue Saves the queue
func (st *Storage) SaveQueue(ctx context.Context, queue *pb.Queue) error {
	_, err := st.client.Collection("queues").Doc(queue.GetQueueName()).Set(ctx, queue)
	return err
}

func (s *server) AddQueue(ctx context.Context, req *pb.AddQueueRequest) (*pb.AddQueueResponse, error) {
	if req.GetQueueName() != "brotherlogic" {
		return nil, status.Errorf(codes.InvalidArgument, "Malformed request")
	}

	if time.Now().Sub(s.lastRequest) < time.Hour {
		return nil, status.Errorf(codes.InvalidArgument, "Only one request per hour please")
	}

	s.lastRequest = time.Now()

	st := NewStorage()
	_, err := st.GetQueue(ctx, req.GetQueueName())
	if err == nil {
		return nil, status.Errorf(codes.InvalidArgument, "This queue already exists")
	}

	queue := &pb.Queue{
		Id:        makeID(),
		QueueName: req.GetQueueName(),
		Github:    req.GetGithub(),
		GithubKey: req.GetGithubKey(),
	}

	return &pb.AddQueueResponse{Added: queue}, st.SaveQueue(ctx, queue)
}
func (s *server) GetTasks(ctx context.Context, req *pb.GetTasksRequest) (*pb.GetTasksResponse, error) {
	st := NewStorage()
	q, err := st.GetQueueID(ctx, req.GetQueueId())
	if err != nil {
		return nil, err
	}

	var tasks []*pb.Task
	for _, task := range q.GetTasks() {
		if (req.GetSince() == 0 || task.GetDateAdded() < req.GetSince()) &&
			(req.GetType() != pb.TaskType_UNKNOWN || task.GetType() == req.GetType()) {
			tasks = append(tasks, task)
		}
	}

	return &pb.GetTasksResponse{
		Tasks: tasks,
	}, nil
}
func (s *server) AddTask(ctx context.Context, req *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
	st := NewStorage()
	q, err := st.GetQueueID(ctx, req.GetQueueId())
	if err != nil {
		return nil, err
	}

	task := req.GetTask()
	task.DateAdded = time.Now().Unix()
	q.Tasks = append(q.Tasks, task)

	return &pb.AddTaskResponse{}, st.SaveQueue(ctx, q)
}
