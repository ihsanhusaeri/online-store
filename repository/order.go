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

func (o *orderRepository) Create(ctx context.Context, order entity.Order) entity.Response {
	err := o.db.WithContext(ctx).Create(&order).Error
	if err != nil {
		fmt.Println(err)
		return entity.NewResponse(http.StatusInternalServerError, consts.InternalServerErrorMessage, struct{}{})
	}
	return entity.NewResponse(http.StatusCreated, consts.CreatedMessage, order)
}
func (o *orderRepository) Update(ctx context.Context, order entity.Order) entity.Response {
	err := o.db.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		for _, item := range order.OrderItems {
			if err := tx.Model(entity.Item{}).Where("id = ?", item.ItemId).UpdateColumn("stock", gorm.Expr("stock - ?", item.ItemQty)).Error; err != nil {
				return err
			}
		}

		if err := o.db.WithContext(ctx).Model(entity.Order{}).Where("id = ?", order.ID).Updates(order).Error; err != nil {
			return err
		}

		// return nil will commit the whole transaction
		return nil
	})
	if err != nil {
		log.Println(err)
		if strings.Contains(err.Error(), "constraint") {
			return entity.NewResponse(http.StatusBadRequest, "Jumlah order melebihi stock item tersedia", struct{}{})
		}
		return entity.NewResponse(http.StatusInternalServerError, consts.InternalServerErrorMessage, struct{}{})
	}
	return entity.NewResponse(http.StatusOK, consts.SuccessMessage, struct{}{})
}

func (o *orderRepository) Get(ctx context.Context, id uint) entity.Response {
	var order entity.Order
	err := o.db.WithContext(ctx).Model(entity.Order{}).Joins("left join order_items on order_items.order_id = orders.id").Preload("OrderItems").Where("orders.id=?", id).First(&order).Error
	if err != nil {
		log.Println(err)
		return entity.NewResponse(http.StatusInternalServerError, consts.InternalServerErrorMessage, struct{}{})
	}
	return entity.NewResponse(http.StatusOK, consts.SuccessMessage, order)
}

func (o *orderRepository) GetExpiredCheckout(ctx context.Context) ([]entity.Order, error) {
	var orders []entity.Order
	err := o.db.WithContext(ctx).Model(entity.Order{}).Where("status=? and checkout_expired_at <= now()", consts.Checkout).Find(&orders).Error
	if err != nil {
		log.Println(err)
		return []entity.Order{}, err
	}
	return orders, nil
}
