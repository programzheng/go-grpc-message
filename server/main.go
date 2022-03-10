package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	pb "github.com/programzheng/go-grpc-message"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement message.GreeterServer.
type Server struct {
	pb.UnimplementedGreeterServer
	mu         sync.Mutex
	routeNotes map[string][]*pb.RouteNote
}

// SayMessage implements message.GreeterServer
func (s *Server) SayMessage(ctx context.Context, in *pb.MessageRequest) (*pb.MessageResponse, error) {
	log.Printf("Received: %v", in.GetMessage())
	return &pb.MessageResponse{Message: in.GetMessage()}, nil
}

func (s *Server) RouteChat(stream pb.Greeter_RouteChatServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		key := serialize(in.Location)

		s.mu.Lock()
		s.routeNotes[key] = append(s.routeNotes[key], in)
		rn := make([]*pb.RouteNote, len(s.routeNotes[key]))
		copy(rn, s.routeNotes[key])
		s.mu.Unlock()

		for _, note := range rn {
			log.Printf("RouteNote: %v", note)
			if err := stream.Send(note); err != nil {
				return err
			}
		}
	}
}

func serialize(point *pb.Point) string {
	return fmt.Sprintf("%d %d", point.Latitude, point.Longitude)
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &Server{routeNotes: make(map[string][]*pb.RouteNote)})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
