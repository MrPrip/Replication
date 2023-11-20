package main

import (
	proto "Replication/proto"
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
}

func main() {
	flag.Parse()

	server := &Server{
		Bidders:     bidders,
		Port:        *port,
		LamportTime: 0,
		Auction: Auction{
			Timer: 120,
			IsOver: false,
			HigestBidder: 0,
			HigestBid: 0,
			Winner: 0,
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

func (s *Server) ConnectToServer(user *proto.User, stream proto.Replication_ConnectToServerServer) error {
	participant := &Participant{
		stream: stream,
		id: user.ClientId,
	}

	s.Bidders = append(s.Bidders, participant)
	return nil
}

func (s *Server) Bid(bid *proto.PlaceBid, stream proto.Replication_BidServer) error {
	acknowledgment := &proto.Acknowledgement{
		AcknowledgementMessage: false,
		Timestamp:              0,
	}
	
	currentBid := bid.BidAmount

	if currentBid > s.Auction.HigestBid {
		s.Auction.HigestBid = currentBid
		s.Auction.HigestBidder = bid.ClientID
		acknowledgment.AcknowledgementMessage = true
		s.BroadcastMessage(&proto.Message{
			Id: "",
			Content: "The higest bid is now at " + strconv.Itoa(int(currentBid)) + " $",
			LamportTime: 0,
		}, )
	}
	
	// Send acknowledgment back to the client
	if err := stream.Send(acknowledgment); err != nil {
		log.Printf("Error sending acknowledgment: %v", err)
	}
	return nil
}

func (s *Server) BroadcastMessage(message *proto.Message, stream proto.Replication_BroadcastMessageServer) error {
	for _, participant := range s.Bidders {
		participant.stream.Send(&proto.Message{
			Id: message.Id,
			Content: message.Content,
			LamportTime: message.LamportTime,
		})
	}
	return nil
}

func (s *Server) die() {
	time.Sleep(5 * time.Second)
	os.Exit(0)
}
