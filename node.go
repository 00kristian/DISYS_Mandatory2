package main

import (
	"bufio"
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/00kristian/DISYS_Mandatory2/tree/docker/proto"
	"google.golang.org/grpc"
)

// Target port for the next node in the ring
var targetPort string

// Port for this node
var port string

// Client for this node
var client proto.TokenRingClient

// Boolean indicating starter node
var starterNode bool
type Connection struct{
	stream proto.TokenRing_ListenServer
}
type NodeServer struct {
	proto.UnimplementedTokenRingServer
	conn *Connection
}

func main() {
	// init port and traget port - start off the service with flag
	reader := bufio.NewReader(os.Stdin)

	starter, _ := reader.ReadString('\n')
	starter = strings.TrimSpace(starter)
	if starter == "Y" {
		starterNode = true
	} else {
		starterNode = false
	}

	temp, _ := reader.ReadString('\n')
	temp = strings.TrimSpace(temp)

	temp2, _ := reader.ReadString('\n')
	temp2 = strings.TrimSpace(temp2)

	temp3, _ := reader.ReadString('\n')
	temp3 = strings.TrimSpace(temp3)

	port = temp
	targetPort = temp2

	start(port, targetPort, temp3)
}

func start(p, tP, cmd string) {
	// Start up of grpc server
	grpcServer := grpc.NewServer()

	// Create a listener for a specific port - in this case the nodes own port
	listener, err := net.Listen("tcp", p)

	// Check if error occurred when trying to listen on port
	if err != nil {
		log.Fatalf("Error listining on port: %v", err)
	}
	server := &NodeServer{}
	// Register this node's server to the grpc
	proto.RegisterTokenRingServer(grpcServer, server)
	
	go func() {
		// Establish a connection to the next node in the ring
		conn, _ := grpc.Dial(tP, grpc.WithInsecure())
		// Close the conneciton when the method exits
		defer conn.Close()
		
		client = proto.NewTokenRingClient(conn)
		
		stream, err := client.Listen(context.Background(), &proto.Empty{})
		if err != nil {
			log.Fatalf("Error getting stream %v", err)
		}
		log.Printf("HELLO %s, %s, %s",p,tP, strconv.FormatBool(starterNode))
		for{
			if starterNode{
				starterToken := &proto.Token{Id: 1}
				log.Print("Starting the token ring")
				log.Print("Entering critical secrtion with the token")
				time.Sleep(3 * time.Second)
				log.Print("Exiting critical zone")
				log.Printf("Passing the token to port: %s", tP)
				log.Printf("Token id: %d", starterToken.Id)
				server.PassToken(context.Background(), starterToken)
				starterNode = false
				continue
			}
			log.Printf("HELLO FROM INSIDE %s",strconv.FormatBool(starterNode))
			
			recieved, err := stream.Recv()
			if err != nil {
				log.Fatalf("Error in recieving message: %v", err)
			}
			log.Printf("Recieved token from port: %s", "xxxx")
			log.Print("Entering critical secrtion with the token")
			time.Sleep(3 * time.Second)
			log.Print("Exiting critical zone")
			log.Printf("Passing the token to port: %s", tP)
			server.PassToken(context.Background(), recieved)
		}
	} ()
		
		// Serve incomming conenctions to the listener
	grpcServer.Serve(listener)
}
	
	func (n *NodeServer) Listen(empty *proto.Empty, stream proto.TokenRing_ListenServer) error{
		conn := &Connection{
			stream: stream,
	}
	n.conn = conn
	return nil
}

func (n *NodeServer) PassToken(ctx context.Context, token *proto.Token) (*proto.Empty, error) {
	err := n.conn.stream.Send(token)
	if err != nil {
		log.Fatalf("Error sending token: %v", err)
	}
	return &proto.Empty{}, nil
}
