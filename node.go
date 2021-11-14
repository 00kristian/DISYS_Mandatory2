package main

import (
	"log"
	"google.golang.org/grpc"
	"net"
)

var targetPort string
var port string

type NodeServer struct{
	UnimplementedTokenServer
}

func main() {
	log.Print("Hey")

	// init port and traget port - start off the service with flag

	start(port, targetPort)
}
 
func start(port, targetPort string) {
	// Start up of grpc server
	grpcServer := grpc.NewServer()

	// Create a listener for a specific port - in this case the nodes own port
	listener, err := net.Listen("tcp", port)
	
	// Check if error occurred when trying to listen on port
	if err != nil {
		log.Fatalf("Error listining on port: %v", err)
	}

	proto.RegisterTokenServer(grpcServer, &NodeServer{})
}