package main

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	pb "todo-list-grpc-server/genproto"
	"todo-list-grpc-server/storage"
)

type ToDoService struct {
	storage *storage.PostgresStorage
	pb.UnimplementedToDoServiceServer
}

func NewToDoService(storage *storage.PostgresStorage) *ToDoService {
	return &ToDoService{storage: storage}
}

func (s *ToDoService) AddTask(ctx context.Context, req *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
	title := req.Title
	description := req.Description

	success, message := s.storage.AddTask(title, description)

	if success {
		return &pb.AddTaskResponse{Succes: true, Text: message}, nil
	} else {
		return &pb.AddTaskResponse{Succes: false, Text: message}, nil
	}
}

func (s *ToDoService) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	taskID := req.TaskId
	title := req.Title
	description := req.Description

	success, message := s.storage.UpdateTask(taskID, title, description)

	if success {
		return &pb.UpdateTaskResponse{Succes: true, Text: message}, nil
	} else {
		return &pb.UpdateTaskResponse{Succes: false, Text: message}, nil
	}
}

func (s *ToDoService) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	taskID := req.TaskId

	success, message := s.storage.DeleteTask(taskID)

	if success {
		return &pb.DeleteTaskResponse{Succes: true, Text: message}, nil
	} else {
		return &pb.DeleteTaskResponse{Succes: false, Text: message}, nil
	}
}

func (s *ToDoService) GetOneTask(ctx context.Context, req *pb.GetOneTaskRequest) (*pb.GetOneTaskResponse, error) {
	taskID := req.TaskId

	id, title, description, err := s.storage.GetOneTask(taskID)
	if err != nil {
		return nil, err
	}

	return &pb.GetOneTaskResponse{
		Task: &pb.Task{
			Id:          id,
			Title:       title,
			Description: description,
		},
	}, nil
}

func (s *ToDoService) GetAllTasks(ctx context.Context, req *pb.GetAllTasksRequest) (*pb.GetAllTasksResponse, error) {
	tasks, err := s.storage.GetAllTasks()
	if err != nil {
		return nil, err
	}

	var pbTasks []*pb.Task
	for _, task := range tasks {
		pbTasks = append(pbTasks, &pb.Task{
			Id:          task.Id,
			Title:       task.Title,
			Description: task.Description,
		})
	}

	return &pb.GetAllTasksResponse{Tasks: pbTasks}, nil
}

func main() {
	db, err := sql.Open("postgres", "user=postgres password=nodirbek dbname=tasks_grpc sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	storage := storage.NewPostgresStorage(db)

	toDoService := NewToDoService(storage)

	grpcServer := grpc.NewServer()

	pb.RegisterToDoServiceServer(grpcServer, toDoService)

	port := 50051
	listenAddr := ":" + strconv.Itoa(port)

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", port, err)
	}

	log.Printf("Server is listening on port %d...", port)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
