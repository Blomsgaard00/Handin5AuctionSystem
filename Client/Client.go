// publish method, 128 max characters
package main

import (
	//proto "Handin5AuctionSystem/gRPC"
	"log"
	"os"

	//"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials/insecure"
)

var timestamp int32

func main() {
	//if invalid cli error is displayed
	if len(os.Args) < 2 || len(os.Args) > 2 {
		log.Fatal("error making client, run: go run client.go <client.ID> ")
	}
	clientID := os.Args[1]
	log.Printf(clientID)
 
	client := proto.AuctionClient(nil)
	ports := [3]string{"localhost:5101", "localhost:5102", "localhost:5103"}

	for _, port := range ports{
		conn, err := grpc.NewClient(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
		client = proto.NewAuctionClient(conn)
	}
}

