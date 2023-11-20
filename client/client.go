package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	//"net"
	proto "Replication/proto"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port = flag.Int("server", 8080, "Port for server")
	name = flag.String("name", "John Doe", "Client name")
)

type Client struct {
	id int64
	clientName string
	timestamp int64
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial("localhost:"+strconv.Itoa(*port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port %d", *port)
	}

	client := proto.NewReplicationClient(conn)
	recivedId, _ := client.GetIdFromServer(context.Background(), &proto.Close{})

	clientStruct := &Client {
		id: recivedId.ClientId,
		clientName: *name,
		timestamp: 0,
	}
	stream, _ := client.ConnectToServer(context.Background(), &proto.User{ClientId: clientStruct.id,})

	go func(str proto.Replication_ConnectToServerClient) {
		for {
			msg, _ := str.Recv()
			fmt.Printf("%s", msg.Content)
		}
	}(stream)
	

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		tempMessageHolder := scanner.Text()
		if tempMessageHolder == "exit" {
			log.Printf("You've left the auction. \n Goodbye!")
			os.Exit(1)
		}

		// check if tempMessageHolder is a number
		returnNumber, isNumber := isNumeric(tempMessageHolder)
		if isNumber {
			reply, err := client.Bid(context.Background(), &proto.PlaceBid{
				ClientID: clientStruct.id,
				BidAmount: returnNumber,
				Timestamp: clientStruct.timestamp,
			})
			if err != nil {
				fmt.Println("There was an error placing the bid")
			}
			t, _ := reply.Recv()
			
			if !t.AcknowledgementMessage {
				fmt.Println("Your bid was lower than the winning bid")
			} else {
				fmt.Println("The bid was placed with the auction")
			}
		} else {
			fmt.Println("You must enter a number to place a bid")
		}
	}
}

func isNumeric(stringToCheck string) (int64, bool) {
number, err := strconv.Atoi(stringToCheck)
	if err != nil {
		return 0, err == nil
	}
	return int64(number), err == nil
}


// Client 1 -> 8080
// 8080 modtager
// 8080 -> backup
// 8080 -> modtaget til Client 1
// Client 1 !-> modtaget = timeout

