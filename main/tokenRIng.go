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
	isPrimary   bool
	nextAddress string
	mu          sync.Mutex
}

var numOfNodes = 6 // actually 5, but arrays are evil
var nodes []*Node

// RequestCriticalSection is used to request entry to the critical section
func (n *Node) RequestCriticalSection(ctx context.Context, req *pb.ReceiveMessageRequest) (*pb.AckMessage, error) {
	if !n.isPrimary {
		request := &pb.ReceiveMessageRequest{
			Message: &pb.ReceiveMessageRequest_Request{
				Request: &pb.RequestMessage{Nodeid: int32(n.id)},
			},
		}
		n.ReceiveToken(context.Background(), request)
	}
	return &pb.AckMessage{Message: "Request sent to primary node"}, nil
}

// ReceiveToken handles token and message processing
func (n *Node) ReceiveToken(ctx context.Context, req *pb.ReceiveMessageRequest) (*pb.AckMessage, error) {
	//n.mu.Lock()
    n.hasToken = true
	log.Printf("Node %d has received the token.", n.id-1)

	switch m := req.Message.(type) {
	case *pb.ReceiveMessageRequest_Token:
		// Token handling logic here
	case *pb.ReceiveMessageRequest_Request:
		if n.isPrimary {
            log.Printf("Primary node reached")
			// Primary node approves the request
			approval := &pb.ReceiveMessageRequest{
				Message: &pb.ReceiveMessageRequest_Ack{},
				Nodeid: m.Request.Nodeid,
			}
			n.PassTokenToNext(ctx, approval)
		} else {
			// Non-primary node just passes the request
			n.PassTokenToNext(ctx, req)
		}
	case *pb.ReceiveMessageRequest_Ack:
		if req.Nodeid-1 == n.id {
            log.Printf("requestee found")
			n.enterCriticalSection()
		} else {
            log.Printf("this was not the requestee")
			n.PassTokenToNext(ctx, req)
        }
	}

	
	//n.mu.Unlock() // Unlock after the work is done
	
	return &pb.AckMessage{Message: "Token received"}, nil
}

// passTokenToNext sends the token to the next node
func (n *Node) PassTokenToNext(ctx context.Context, msg *pb.ReceiveMessageRequest) (*pb.AckMessage, error) {
	conn, err := grpc.Dial(n.nextAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
    //log.Printf("passTokenToNext (checkmark)")
	if err != nil {
		log.Fatalf("Failed to connect to next node: %v", err)
	}
	defer conn.Close()

	client := pb.NewTokenRingClient(conn)
	_, err = client.ReceiveToken(ctx, msg)
	if err != nil {
		log.Fatalf("Error occurred while trying to pass token: %v", err)
	}

	n.hasToken = false
	// Update token state to indicate this node no longer has it
	return &pb.AckMessage{Message: "Token passed to next node"}, nil
}

// enterCriticalSection simulates work done in the critical section
func (n *Node) enterCriticalSection() {
	if n.hasToken {
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

		// Release the token after work is done
		n.hasToken = false
		
	}
}

// Set up and start the servers for each node
func startServers() {
	startPort := 5001
	for i := 0; i < numOfNodes; i++ {
		nodes[i] = &Node{
			id:          int32(i +1),
			hasToken:    false,
			isPrimary:   i == 0,
			nextAddress: "localhost:" + strconv.Itoa(startPort + (i+1)%numOfNodes), // Wrap around to next node
		}

		// Start listening on a TCP port specific to the node
		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:500%d", nodes[i].id))
		if err != nil {
			log.Fatalf("Failed to listen on port: %v", err)
		}

		grpcServer := grpc.NewServer()
		pb.RegisterTokenRingServer(grpcServer, nodes[i])

		// Start the gRPC server in a new goroutine for each node
		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				log.Fatalf("Failed to serve: %v", err)
			}
		}()
	}
}

// Gracefully shut down the server when interrupt is received
func gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit // Wait for interrupt signal
	log.Println("Server is shutting down...")

	for _, node := range nodes {
		node.mu.Lock() // Ensure no operations are happening when stopping
	}
	log.Println("All servers stopped.")
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
				sNode, err := strconv.Atoi(text[1])
				if err != nil || sNode <= 0 || sNode > len(nodes) {
					log.Printf("Invalid node ID: %v", err)
					continue
				}
				reqNode := nodes[sNode]
				_, err = reqNode.RequestCriticalSection(context.Background(), &pb.ReceiveMessageRequest{})
				if err != nil {
					log.Printf("Error requesting critical section: %v", err)
				}
			}
		}
	}
}

func main() {
	nodes = make([]*Node, numOfNodes)

	// Start the servers for each node
	startServers()

	// Start handling command line input in a separate goroutine
	go CommandlineInput()

	// Wait for the server to be interrupted and shut it down gracefully
	gracefulShutdown()
}
