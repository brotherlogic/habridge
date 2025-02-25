package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	pb "github.com/brotherlogic/habridge/proto"

	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

var (
	port        = flag.Int("port", 8080, "Server port for grpc traffic")
	metricsPort = flag.Int("metrics_port", 8081, "Metrics port")
)

type Server struct {
	token  string
	url    string
	client *resty.Client
}

func NewServer(token string, url string) *Server {
	return &Server{
		token:  token,
		client: resty.New(),
		url:    url}
}

func main() {
	flag.Parse()

	token := os.Getenv("HA_TOKEN")
	if token == "" {
		log.Fatalf("Missing HA_TOKEN")
	}
	url := os.Getenv("HA_URL")
	if token == "" {
		log.Fatalf("Missing HA_URL")
	}
	s := NewServer(token, url)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("gramophile is unable to listen on the grpc port %v: %v", *port, err)
	}
	gs := grpc.NewServer()
	pb.RegisterHabridgeServiceServer(gs, s)
	go func() {
		if err := gs.Serve(lis); err != nil {
			log.Fatalf("gramophile is unable to serve grpc: %v", err)
		}
		log.Fatalf("gramophile has closed the grpc port for some reason")
	}()

	http.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(fmt.Sprintf(":%v", *metricsPort), nil)
	if err != nil {
		log.Fatalf("gramophile is unable to serve metrics: %v", err)
	}
	log.Printf("Exiting after safe shutdown")
}
