package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/brotherlogic/habridge/proto"
)

type EntityResponse struct {
	EntityId string `json:"entity_id"`
	State    string `json:"state"`
}

func (s *Server) GetState(ctx context.Context, req *pb.GetStateRequest) (*pb.GetStateResponse, error) {
	respv := &EntityResponse{}

	url := fmt.Sprintf("http://%v/api/states/%v", s.url, req.GetButtonId())

	response, err := s.client.R().
		SetAuthToken(s.token).
		SetResult(respv).Get(url)

	if err != nil {
		return nil, fmt.Errorf("unable to read (%v): %w", url, err)
	}

	if response.Status() != "200 OK" {
		return nil, status.Errorf(codes.Internal, "Unable to read: %v", response.Status())
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
