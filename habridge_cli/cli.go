package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	pb "github.com/brotherlogic/habridge/proto"
)

func extractServer(s string) string {
	return strings.Split(strings.Split(s, "//")[1], ":")[0]
}

func dialService(service string) (*grpc.ClientConn, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting user home dir: %v\n", err)
		os.Exit(1)
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		fmt.Printf("error getting Kubernetes config: %v\n", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		fmt.Printf("error getting Kubernetes clientset: %v\n", err)
		os.Exit(1)
	}

	services, err := clientset.CoreV1().Services(service).List(context.Background(), v1.ListOptions{})
	if err != nil {
		fmt.Printf("error getting Kubernetes services: %v\n", err)
		os.Exit(1)
	}

	for _, rservice := range services.Items {
		if rservice.ObjectMeta.Name == service {
			for _, port := range rservice.Spec.Ports {
				if port.Name == "grpc" {
					return grpc.NewClient(fmt.Sprintf("%v:%v", extractServer(kubeConfig.Host), port.NodePort), grpc.WithTransportCredentials(insecure.NewCredentials()))
				}
			}
		}
	}

	return nil, fmt.Errorf("Unable to dial %v", service)
}

func main() {
	conn, err := dialService("habridge")
	if err != nil {
		log.Fatalf("Dial fail: %v", err)
	}

	haclient := pb.NewHabridgeServiceClient(conn)
	status, err := haclient.GetState(context.Background(), &pb.GetStateRequest{
		ButtonId: "device_tracker.pixel_7",
	})

	fmt.Printf("%v [%v]\n", status, err)
}
