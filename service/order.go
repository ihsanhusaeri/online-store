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

//Create digunakan untuk membuat data order
func (o *orderService) Create(ctx context.Context, order entity.Order) entity.Response {
	//mapping data order items
	for _, item := range order.OrderItems {
		//dapatkan data item
		responseItem := o.itemRepo.Get(ctx, item.ItemId)
		// jika terjadi error (response code != 200) maka return response
		if responseItem.Code != http.StatusOK {
			return responseItem
		}
		//cast data interface{} ke Item
		itemData := responseItem.Data.(entity.Item)

		//cek apakah qty order melebihi stock item
		if item.ItemQty > itemData.Stock {
			return entity.NewResponse(http.StatusBadRequest, "Jumlah order melebihi stock item tersedia", struct{}{})
		}
	}
	response := o.orderRepo.Create(ctx, order)
	return response
}

// Update digunakan untuk mengubah data order
func (o *orderService) Update(ctx context.Context, id uint, status string) entity.Response {
	//get existing order
	response := o.orderRepo.Get(ctx, id)

	// jika terjadi error (response code != 200) maka return response
	if response.Code != http.StatusOK {
		return response
	}
	//cast data interface{} ke Order
	order := response.Data.(entity.Order)
	order.Status = status

	//jika status order diubah ke checkout maka set checkout_expired_at +5 menit setelah checkout
	if order.Status == consts.Checkout {
		expired := time.Now().Local().Add(time.Minute * time.Duration(5))
		order.CheckoutExpiredAt = &expired
	}
	return o.orderRepo.Update(ctx, order)
}

//CheckExpiredCheckout adalah cronjob function yang dijalankan setiap jangka waktu tertentu (sesuai parameter intervalMinute)
// function ini digunakan untuk mengecek order-order yang masih checkout dan sudah melewati batas checkout_expired_at
// order-order tersebut akan diubah statusnya menjadi expired dan stock item akan dikembalikan ke angka semula
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

	//get semua expired order
	expiredOrders, err := o.orderRepo.GetExpiredCheckout(ctx)

	if err != nil {
		return err
	}

	for _, order := range expiredOrders {
		//set status order ke expired
		order.Status = consts.Expired
		response := o.orderRepo.Update(ctx, order)

		if response.Code != http.StatusOK {
			return errors.New(response.Message)
		}
		log.Printf("Order [%d] updated to %s", order.ID, consts.Expired)
	}
	return nil
}
