package main

import (
	"log"
	"net"

	cat "github.com/joune/zenly/categories"
	"github.com/joune/zenly/data"
	"github.com/joune/zenly/db"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

var (
	cql = db.InitDB()
)

type server struct{}

func (s *server) GetCategories(ctx context.Context, in *data.CategoriesRequest) (*data.CategoriesReply, error) {
	user1 := in.GetUser()
	byUser2, err := db.FetchSessionsGrouped(cql, user1, 7)
	if err != nil {
		//don't report errors to client, just yield empty result (wdyt?)
		return &data.CategoriesReply{}, nil
	}
	mostSeen := cat.MostSeen(user1, byUser2)
	return &data.CategoriesReply{
		MostSeen:         mostSeen,
		BestFriend:       cat.BestFriend(user1, byUser2),
		Crush:            cat.Crush(user1, byUser2),
		MutualLove_7Days: cat.MutualLove_7Days(cql, user1, mostSeen),
		MutualLoveGlobal: cat.MutualLoveGlobal(cql, user1),
	}, nil
}

func main() {
	defer cql.Close()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	data.RegisterCategoriesServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
