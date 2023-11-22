package main

import (
	proto "Replication/proto"
	"context"
	"flag"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port    = flag.Int64("sPort", 8080, "Centrel server port")
	bidders = make(map[string]*Participant)
)

type Server struct {
	proto.UnimplementedReplicationServer
	Bidders            map[string]*Participant
	idCounter          int64
	Port               int64
	BackupServerPort   int64
	BackupServerClient proto.ReplicationClient
	Auction            Auction
}

type Auction struct {
	Timer        int
	IsOver       bool
	HigestBidder string
	HigestBid    int64
	Winner       int64
}

type Participant struct {
	id   int64
	name string
}

func main() {
	flag.Parse()

	server := &Server{
		Bidders:            bidders,
		idCounter:          0,
		Port:               *port,
		BackupServerPort:   *port + 1,
		BackupServerClient: nil,
		Auction: Auction{
			Timer:        120,
			IsOver:       false,
			HigestBidder: "",
			HigestBid:    0,
			Winner:       0,
		},
	}

	grpcServer := grpc.NewServer()

	// Nessecery for while loop
	var listener net.Listener
	var err error
	for {
		listener, err = net.Listen("tcp", ":"+strconv.Itoa(int(server.Port)))
		if err != nil {
			server.Port++
			continue
		}
		break
	}

	time.Sleep(5 * time.Second)

	conn, err := grpc.Dial("localhost:"+strconv.Itoa(int(server.BackupServerPort)), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("No packup server started at port: %d", server.BackupServerPort)
	} else {
		server.BackupServerClient = proto.NewReplicationClient(conn)
	}

	log.Printf("Started server at port: %d\n", server.Port)

	// Simulate crash if server is running on port 8080 (main server)
	if server.Port == 8080 {
		go server.die()
	}

	proto.RegisterReplicationServer(grpcServer, server)
	grpcServer.Serve(listener)
}

func (s *Server) GetIdFromServer(context context.Context, close *proto.Close) (*proto.User, error) {
	user := &proto.User{ClientId: s.idCounter}
	s.idCounter++
	return user, nil
}

func (s *Server) Bid(bid *proto.PlaceBid, stream proto.Replication_BidServer) error {
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
			id:   bid.ClientID,
			name: bid.ClientName,
		}
		s.Bidders[bid.ClientName] = participant
	}


	acknowledgment := &proto.Acknowledgement{
		AcknowledgementMessage: "Your bid was lower than the winning bid. The winning bid is currently at $" + strconv.Itoa(int(s.Auction.HigestBid)),
	}

	if s.Auction.IsOver {
		acknowledgment.AcknowledgementMessage = "The auction is over. Type 'result' to view the winner."
	} else {
		currentBid := bid.BidAmount

		if currentBid > s.Auction.HigestBid {
			s.Auction.HigestBid = currentBid
			s.Auction.HigestBidder = bid.ClientName
			acknowledgment.AcknowledgementMessage = "The bid was placed with the auction " + strconv.Itoa(int(bid.ClientID))
			s.BackupServerClient.Bid(context.Background(), &proto.PlaceBid{
				ClientID:  bid.ClientID,
				BidAmount: bid.BidAmount,
			})
		}
	}

	// Send acknowledgment back to the client
	if err := stream.Send(acknowledgment); err != nil {
		log.Printf("Error sending acknowledgment: %v", err)
	}
	return nil
}

func (s *Server) Result(close *proto.Close, stream proto.Replication_ResultServer) error {
	outcome := &proto.Outcome{
		RepyMessage: "The winner of the auction is " + s.Auction.HigestBidder,
	}

	if !s.Auction.IsOver {
		outcome.RepyMessage = "The higest bid is currently $" + strconv.Itoa(int(s.Auction.HigestBid))
	}

	// Send acknowledgment back to the client
	if err := stream.Send(outcome); err != nil {
		log.Printf("Error sending acknowledgment: %v", err)
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
	time.Sleep(30 * time.Second)
	os.Exit(0)
}
