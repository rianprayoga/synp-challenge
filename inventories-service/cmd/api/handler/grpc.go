package handler

import (
	"context"
	appError "inventories-app/internal/error"
	"inventories-app/internal/repository"
	pb "rpc"

	"google.golang.org/grpc"
)

type ItemsGrpcHandler struct {
	DB repository.DBRepo
	pb.UnimplementedInventoryRpcServer
}

func NewGrpcItemService(grpc *grpc.Server, repo repository.DBRepo) {

	pb.RegisterInventoryRpcServer(grpc, &ItemsGrpcHandler{
		DB: repo,
	})
}

func (it *ItemsGrpcHandler) CheckStock(c context.Context, i *pb.Item) (*pb.ItemStock, error) {
	item, err := it.DB.GetItem(i.ItemId)
	if err != nil {
		return nil, appError.ErrItemNotFound
	}

	return &pb.ItemStock{
		ItemId: item.ID,
		Stock:  int32(item.Stock),
	}, nil
}
