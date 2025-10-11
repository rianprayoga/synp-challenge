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

func (it *ItemsGrpcHandler) Reservestock(c context.Context, r *pb.ReserveStockRequest) (*pb.ReserverResponse, error) {
	item, err := it.DB.ReduceStock(r.ItemId, int(r.Qty), int(r.Version))
	if err != nil {
		return nil, err
	}

	return &pb.ReserverResponse{
		ItemId:         item.ID,
		Version:        int32(item.Version),
		StockReduced:   r.Qty,
		StockRemaining: int32(item.Stock),
	}, nil
}

func (it *ItemsGrpcHandler) CheckStock(c context.Context, i *pb.CheckStockRequest) (*pb.ItemStock, error) {
	item, err := it.DB.GetItem(i.ItemId)
	if err != nil {
		return nil, appError.ErrItemNotFound
	}

	return &pb.ItemStock{
		ItemId:  item.ID,
		Stock:   int32(item.Stock),
		Version: int32(item.Version),
	}, nil
}
