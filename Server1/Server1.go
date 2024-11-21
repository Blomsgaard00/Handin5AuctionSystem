package main

import (
	proto "Handin5AuctionSystem/gRPC"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type AuctionServer struct {
	proto.UnimplementedAuctionServer
	auctionOpen bool
	highestBid proto.Bid
}

// bididng
func (s *AuctionServer) Bidding(ctx context.Context, currentBid *proto.Bid) (*proto.Ack, error) {
	
	if (s.highestBid.Amount < currentBid.Amount && s.auctionOpen){
		s.highestBid.Amount = currentBid.Amount
		s.highestBid.Clientid = currentBid.Clientid
		acknowledgement := &proto.Ack{
			BidAccepted: "Success: Bid accepted",
		}
		return acknowledgement, nil
	} else{
		acknowledgement := &proto.Ack{
			BidAccepted: "Fail: Bid not accepted",
		}
		return acknowledgement, nil
	}

}

// ctx context.Context, msg *proto.Message) (*proto.Close, error
// result
func (s *AuctionServer) GetResult(ctx context.Context, empty *proto.Empty) (*proto.Result, error) {
	if(s.auctionOpen){
		result := &proto.Result{
			Result: "Current highest bid is " + fmt.Sprint(s.highestBid.Amount) + " by client: " + s.highestBid.Clientid,
		}
		return result, nil//, s.currentWinner
	}else{
		result := &proto.Result{
			Result: "The auction is over and the winner was client " + s.highestBid.Clientid + " with the bid " + fmt.Sprint(s.highestBid.Amount),
		}
		return result, nil//, s.currentWinner
	}
}

func main() {

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	server := &AuctionServer{
		auctionOpen: true,
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
}