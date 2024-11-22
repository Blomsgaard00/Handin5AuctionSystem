// publish method, 128 max characters
package main

import (
	proto "Handin5AuctionSystem/gRPC"
	"bufio"
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var timestamp int32
var ports [3]string
var clientID string

func main() {
	//if invalid cli error is displayed
	if len(os.Args) < 2 || len(os.Args) > 2 {
		log.Fatal("error making client, run: go run client.go <client.ID> ")
	}
	clientID = os.Args[1]

	client := proto.AuctionClient(nil)
	ports = [3]string{"localhost:5101", "localhost:5102", "localhost:5103"}
	log.Println(clientID + " has now joined the auction and you have 2 options:")
	log.Println("write 'bid' to enter a bid in the auction")
	log.Println("write 'get result' to get the auction result")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(os.Stdin)
		for {
			text, _ := reader.ReadString('\n')
			// convert CRLF to LF
			text = strings.Replace(text, "\n", "", -1)
			if strings.Compare("bid", text) == 0 {
				readbid(client)
			} else if strings.Compare("get result", text) == 0 || strings.Compare("getresult", text) == 0 {
				//getresult()
			} else {
				log.Println("invalid command. You have two options: ")
				log.Println("'bid' to enter a bid in the auction")
				log.Println("'get result' to get the auction result")
			}
		}
	}()

	wg.Wait()
}

func readbid(client proto.AuctionClient) {
	reader := bufio.NewReader(os.Stdin)
	//var amount int32
	log.Print("Enter bid amount: ")
	out, _ := reader.ReadString('\n')
	out = strings.Replace(out, "\n", "", -1)
	amount, _ := strconv.Atoi(out)
	responses := [3]*proto.Ack{}

	for _, port := range ports {
		conn, err := grpc.NewClient(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Error creating the server %v", err)
		}
		currentbid := &proto.Bid{
			Amount:   int32(amount),
			Clientid: clientID,
		}

		client = proto.NewAuctionClient(conn)
		respons, err := client.Bidding(context.Background(), currentbid)
		if err != nil {
			log.Fatalf("Error creating the server %v", err)
		}
		responses[1] = respons
	}

}

/*//set client name by reading cli command


}
*/
