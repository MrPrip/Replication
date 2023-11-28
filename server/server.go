package main

import (
	proto "Replication/proto"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	bidders = make(map[string]*Participant)
	port = 8080
)

type Server struct {
	proto.UnimplementedReplicationServer
	Bidders            map[string]*Participant
	BackupServerPort   int64
	BackupServerClient proto.ReplicationClient
	Auction            Auction
}

type Auction struct {
	Timer        int
	IsOver       bool
	HigestBidder string
	HigestBid    int64
	mu           sync.Mutex
}

type Participant struct {
	name string
}

func main() {
	flag.Parse()

	server := &Server{
		Bidders:            bidders,
		BackupServerPort:   int64(port + 1),
		BackupServerClient: nil,
		Auction: Auction{
			Timer:        120,
			IsOver:       false,
			HigestBidder: "",
			HigestBid:    0,
		},
	}

	grpcServer := grpc.NewServer()

	// Nessecery for while loop
	var listener net.Listener
	var err error
	for {
		listener, err = net.Listen("tcp", ":"+strconv.Itoa(int(port)))
		if err != nil {
			port++
			continue
		}
		break
	}
	logFilePath := "../programLog.txt"

	_, logerr := os.Stat(logFilePath)
	if os.IsNotExist(logerr) {
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

	time.Sleep(5 * time.Second)

	conn, err := grpc.Dial("localhost:"+strconv.Itoa(int(server.BackupServerPort)), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("No backup server started at port: %d", server.BackupServerPort)
		fmt.Printf("No backup server started at port: %d", server.BackupServerPort)
	} else {
		server.BackupServerClient = proto.NewReplicationClient(conn)
	}

	log.Printf("Started server at port: %d\n", port)
	fmt.Printf("Started server at port: %d\n", port)

	// Simulate crash if server is running on port 8080 (main server)
	if port == 8080 {
		go server.die()
	}

	proto.RegisterReplicationServer(grpcServer, server)
	grpcServer.Serve(listener)
}

func (s *Server) Bid(bid *proto.PlaceBid, stream proto.Replication_BidServer) error {
	//s.queue
	/*bid.lamportTime > s.lamportTime {
		s.queue.end += bid
		s.lamportTime = bid.lamportTime
		break
	} else {
		for item t : s.queue {
			!(bid.lamportTime < t.lamportTime) {
				s.queue.here +=
			}
		}
		break
	}*/

	//s.Auction.mu.Lock()
	
	// start auction timer if first bid is made
	if s.Auction.Timer == 120 {
		go func() {
			for s.Auction.Timer >= 0 {
				time.Sleep(1 * time.Second)
				s.Auction.Timer = s.Auction.Timer - 1
			}
			s.Auction.IsOver = true
		}()
	}

	// Register user when they bid
	if !checkIfParticipantIsRegisered(s.Bidders, bid.ClientName) {
		participant := &Participant{
			name: bid.ClientName,
		}
		s.Bidders[bid.ClientName] = participant
	}

	acknowledgment := &proto.Acknowledgement{
		AcknowledgementMessage: "Hi " + bid.ClientName + ". Your bid was lower than the winning bid. The winning bid is currently at $" + strconv.Itoa(int(s.Auction.HigestBid)),
	}

	if s.Auction.IsOver {
		acknowledgment.AcknowledgementMessage = "The auction is over. Type 'result' to view the winner."
	} else {
		currentBid := bid.BidAmount

		if currentBid > s.Auction.HigestBid {
			s.Auction.HigestBid = currentBid
			s.Auction.HigestBidder = bid.ClientName
			acknowledgment.AcknowledgementMessage = "Hi " + bid.ClientName + " the bid was placed with the auction"
			s.BackupServerClient.Bid(context.Background(), &proto.PlaceBid{
				ClientName: bid.ClientName,
				BidAmount:  bid.BidAmount,
			})
			log.Println("A bid of $" + strconv.Itoa(int(s.Auction.HigestBid)) + " was placed by " + s.Auction.HigestBidder + " at server " + strconv.Itoa(int(port)))
		}
	}
	//s.Auction.mu.Unlock()
	// Send acknowledgment back to the client
	if err := stream.Send(acknowledgment); err != nil {
		log.Printf("Error sending acknowledgment: %v", err)
		fmt.Printf("Error sending acknowledgment: %v", err)
	}
	return nil
}

func (s *Server) Result(close *proto.Close, stream proto.Replication_ResultServer) error {
	//s.Auction.mu.Lock()
	 

	outcome := &proto.Outcome{
		RepyMessage: "The winner of the auction is " + s.Auction.HigestBidder,
	}

	if !s.Auction.IsOver {
		outcome.RepyMessage = "The higest bid is currently $" + strconv.Itoa(int(s.Auction.HigestBid))
	}
	//s.Auction.mu.Unlock()
	// Send acknowledgment back to the client
	if err := stream.Send(outcome); err != nil {
		log.Printf("Error sending acknowledgment: %v", err)
		fmt.Printf("Error sending acknowledgment: %v", err)
	}
	return nil
}

// Check if client is registered with the auction
func checkIfParticipantIsRegisered(list map[string]*Participant, clientName string) bool {
	for _, participant := range list {
		if participant.name == clientName {
			return true
		}
	}
	return false
}

// Make the server die after 30 seconds
func (s *Server) die() {
	log.Println("Killing primary server (" + strconv.Itoa(int(port)) + ") after 30 seconds")
	fmt.Println("Killing primary server (" + strconv.Itoa(int(port)) + ") after 30 seconds")
	time.Sleep(30 * time.Second)
	log.Println("Killing primary server (" + strconv.Itoa(int(port)) + ") now")
	fmt.Println("Killing primary server (" + strconv.Itoa(int(port)) + ") now")
	os.Exit(0)
}
