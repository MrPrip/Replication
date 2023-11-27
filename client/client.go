package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	proto "Replication/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	name = flag.String("name", "John Doe", "Client name")
)

const (
	port = 8080
	logFilePath = "../programLog.txt"
)

type Client struct {
	clientName string
}

func main() {
	flag.Parse()
	
	var client proto.ReplicationClient

	client = connectToServer(port)
	

	clientStruct := &Client {
		clientName: *name,
	}

	

	_, err := os.Stat(logFilePath)
	if os.IsNotExist(err) {
		// If the file does not exist, create it
		file, createErr := os.Create(logFilePath)
		if createErr != nil {
			log.Fatal("Error creating log file:", createErr)
		}
		defer file.Close()
	}
	
	// Open the log file
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	
	log.Println("You've entered an auction: To start bidding enter a bid. If you wish to exit type 'exit'.");
	fmt.Println("You've entered an auction: To start bidding enter a bid. If you wish to exit type 'exit'.");
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		tempMessageHolder := scanner.Text()
		if tempMessageHolder == "exit" {
			log.Printf("You have left the auction. \n Goodbye!")
			fmt.Printf("You have left the auction. \n Goodbye!")
			os.Exit(1)
		} else if tempMessageHolder == "result" {
			var err error
			err = getResult(*clientStruct, client)
			for err != nil {
				client = connectToServer(port+1)
				err = getResult(*clientStruct, client)
			}
		} else {
			returnNumber := isNumeric(tempMessageHolder)
			if returnNumber > 0 {
				var reply proto.Replication_BidClient
				var err error
				reply, err = sendBid(*clientStruct, client, returnNumber)
				for err != nil {
					client = connectToServer(port+1)
					reply, err = sendBid(*clientStruct, client, returnNumber)
				}
				replyMessage, _ := reply.Recv()
				log.Println(replyMessage.AcknowledgementMessage)
				fmt.Println(replyMessage.AcknowledgementMessage)
			} else {
				fmt.Println("You can either enter a number to place a bid or use the commands ('result' or 'exit')")
				fmt.Println("You can either enter a number to place a bid or use the commands ('result' or 'exit')")
			}
		}
	}
}

func connectToServer(port int) (proto.ReplicationClient) {
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port %d", port)
	}
	
	client := proto.NewReplicationClient(conn)
	
	return client
}

func getResult(clientStruct Client, client proto.ReplicationClient) error {
	reply, err := client.Result(context.Background(), &proto.Close{})

	if err != nil {
		return err
	}
	replyMessage, _ := reply.Recv()
	log.Println(replyMessage.RepyMessage)
	return  nil
}

func sendBid(clientStruct Client, client proto.ReplicationClient, bidAmount int64) (proto.Replication_BidClient, error){
	reply, err := client.Bid(context.Background(), &proto.PlaceBid{
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

