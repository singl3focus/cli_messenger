package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/singl3focus/cli_messenger/auth/internal/adapters/postgres"
	"github.com/singl3focus/cli_messenger/auth/internal/config"
	"github.com/singl3focus/cli_messenger/auth/internal/core/interfaces"
	"github.com/singl3focus/cli_messenger/auth/internal/core/models"
	desc "github.com/singl3focus/cli_messenger/auth/pkg/auth_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedAuthV1Server

	Repostory interfaces.Repostory
}

func (s *server) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	const op = "handler.CreateUser"

	if req.User.GetPassword() != req.GetPasswordConfirm() {
		return nil, status.Errorf(codes.InvalidArgument, "%s: invalid password confirm", op)
	}

	user, err := models.NewUser(
		req.User.GetName(),
		req.User.GetEmail(),
		req.User.GetPassword(),
		int(req.User.GetRole()),
	)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s: creating user model has failed", op)
	}

	log.Printf("Create User: %v", user)

	id, err := s.Repostory.CreateUser(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "unexpected error: %s", err.Error())
	}
	
	return &desc.CreateUserResponse{
		Id: id,
	}, nil
}

func (s *server) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	const op = "handler.GetUser"

	id := req.GetId()
	if id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%s: invalid user ID: must be positive integer, got %d", op, id)
	}
	
	log.Printf("Get User id: %d", id)

	user, err := s.Repostory.GetUser(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "unexpected error: %s", err.Error())
	}

	return &desc.GetUserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      desc.Role(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

func (s *server) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	const op = "handler.UpdateUser"

	id := req.GetId()
	if id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%s: invalid user ID: must be positive integer, got %d", op, id)
	}

	var name *string
	if req.Name != nil {
        n := strings.TrimSpace(req.Name.GetValue())
        if n == "" {
            return nil, status.Errorf(codes.InvalidArgument, "%s: name cannot be empty", op)
        }
        name = &n
    }

	var email *string
	if req.Email != nil {
        e := strings.TrimSpace(req.Email.GetValue())
        if e == "" {
            return nil, status.Errorf(codes.InvalidArgument, "%s: invalid email format", op)
        }
        email = &e
    }
	
	log.Printf("Update User id: %d", id)

	err := s.Repostory.UpdateUser(ctx, id, name, email)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "unexpected error: %s", err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	const op = "handler.DeleteUser"

	id := req.GetId()
	if id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%s: invalid user ID: must be positive integer, got %d", op, id)
	}
	
	log.Printf("Delete User id: %d", id)

	err := s.Repostory.DeleteUser(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "unexpected error: %s", err.Error())
	}

	return &emptypb.Empty{}, nil
}

func main() {
	cfg := config.NewConfig(config.ENV)

	if err := cfg.Load(configPath); err != nil {
		log.Fatalf("%s (path %s)", err, configPath)
	}

	db, err := postgres.NewPostgres(cfg.PGDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{Repostory: db})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
