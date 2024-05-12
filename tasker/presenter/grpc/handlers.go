//go:build grpc

package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/teran/backupnizza/tasker/presenter/grpc/proto"
	"github.com/teran/backupnizza/tasker/service"
)

type Handlers interface {
	proto.TaskerServiceServer

	Register(*grpc.Server)
}

type handlers struct {
	svc service.Tasker
}

func New(svc service.Tasker) Handlers {
	return &handlers{
		svc: svc,
	}
}

func (h *handlers) RunTask(ctx context.Context, in *proto.RunTaskRequest) (*proto.RunTaskResponse, error) {
	task, err := h.svc.GetByName(in.GetName())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = task.Run(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.RunTaskResponse{}, nil
}

func (h *handlers) Register(srv *grpc.Server) {
	proto.RegisterTaskerServiceServer(srv, h)
}
