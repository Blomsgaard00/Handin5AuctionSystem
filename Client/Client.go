// publish method, 128 max characters
package main

import (
	proto "DSMandatoryActivity3TIM/gRPC"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"unicode/utf8"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var timestamp int32

func main() {
	conn, err := grpc.NewClient("localhost:5100", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Client not working")
	}
	client := proto.NewChittyChatClient(conn)
	var wg sync.WaitGroup

	
	timestamp = 0
	connectionMessage := &proto.Connect{
		Active:    true,
		Timestamp: timestamp,
	}

	stream, err := client.CreateStream(context.Background(), connectionMessage)
	if err != nil {
		log.Fatalf("Not working")
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		receiveMessage(stream)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		SendMessage(client)
	}()

	wg.Wait()
}

func receiveMessage(stream grpc.ServerStreamingClient[proto.Message]) {
	for {
		input, err := stream.Recv()
		if err != nil {
			log.Fatalf("Not working")

		}
		if input.Timestamp > timestamp {
			timestamp = input.Timestamp
		}
		log.Println("Lamport timestamp: " + fmt.Sprint(timestamp) + ", Message: " + input.Message)

	}
}

func SendMessage(client proto.ChittyChatClient) {
	reader := bufio.NewReader(os.Stdin)
	for {
		//code for reading terminal input based on https://tutorialedge.net/golang/reading-console-input-golang/
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		//checks if message is shorter that 128 charactors and a valid UTF-8 string.
		//The UTF-8 check is from a method https://henvic.dev/posts/go-utf8/
		if len(text) < 128 && utf8.ValidString(text) {
			sendMessage := &proto.Message{
				Message:   text,
				Timestamp: timestamp,
			}
			client.BroadcastMessage(context.Background(), sendMessage)
		} else {
			log.Println("invalid message")
		}

	}

}
