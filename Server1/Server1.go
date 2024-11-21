package main

import (
	proto "Handin5AuctionSystem/gRPC"
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type AuctionServer struct {
	proto.UnimplementedAuctionServer
	auctionOpen bool
	highestBid  proto.Bid
	timestamp   int32
}

// bididng
func (s *AuctionServer) Bidding(ctx context.Context, currentBid *proto.Bid) (*proto.Ack, error) {
	s.timestamp = compareTimestamps(s.timestamp, currentBid.Timestamp)
	s.timestamp++
	if s.highestBid.Amount < currentBid.Amount && s.auctionOpen {
		s.highestBid.Amount = currentBid.Amount
		s.highestBid.Clientid = currentBid.Clientid
		acknowledgement := &proto.Ack{
			BidAccepted: "Success: Bid accepted",
			Timestamp:   s.timestamp,
		}
		return acknowledgement, nil
	} else {
		acknowledgement := &proto.Ack{
			BidAccepted: "Fail: Bid not accepted",
			Timestamp:   s.timestamp,
		}
		return acknowledgement, nil
	}

}

// ctx context.Context, msg *proto.Message) (*proto.Close, error
// result
func (s *AuctionServer) GetResult(ctx context.Context, empty *proto.Empty) (*proto.Result, error) {
	s.timestamp++
	if s.auctionOpen {
		result := &proto.Result{
			Result:    "Current highest bid is " + fmt.Sprint(s.highestBid.Amount) + " by client: " + s.highestBid.Clientid,
			Timestamp: s.timestamp,
		}
		return result, nil //, s.currentWinner
	} else {
		result := &proto.Result{
			Result:    "The auction is over and the winner was client " + s.highestBid.Clientid + " with the bid " + fmt.Sprint(s.highestBid.Amount),
			Timestamp: s.timestamp,
		}
		return result, nil //, s.currentWinner
	}
}

func main() {

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	server := &AuctionServer{
		auctionOpen: true,
		timestamp:   0,
	}

	// Register the pool with the gRPC server
	proto.RegisterAuctionServer(grpcServer, server)

	// Create a TCP listener at port 5101
	listener, err := net.Listen("tcp", ":5101")
	if err != nil {
		log.Fatalf("Error creating the server %v", err)
	}

	log.Println("Server started at port :5101")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error creating the server %v", err)
	}
	server.timestamp++

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		auctionTime(server)
	}()
	wg.Wait()

}

func compareTimestamps(timestampServer int32, timeStampClient int32) int32 {
	highestTimestamp := timestampServer

	if timeStampClient > highestTimestamp {
		highestTimestamp = timeStampClient
	}

	return highestTimestamp

}

func auctionTime(s *AuctionServer) {

	for {
		if s.timestamp > 20 {
			s.auctionOpen = false
			break
		}
	}

}
