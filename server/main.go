package main

import (
	"io"
	"log"
	"net"

	pb "github.com/q2/go-pb-grpc/message"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement message.MessageServiceServer
type server struct{}

// SendMessageSimple implements message.MessageServiceServer
func (s *server) SendMessageSimple(ctx context.Context, in *pb.MessageRequest) (*pb.MessageReply, error) {
	msg := in.Message
	log.Printf("received: %v", msg)
	return &pb.MessageReply{Message: msg}, nil
}

// SendMessage implements message.MessageServiceServer as a bidirectional streaming RPC service
func (s *server) SendMessage(stream pb.MessageService_SendMessageServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		msg := in.Message
		log.Printf("received: %v", msg)
		msg = msg + "_read"
		if err := stream.Send(&pb.MessageReply{Message: msg}); err != nil {
			return err
		}
		log.Printf("sent: %v", msg)
	}
}

func main() {
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listening on: %s", port)
	s := grpc.NewServer()
	pb.RegisterMessageServiceServer(s, &server{})
	s.Serve(l)
}
