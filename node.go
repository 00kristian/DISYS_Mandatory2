package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/00kristian/DISYS_Mandatory2/tree/docker/proto"
	"google.golang.org/grpc"
)

// Target port for the next node in the ring
var targetPort string

// Port for this node
var port string

var recievingPort string

// Client for this node
var client proto.TokenRingClient

// Boolean indicating starter node
var starterNode bool

type Connection struct{
	stream proto.TokenRing_ListenServer
	error chan error
}
type NodeServer struct {
	proto.UnimplementedTokenRingServer
	conn *Connection
}

func main() {
	// init port and traget port - start off the service with flag
	
	starter := os.Args[1]
	if starter == "Y" {
		starterNode = true
	} else {
		starterNode = false
	}

	port = os.Args[2]
	targetPort = os.Args[3]
	recievingPort = os.Args[4]

	log.Print(os.Args[1], os.Args[2], os.Args[3], os.Args[4])
	start(port, targetPort, recievingPort)
}

func start(p, tP, rP string) {
	// Start up of grpc server
	grpcServer := grpc.NewServer()

	// Create a listener for a specific port - in this case the nodes own port
	listener, err := net.Listen("tcp", tP)

	// Check if error occurred when trying to listen on port
	if err != nil {
		log.Fatalf("Error listining on port: %v", err)
	}
	server := &NodeServer{}
	// Register this node's server to the grpc
	proto.RegisterTokenRingServer(grpcServer, server)
	
	go func() {
		// Establish a connection to the next node in the ring
		conn, _ := grpc.Dial(p, grpc.WithInsecure())
		// Close the conneciton when the method exits
		defer conn.Close()
		
		client = proto.NewTokenRingClient(conn)
		
		stream, err := client.Listen(context.Background(), &proto.Empty{})
		if err != nil {
			log.Fatalf("Error getting stream %v", err)
		}

		for{
			if starterNode{
				starterToken := &proto.Token{Id: 1}
				log.Print("Starting the token ring")
				log.Print("Entering critical section with the token")
				time.Sleep(3 * time.Second)
				log.Print("Exiting critical section")
				log.Printf("Passing the token to port: %s", tP)
				server.PassToken(context.Background(), starterToken)
				starterNode = false
				log.Print("Waiting to recieve token...")
				continue
			}
			
			recieved, err2 := stream.Recv()
			if err != nil {
				log.Fatalf("Error in recieveing token: %v", err2)
			}
			log.Printf("Recieved token from port: %s", rP)
			log.Print("Entering critical section with the token")
			time.Sleep(3 * time.Second)
			log.Print("Exiting critical section")
			log.Printf("Passing the token to port: %s", tP)
			server.PassToken(context.Background(), recieved)
			log.Print("Waiting to recieve token...")
		}
	} ()
		// Serve incomming conenctions to the listener
	grpcServer.Serve(listener)
	// Just ensuring that method does not terminate
	for{}
}

func (n *NodeServer) Listen(empty *proto.Empty, stream proto.TokenRing_ListenServer) error{
	conn := &Connection{
		stream: stream,
		error: make(chan error),
	}
	n.conn = conn
	return <- n.conn.error
}

func (n *NodeServer) PassToken(ctx context.Context, token *proto.Token) (*proto.Empty, error) {
	err := n.conn.stream.Send(token)
	if err != nil {
		log.Fatalf("Error sending token: %v", err)
	}
	return &proto.Empty{}, nil
}
