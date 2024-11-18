package main

import (
	proto "Handin5AuctionSystem/gRPC"
	"context"
	//"fmt"
	"log"
	"net"
	//"sync"

	"google.golang.org/grpc"
)

type AuctionServer struct {
	proto.UnimplementedAuctionServer
	//stream?
	//active bool
	currentWinner string
	highestBid    int32
	clients map[string]bool
	error chan error
}

// bididng
func (s *AuctionServer) Bidding(ctx context.Context, msg *proto.Bid) (*proto.Ack) {
	if (s.highestBid < msg.Bid){
		s.highestBid = msg.Bid
		return "bid valid"
	} else{
		return "bid invalid"
	}
}

// ctx context.Context, msg *proto.Message) (*proto.Close, error
// result
func (s *AuctionServer) GetResult(ctx context.Context, msg *proto.Empty) (*proto.Result) {
	return s.highestBid //, s.currentWinner
}

func main() {

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	//create instance

	// Register the pool with the gRPC server
	proto.RegisterAuctionServer(grpcServer, AuctionServer)

	// Create a TCP listener at port 5103
	listener, err := net.Listen("tcp", ":5103")

	if err != nil {
		log.Fatalf("Error creating the server %v", err)
	}

	log.Println("Server started at port :5103")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error creating the server %v", err)
	}
	/*
		type Pool struct {
			clientCount int32
			proto.UnimplementedAuctionServer
			Connection      []*Connection
			serverTimestamp int32
		}

		func (p *Pool) CreateStream(pconn *proto.Connect, stream proto.ChittyChat_CreateStreamServer) error {
			p.clientCount++

			ClientID := p.clientCount
			conn := &Connection{
				//stream: stream,
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
		}*/
	/*
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
	   }*/
}
