package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/brotherlogic/habridge/proto"
)

func main() {
	conn, err := grpc.Dial(os.Args[1], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Dial fail: %v", err)
	}

	haclient := pb.NewHabridgeServiceClient(conn)
	status, err := haclient.GetState(context.Background(), &pb.GetStateRequest{
		ButtonId: "5ee1e8c14cbe2463f12a49829ec1415d",
	})

	fmt.Printf("%v [%v]\n", status, err)
}
