package main

import (
	proto "DSMandatoryActivity3TIM/gRPC"
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

// Heavily inspired by https://medium.com/@bhadange.atharv/building-a-real-time-chat-application-with-grpc-and-go-aa226937ad3c
type Connection struct {
	proto.UnimplementedChittyChatServer
	stream proto.ChittyChat_CreateStreamServer
	active bool
	error  chan error
}

type Pool struct {
	clientCount int32
	proto.UnimplementedChittyChatServer
	Connection      []*Connection
	serverTimestamp int32
}

func (p *Pool) CreateStream(pconn *proto.Connect, stream proto.ChittyChat_CreateStreamServer) error {
	p.clientCount++

	ClientID := p.clientCount
	conn := &Connection{
		stream: stream,
		active: true,
		error:  make(chan error),
	}

	initialConnectMessage := &proto.Message{
		Message:   "Participant " + fmt.Sprint(ClientID) + " joined Chitty-Chat",
		Timestamp: p.serverTimestamp,
	}

	p.BroadcastMessage(context.Background(), initialConnectMessage)

	p.Connection = append(p.Connection, conn)

	<-conn.stream.Context().Done()
	p.serverTimestamp++
	conn.active = false
	clientShutdownMessage := &proto.Message{
		Message:   "Participant " + fmt.Sprint(ClientID) + " has left the server",
		Timestamp: p.serverTimestamp,
	}

	p.BroadcastMessage(context.Background(), clientShutdownMessage)

	return <-conn.error
}

func (s *Pool) BroadcastMessage(ctx context.Context, msg *proto.Message) (*proto.Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	if s.serverTimestamp > msg.Timestamp {
		msg.Timestamp = s.serverTimestamp
	} else if msg.Timestamp > s.serverTimestamp {
		s.serverTimestamp = msg.Timestamp
	}
	s.serverTimestamp++
	msg.Timestamp++
	log.Println("Lamport timestamp: " + fmt.Sprint(s.serverTimestamp) + ", Server recieved message: " + msg.Message)
	s.serverTimestamp++
	msg.Timestamp++
	for _, conn := range s.Connection {
		wait.Add(1)

		go func(msg *proto.Message, conn *Connection) {
			defer wait.Done()

			if conn.active {
				err := conn.stream.Send(msg)

				if err != nil {
					log.Printf("Error with Stream: %v - Error: %v\n", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}
		}(msg, conn)
	}

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
	return &proto.Close{}, nil
}

func main() {
	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Create a new connection pool
	var conn []*Connection

	pool := &Pool{
		Connection: conn,
	}

	// Register the pool with the gRPC server
	proto.RegisterChittyChatServer(grpcServer, pool)

	// Create a TCP listener at port 8080
	listener, err := net.Listen("tcp", ":5100")

	if err != nil {
		log.Fatalf("Error creating the server %v", err)
	}

	log.Println("Server started at port :5100")

	// Start serving requests at port 8080
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error creating the server %v", err)
	}
}
