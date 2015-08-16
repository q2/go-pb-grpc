package main

import (
	"io"
	"log"
	"os"
	"strconv"
	"time"

	pb "github.com/q2/go-pb-grpc/message"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address        = "localhost:50051"
	defaultMessage = "message_"
)

func main() {
	conn, err := grpc.Dial(address)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewMessageServiceClient(conn)

	message := defaultMessage + "simple"
	if len(os.Args) > 1 {
		message = os.Args[1]
	}

	log.Println("test sending a simple message and get a response...")
	r, err := client.SendMessageSimple(context.Background(), &pb.MessageRequest{Message: message})
	if err != nil {
		log.Fatalf("could not send message: %v", err)
	}
	log.Printf("sent: %v", message)
	msg := r.Message
	log.Printf("response: %v", msg)

	log.Println("starting stream for 100k msgs in 5 sec...")
	time.Sleep(5 * time.Second)

	stream, err := client.SendMessage(context.Background())
	waitc := make(chan struct{})

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("failed to receive message: %v", err)
			}
			log.Printf("receiving: %s", in.Message)
		}
	}()

	for i := 1; i <= 100000; i++ {
		num := strconv.Itoa(i)
		msg := defaultMessage + num
		log.Printf("sending: %s", msg)
		if err := stream.Send(&pb.MessageRequest{Message: msg}); err != nil {
			log.Fatalf("failed to send message: %v", err)
		}
	}

	stream.CloseSend()
	<-waitc
}
