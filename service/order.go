package service

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/online-store/consts"
	"github.com/online-store/entity"
	"github.com/online-store/helper"
	"github.com/online-store/repository"
)

type orderService struct {
	orderRepo repository.OrderRepository
	itemRepo  repository.ItemRepository
}

type OrderService interface {
	Create(ctx context.Context, order entity.Order) entity.Response
	Update(ctx context.Context, id uint, status string) entity.Response
}

func NewOrderService(repo repository.OrderRepository, itemR repository.ItemRepository) OrderService {
	return &orderService{
		orderRepo: repo,
		itemRepo:  itemR,
	}
}

func (o *orderService) Create(ctx context.Context, order entity.Order) entity.Response {
	for _, item := range order.OrderItem {
		responseItem := o.itemRepo.Get(ctx, item.ItemId)
		if responseItem.Code != http.StatusOK {
			return responseItem
		}
		itemData := responseItem.Data.(entity.Item)
		if item.ItemQty > itemData.Stock {
			return helper.NewResponse(http.StatusBadRequest, "Jumlah order melebihi stock item tersedia", struct{}{})
		}
	}
	response := o.orderRepo.Create(ctx, order)
	return response
}

func (o *orderService) Update(ctx context.Context, id uint, status string) entity.Response {
	response := o.orderRepo.Get(ctx, id)
	if response.Code != http.StatusOK {
		return response
	}
	log.Println(response)
	order := response.Data.(entity.Order)
	order.Status = status
	if order.Status == consts.Checkout {
		expired := time.Now().Local().Add(time.Minute * time.Duration(5))
		order.CheckoutExpiredAt = &expired
	}
	return o.orderRepo.Update(ctx, order)
}
