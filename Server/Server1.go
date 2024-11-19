package main

import (
	proto "Handin5AuctionSystem/gRPC"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type AuctionServer1 struct {
	proto.UnimplementedAuctionServer
	//stream?
	//active bool
	highestBid proto.Bid
}

// bididng
func (s *AuctionServer1) Bidding(ctx context.Context, msg *proto.Bid) (*proto.Ack, error) {
	if (s.highestBid < msg.Bid){
		s.highestBid = msg.Bid
		s.highestBid = msg.Clientid
		return "bid valid", nil
	} else{
		return "bid invalid", nil
	}
}

// ctx context.Context, msg *proto.Message) (*proto.Close, error
// result
func (s *AuctionServer1) GetResult(ctx context.Context, msg *proto.Empty) (*proto.Result, error) {
	return s.highestBid, nil//, s.currentWinner
}

func main() {

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	//create instance

	// Register the pool with the gRPC server
	proto.RegisterAuctionServer1(grpcServer, AuctionServer1)

	// Create a TCP listener at port 5101
	listener, err := net.Listen("tcp", ":5101")

	if err != nil {
		log.Fatalf("Error creating the server %v", err)
	}

	log.Println("Server started at port :5101")

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
