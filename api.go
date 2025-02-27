package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/brotherlogic/habridge/proto"
)

type EntityResponse struct {
	Entity_id string
	State     string
}

func (s *Server) GetState(ctx context.Context, req *pb.GetStateRequest) (*pb.GetStateResponse, error) {
	var respv *EntityResponse

	response, err := s.client.R().
		SetAuthToken(s.token).
		SetResult(respv).Get(fmt.Sprintf("http://%v/api/states/%v", s.url, req.GetButtonId()))

	if err != nil {
		return nil, err
	}

	if response.Status() != "OK" {
		return nil, status.Errorf(codes.Internal, "Unable to  read: %v", response.Status())
	}

	if respv == nil {
		return nil, status.Errorf(codes.Internal, "Bad result: %v", respv)
	}

	if respv.State == "on" {
		return &pb.GetStateResponse{ButtonState: pb.GetStateResponse_BUTTON_STATE_ON}, nil
	}
	if respv.State == "not_home" {
		return &pb.GetStateResponse{UserState: pb.GetStateResponse_USER_STATE_AWAY}, nil
	}
	if respv.State == "home" {
		return &pb.GetStateResponse{UserState: pb.GetStateResponse_USER_STATE_HOME}, nil
	}
	return nil, status.Errorf(codes.Unimplemented, "Can't interpret: %+v", respv)
}
