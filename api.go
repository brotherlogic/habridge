package main

import (
	"context"

	pb "github.com/brotherlogic/habridge/proto"
)

func (s *Server) GetState(ctx context.Context, req *pb.GetStateRequest) (*pb.GetStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Not yet")
}
