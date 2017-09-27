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
	mostSeen := make(chan uint64)
	go func() { mostSeen <- cat.MostSeen(user1, byUser2) }()
	bf := make(chan uint64)
	go func() { bf <- cat.BestFriend(user1, byUser2) }()
	cr := make(chan uint64)
	go func() { cr <- cat.Crush(user1, byUser2) }()
	ms := <-mostSeen
	ml7 := make(chan uint64)
	go func() { ml7 <- cat.MutualLove_7Days(cql, user1, ms) }()
	mlg := make(chan uint64)
	go func() { mlg <- cat.MutualLoveGlobal(cql, user1) }()

	return &data.CategoriesReply{
		MostSeen:         ms,
		BestFriend:       <-bf,
		Crush:            <-cr,
		MutualLove_7Days: <-ml7,
		MutualLoveGlobal: <-mlg,
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
