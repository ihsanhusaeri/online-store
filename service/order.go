package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/online-store/consts"
	"github.com/online-store/entity"
	"github.com/online-store/repository"
	"github.com/robfig/cron"
)

type orderService struct {
	orderRepo repository.OrderRepository
	itemRepo  repository.ItemRepository
}

type OrderService interface {
	Create(ctx context.Context, order entity.Order) entity.Response
	Update(ctx context.Context, id uint, status string) entity.Response
	CheckExpiredCheckout(intervalMinute int) error
}

func NewOrderService(repo repository.OrderRepository, itemR repository.ItemRepository) OrderService {
	return &orderService{
		orderRepo: repo,
		itemRepo:  itemR,
	}
}

func (o *orderService) Create(ctx context.Context, order entity.Order) entity.Response {
	for _, item := range order.OrderItems {
		responseItem := o.itemRepo.Get(ctx, item.ItemId)
		if responseItem.Code != http.StatusOK {
			return responseItem
		}
		itemData := responseItem.Data.(entity.Item)
		if item.ItemQty > itemData.Stock {
			return entity.NewResponse(http.StatusBadRequest, "Jumlah order melebihi stock item tersedia", struct{}{})
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
	order := response.Data.(entity.Order)
	order.Status = status
	if order.Status == consts.Checkout {
		expired := time.Now().Local().Add(time.Minute * time.Duration(5))
		order.CheckoutExpiredAt = &expired
	}
	return o.orderRepo.Update(ctx, order)
}

func (o *orderService) CheckExpiredCheckout(intervalMinute int) error {
	log.Printf("Running cron on every %d minutes to check expired checkout\n", intervalMinute)
	c := cron.New()
	schedule := fmt.Sprintf("@every 0h%dm0s", intervalMinute)
	if err := c.AddFunc(schedule,
		func() {
			err := o.checkExpiredCheckout(intervalMinute)
			if err != nil {
				log.Println(err)
			}
		}); err != nil {
		log.Println("Error cronjob:", err)
		return err
	}
	c.Start()
	return nil
}

func (o *orderService) checkExpiredCheckout(intervalMinute int) error {
	ctx := context.Background()
	expiredOrders, err := o.orderRepo.GetExpiredCheckout(ctx)

	if err != nil {
		return err
	}

	for _, order := range expiredOrders {
		order.Status = consts.Expired
		response := o.orderRepo.Update(ctx, order)

		if response.Code != http.StatusOK {
			return errors.New(response.Message)
		}
		log.Printf("Order [%d] updated to %s", order.ID, consts.Expired)
	}
	return nil
}
