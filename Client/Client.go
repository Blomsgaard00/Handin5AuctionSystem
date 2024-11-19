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
	client := proto.AuctionClient(nil)
	ports := [3]string{"localhost:5101", "localhost:5102", "localhost:5103"}

	for _, port := range ports{
		conn, err := grpc.NewClient(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
		client = proto.NewAuctionClient(conn)
	}
	
	
	
}

