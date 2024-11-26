// publish method, 128 max characters
package main

import (
	proto "Handin5AuctionSystem/gRPC"
	"bufio"
	"context"
	"log"
	"os"
	"reflect"
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

	timestamp = 0

	ports = [3]string{"localhost:5101", "localhost:5102", "localhost:5103"}
	log.Println(clientID + " has now joined the auction and you have 2 options:")
	log.Println("write 'bid' to enter a bid in the auction")
	log.Println("write 'get result' to get the auction result")

	timestamp++

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
				getresult(client)
			} else {
				log.Println("invalid command. You have two options: ")
				log.Println("'bid' to enter a bid in the auction")
				log.Println("'get result' to get the auction result")
			}
		}
	}()

	wg.Wait()
}

func getresult(client proto.AuctionClient) {
	timestamp++
	responses := [3]*proto.Result{}
	for index, port := range ports {
		conn, err := grpc.NewClient(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Println("The server at port " + port + " has crashed")
			continue
		}
		empty := &proto.Empty{}

		client = proto.NewAuctionClient(conn)
		respons, err := client.GetResult(context.Background(), empty)
		if err != nil {
			log.Println("The server at port " + port + " has crashed")
			continue
		}
		if timestamp > respons.Timestamp {
			respons.Timestamp = timestamp
		}
		responses[index] = respons
	}
	if reflect.DeepEqual(responses[0], responses[1]) {
		log.Println(responses[0].Result)
	} else if reflect.DeepEqual(responses[0], responses[2]) {
		log.Println(responses[0].Result)
	} else if reflect.DeepEqual(responses[1], responses[2]) {
		log.Println(responses[1].Result)
	} else {
		log.Println("Servers cant reach consensus. More than one server dont work")
	}
}

func readbid(client proto.AuctionClient) {
	timestamp++
	reader := bufio.NewReader(os.Stdin)
	//var amount int32
	log.Print("Enter bid amount: ")
	out, _ := reader.ReadString('\n')
	out = strings.Replace(out, "\n", "", -1)
	amount, _ := strconv.Atoi(out)
	responses := [3]*proto.Ack{}

	for index, port := range ports {
		conn, err := grpc.NewClient(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Println("The server at port " + port + " has crashed")
			continue
		}
		currentbid := &proto.Bid{
			Amount:    int32(amount),
			Clientid:  clientID,
			Timestamp: timestamp,
		}

		client = proto.NewAuctionClient(conn)
		respons, err := client.Bidding(context.Background(), currentbid)
		if err != nil {
			log.Println("The server at port " + port + " has crashed")
			continue
		}
		responses[index] = respons
	}
	if reflect.DeepEqual(responses[0], responses[1]) {
		log.Println(responses[0].BidAccepted)
	} else if reflect.DeepEqual(responses[0], responses[2]) {
		log.Println(responses[0].BidAccepted)
	} else if reflect.DeepEqual(responses[1], responses[2]) {
		log.Println(responses[1].BidAccepted)
	} else {
		log.Println("Servers cant reach consensus. More than one server dont work")
	}
}

func compareTimestamps(timestampServer int32, timeStampClient int32) int32 {
	highestTimestamp := timestampServer

	if timeStampClient > highestTimestamp {
		highestTimestamp = timeStampClient
	}

	return highestTimestamp

}
