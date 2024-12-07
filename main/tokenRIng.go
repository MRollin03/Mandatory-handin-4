package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"
	"token-ring/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Node struct {
	pb.UnimplementedTokenRingServer
	id          int32
	hasToken    bool
	nextAddress string
	requested   bool 
	mu          sync.Mutex
	server      *grpc.Server
}

var nodes []*Node
var numOfNodes = 5

// startServers initializes and starts the gRPC servers for each node.
func startServers(ctx context.Context) {
	startPort := 5001
	for i := 0; i < numOfNodes; i++ {
			nodes[i] = &Node{
				id:          int32(i + 1),
				hasToken:    i == 0, // The first node starts with the token
				nextAddress: fmt.Sprintf("localhost:%d", startPort+((i+1)%numOfNodes)),
			}
			

		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", startPort+i))
		if err != nil {
			log.Fatalf("Failed to listen on port %d: %v", startPort+i, err)
		}

		grpcServer := grpc.NewServer()
		pb.RegisterTokenRingServer(grpcServer, nodes[i])
		nodes[i].server = grpcServer

		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				log.Fatalf("Failed to serve: %v", err)
			}
		}()
	}

	// Token passing will start when the server is up and running
	nodes[0].passToken()
}

// ReceiveToken processes the token and allows a node to enter the critical section if requested.
func (n *Node) ReceiveToken(ctx context.Context, req *pb.ReceiveMessageRequest) (*pb.AckMessage, error) {
	log.Printf("Node %d has received the token", n.id)
	n.hasToken = true;
	time.Sleep(2 * time.Second);

	n.mu.Lock()
	defer n.mu.Unlock()

	// Check if this node has requested access to the critical section
	if n.requested {
		n.enterCriticalSection()
	}

	// Pass the token to the next node
	n.passToken()

	return &pb.AckMessage{Message: "Token received and processed"}, nil
}

// passToken sends the token to the next node in the ring.
func (n *Node) passToken() {
	log.Printf("Node %d is passing the token", n.id)

	// Use a fresh context to avoid cancellation issues
	time.Sleep(4 * time.Second); // ---------------- PUT IN TO GIVE THE TERMINAL USER A POSSIBLE WRITING REQUEST OPPORTUNITY --------------------------------------------------------------

	// Pass the token asynchronously to avoid blocking the current node
	go func() {
		conn, err := grpc.Dial(n.nextAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("Failed to connect to next node: %v", err)
			return
		}
		defer conn.Close()

		client := pb.NewTokenRingClient(conn)
		_, err = client.ReceiveToken(context.Background(), &pb.ReceiveMessageRequest{})
		if err != nil {
			log.Printf("Failed to pass token: %v", err)
			return
		}

		// Only set hasToken to false after passing the token
		n.hasToken = false
		log.Printf("Node %d passed the token to node %d", n.id, (int(n.id))%numOfNodes+1)
	}()
}


// enterCriticalSection simulates work done in the critical section
func (n *Node) enterCriticalSection() {
	if(n.hasToken){
		
		log.Printf("Node %d is entering the critical section", n.id)
		time.Sleep(5 * time.Second) // Simulate work with a sleep
		
		// Get the current working directory
		dir, err := os.Getwd()
		if err != nil {
			log.Printf("Error getting current working directory: %v", err)
			return
		}

		// Create the file in the current working directory
		filePath := fmt.Sprintf("%s/criticalSection", dir)
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Error creating file: %v", err)
			return
		}
		defer file.Close()

		// Write to the file
		_, err = file.WriteString(fmt.Sprintf("Critical section accessed by node %d \n", n.id))
		if err != nil {
			log.Printf("Error writing to file: %v", err)
		}

		log.Printf("Node %d wrote to the file", n.id);
		n.requested = false; //remove Request after completion
	}
}

// CommandlineInput handles user input to request critical section
func CommandlineInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter command: ")
		if scanner.Scan() {
			scanText := scanner.Text()
			text := strings.Fields(scanText) // Use Fields to handle spaces correctly
			if len(text) > 1 && text[0] == "request" {
				nodeID, err := strconv.Atoi(text[1])
				if err != nil || nodeID <= 0 || nodeID > len(nodes) {
					log.Printf("Invalid node ID: %v", err)
					continue
				}

				node := nodes[nodeID-1]
				node.mu.Lock()
				node.requested = true
				node.mu.Unlock()
				log.Printf("Node %d has requested access to the critical section", nodeID)
			} else {
				log.Println("Invalid command. Use 'request <nodeID>'.")
			}
			}
		}
	}

// gracefulShutdown gracefully stops all servers when an interrupt signal is received.
func gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down servers...")

	// Gracefully stop each gRPC server
	for _, node := range nodes {
		node.server.GracefulStop() // Graceful stop the gRPC server
	}
	log.Println("All servers stopped.")
}

func main() {
	nodes = make([]*Node, numOfNodes)

	// Create context to pass around
	ctx := context.Background()

	// Start servers for all nodes
	startServers(ctx)

	// Start handling user requests in a separate goroutine
	go CommandlineInput()

	// Wait for graceful shutdown signal
	gracefulShutdown()
}
