package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	//"time"
	"strconv"
	proto "Replication/proto"
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
}

func main() {
	flag.Parse()

	var client proto.ReplicationClient
	var recivedId *proto.User

	client, recivedId = connectToServer(*port)
	

	clientStruct := &Client {
		id: recivedId.ClientId,
		clientName: *name,
	}
	
	
	log.Println("You've entered an auction: To start bidding enter a bid. If you wish to exit type 'exit'.");
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		tempMessageHolder := scanner.Text()
		if tempMessageHolder == "exit" {
			log.Printf("You have left the auction. \n Goodbye!")
			os.Exit(1)
		} else if tempMessageHolder == "result" {
			var err error
			err = getResult(*clientStruct, client)
			for err != nil {
				client, _ = connectToServer(*port+1)
				err = getResult(*clientStruct, client)
			}
		} else {
			// check if tempMessageHolder is a number
			returnNumber := isNumeric(tempMessageHolder)
			if returnNumber > 0 {
				var reply proto.Replication_BidClient
				var err error
				reply, err = sendBid(*clientStruct, client, returnNumber)
				for err != nil {
					client, _ = connectToServer(*port+1)
					reply, err = sendBid(*clientStruct, client, returnNumber)
				}
				replyMessage, _ := reply.Recv()
				fmt.Println(replyMessage.AcknowledgementMessage)
			} else {
				fmt.Println("You must enter a number to place a bid")
			}
		}
	}
}

func connectToServer(port int) (proto.ReplicationClient, *proto.User) {
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port %d", port)
	}

	client := proto.NewReplicationClient(conn)

	recivedId, errGetID := client.GetIdFromServer(context.Background(), &proto.Close{})

	if errGetID != nil {
		client, recivedId := connectToServer(port+1)	
		return client, recivedId
	}	

	return client, recivedId
}

func getResult(clientStruct Client, client proto.ReplicationClient) error {
	reply, err := client.Result(context.Background(), &proto.Close{})

	if err != nil {
		return err
	}
	replyMessage, _ := reply.Recv()
	fmt.Println(replyMessage.RepyMessage)
	return  nil
}

func sendBid(clientStruct Client, client proto.ReplicationClient, bidAmount int64) (proto.Replication_BidClient, error){
	reply, err := client.Bid(context.Background(), &proto.PlaceBid{
		ClientID: clientStruct.id,
		ClientName: clientStruct.clientName,
		BidAmount: bidAmount,
	})

	if err != nil {
		return reply, err
	}
	return reply, nil
}

func isNumeric(stringToCheck string) int64 {
	number, err := strconv.Atoi(stringToCheck)
	if err != nil {
		return -1
	}
	return int64(number)
}

