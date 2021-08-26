package repository

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/online-store/consts"
	"github.com/online-store/entity"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

type OrderRepository interface {
	Create(ctx context.Context, order entity.Order) entity.Response
	Update(ctx context.Context, order entity.Order) entity.Response
	Get(ctx context.Context, id uint) entity.Response
	GetExpiredCheckout(ctx context.Context) ([]entity.Order, error)
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

// method Create digunakan untuk membuat data order
func (o *orderRepository) Create(ctx context.Context, order entity.Order) entity.Response {
	err := o.db.WithContext(ctx).Create(&order).Error
	if err != nil {
		fmt.Println(err)
		return entity.NewResponse(http.StatusInternalServerError, consts.InternalServerErrorMessage, struct{}{})
	}
	return entity.NewResponse(http.StatusCreated, consts.CreatedMessage, order)
}

// method Update digunakan untuk mengubah data order
func (o *orderRepository) Update(ctx context.Context, order entity.Order) entity.Response {
	err := o.db.Transaction(func(tx *gorm.DB) error {

		// jika status akan diubah ke checkout atau expired maka update stock item
		if order.Status == consts.Checkout || order.Status == consts.Expired {
			var expression string

			// jika status akan diubah ke checkout maka kurangi stock item
			// jika status akan diubah ke expired maka kembalikan stock item ke angka semula
			// dengan cara menambahkan stock sesuai qty order item
			if order.Status == consts.Checkout {
				expression = "stock - ?"
			} else {
				expression = "stock + ?"
			}
			for _, item := range order.OrderItems {
				//query untuk update stock item
				if err := tx.Model(entity.Item{}).Where("id = ?", item.ItemId).UpdateColumn("stock", gorm.Expr(expression, item.ItemQty)).Error; err != nil {
					return err
				}
			}
		}

		// query untuk update data order
		if err := o.db.WithContext(ctx).Model(entity.Order{}).Where("id = ?", order.ID).Updates(order).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Println(err)
		// jika error yang terjadi karena stock item < 0 maka kondisi return error berikut
		// notes: akan lebih tepat jika menggunakan fungsi errors.Is tapi saya belum menemukan error type untuk contrainst violate error
		if strings.Contains(err.Error(), "constraint") {
			return entity.NewResponse(http.StatusBadRequest, "Jumlah order melebihi stock item tersedia", struct{}{})
		}
		//return internal server error
		return entity.NewResponse(http.StatusInternalServerError, consts.InternalServerErrorMessage, struct{}{})
	}
	return entity.NewResponse(http.StatusOK, consts.SuccessMessage, struct{}{})
}

// Get digunakan untuk mendapatkan satu data order
func (o *orderRepository) Get(ctx context.Context, id uint) entity.Response {
	var order entity.Order
	// get order dan join kan dengan order item
	err := o.db.WithContext(ctx).Model(entity.Order{}).Joins("left join order_items on order_items.order_id = orders.id").Preload("OrderItems").Where("orders.id=?", id).First(&order).Error
	if err != nil {
		log.Println(err)
		return entity.NewResponse(http.StatusInternalServerError, consts.InternalServerErrorMessage, struct{}{})
	}
	return entity.NewResponse(http.StatusOK, consts.SuccessMessage, order)
}

// GetExpiredCheckout digunakan untuk mendapatkan list order yang masih checkout dan sudah melewati batas checkout_expired_at
func (o *orderRepository) GetExpiredCheckout(ctx context.Context) ([]entity.Order, error) {
	var orders []entity.Order
	err := o.db.WithContext(ctx).Model(entity.Order{}).Joins("left join order_items on order_items.order_id = orders.id").Preload("OrderItems").Where("status=? and checkout_expired_at <= now()", consts.Checkout).Find(&orders).Error
	if err != nil {
		log.Println(err)
		return []entity.Order{}, err
	}
	return orders, nil
}
