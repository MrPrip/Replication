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
)

var (
	port    = flag.Int64("sPort", 8080, "Centrel server port")
	bidders []*Participant
)

type Server struct {
	proto.UnimplementedReplicationServer
	Bidders     []*Participant
	Port        int64
	LamportTime int64
	Auction     Auction
}

type Auction struct {
	Timer        int
	IsOver       bool
	HigestBidder int64
	HigestBid    int64
	Winner       int64
}

type Participant struct {
	stream proto.Replication_ConnectToServerServer
	id     int64
	active bool
	error  chan error
}

func main() {
	flag.Parse()

	server := &Server{
		Bidders:     bidders,
		Port:        *port,
		LamportTime: 0,
		Auction: Auction{
			Timer:        120,
			IsOver:       false,
			HigestBidder: 0,
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

	log.Printf("Started server at port: %d\n", server.Port)

	// Simulate crash one minute after start
	//if server.Port == 8080 {
	//	go server.die()
	//}

	proto.RegisterReplicationServer(grpcServer, server)
	grpcServer.Serve(listener)
}

func (s *Server) GetIdFromServer(context context.Context, close *proto.Close) (*proto.User, error){
	return &proto.User{ClientId: int64(len(s.Bidders))}, nil
}

func (s *Server) ConnectToServer(user *proto.User, stream proto.Replication_ConnectToServerServer) error {
	participant := &Participant{
		stream: stream,
		id:     user.ClientId,
		active: true,
	}

	s.Bidders = append(s.Bidders, participant)
	return <-participant.error
}

func (s *Server) Bid(bid *proto.PlaceBid, stream proto.Replication_BidServer) error {
	if s.Auction.Timer == 120 { 
		go func() {
			time.Sleep(time.Duration(s.Auction.Timer) * time.Second)
			s.Auction.IsOver = true;
		}()
	}
	acknowledgment := &proto.Acknowledgement{
		AcknowledgementMessage: "Your bid was lower than the winning bid. The winning bid is currently at $" + strconv.Itoa(int(s.Auction.HigestBid)),
		Timestamp:              0,
	}

	if s.Auction.IsOver {
		acknowledgment.AcknowledgementMessage = "The auction is over. Type 'result' to view the winner."
	} else {
		currentBid := bid.BidAmount

		if currentBid > s.Auction.HigestBid {
			s.Auction.HigestBid = currentBid
			s.Auction.HigestBidder = bid.ClientID
			acknowledgment.AcknowledgementMessage = "The bid was placed with the auction"
			log.Println(currentBid)
		}
	}

	// Send acknowledgment back to the client
	if err := stream.Send(acknowledgment); err != nil {
		log.Printf("Error sending acknowledgment: %v", err)
	}
	return nil
}

func (s *Server) Result(close *proto.Close, stream proto.Replication_ResultServer) (error) {
	outcome := &proto.Outcome{
		RepyMessage: "The winner of the auction is " + strconv.Itoa(int(s.Auction.HigestBidder)),
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

func (s *Server) die() {
	time.Sleep(5 * time.Second)
	os.Exit(0)
}
