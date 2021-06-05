package grpcserver

import (
	"apptica/commander"
	"apptica/model"

	"context"
	"errors"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	appID       = 1421444
	countryCode = 1
)

type srv struct {
	cmd commander.Commander
}

var (
	errBadDate        = errors.New("missing date")
	errFailedToGetTop = errors.New("failed to get response")
)

func (s srv) AppTopCategories(ctx context.Context, request *AppTopCategoryRequest) (*AppTopCategoryResponse, error) {
	date := request.GetDate()
	if date == "" {
		return nil, errBadDate
	}

	topPositions, err := s.cmd.GetTopPositions(
		ctx,
		commander.Request{
			CountryCode: countryCode,
			AppID:       appID,
			Date:        model.Date(date),
		},
	)
	if err != nil {
		return nil, errFailedToGetTop
	}

	rows := []*Row{}
	for cat, pos := range topPositions {
		rows = append(rows, &Row{
			CategoryId: string(cat),
			Position:   int32(pos),
		})
	}

	return &AppTopCategoryResponse{Data: rows}, nil
}

func New(addr string, cmd commander.Commander) *grpc.Server {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	RegisterAppticaServiceServer(s, &srv{cmd: cmd})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return s
}
