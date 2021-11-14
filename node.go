package main

import (
	"bufio"
	"context"
	"log"
	"net"
	"os"
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

type NodeServer struct{
	proto.UnimplementedTokenRingServer
	stream proto.TokenRing_RecieveTokenServer
}

func recieve(s *NodeServer){
	ctx := context.Background()
	stream, _ := client.RecieveToken(ctx, &proto.Empty{})

	for{
		token, _ := stream.Recv()
		s.passToken(ctx, token)
		continue
	}
}
func (s *NodeServer) recieveToken(ctx context.Context, stream proto.TokenRing_RecieveTokenServer)  error  {
	return nil
}
func (s *NodeServer) passToken(ctx context.Context, token *proto.Token) (*proto.Empty, error){
	log.Printf("Recieved token from port: %s", "xxxx")
	log.Print("Entering critical secrtion with the token")
	time.Sleep(3*time.Second)
	log.Print("Exiting critical zone")
	log.Printf("Passing the token to port: %s", "xxxx")
	client.PassToken(ctx, token)
	return &proto.Empty{}, nil
}

func main() {
	log.Print("Hey")

	// init port and traget port - start off the service with flag
	reader := bufio.NewReader(os.Stdin)

	temp, _ := reader.ReadString('\n')
	temp = strings.TrimSpace(temp)

	temp2, _ := reader.ReadString('\n')
	temp2 = strings.TrimSpace(temp)

	temp3, _ := reader.ReadString('\n')
	temp3 = strings.TrimSpace(temp)

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

	// Serve incomming conenctions to the listener
	grpcServer.Serve(listener)

	// Establish a connection to the next node in the ring
	conn, _ := grpc.Dial(tP, grpc.WithInsecure())
	// Close the conneciton when the method exits
	defer conn.Close()

	client = proto.NewTokenRingClient(conn)
	
	token := &proto.Token{Id: 1}

	if (cmd == "\\start"){
		client.PassToken(context.Background(), token)
	}else if( cmd == "\\recieve"){
		recieve(server)
	}
}