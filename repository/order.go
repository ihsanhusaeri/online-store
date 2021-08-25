package repository

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/online-store/consts"
	"github.com/online-store/entity"
	"github.com/online-store/helper"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

type OrderRepository interface {
	Create(ctx context.Context, order entity.Order) entity.Response
	Update(ctx context.Context, order entity.Order) entity.Response
	Get(ctx context.Context, id uint) entity.Response
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

var response entity.Response

func (o *orderRepository) Create(ctx context.Context, order entity.Order) entity.Response {
	err := o.db.WithContext(ctx).Create(&order).Error
	if err != nil {
		fmt.Println(err)
		return helper.NewResponse(http.StatusInternalServerError, consts.InternalServerErrorMessage, struct{}{})
	}
	return helper.NewResponse(http.StatusCreated, consts.CreatedMessage, order)
}
func (o *orderRepository) Update(ctx context.Context, order entity.Order) entity.Response {
	err := o.db.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		for _, item := range order.OrderItem {
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
		return helper.NewResponse(http.StatusInternalServerError, consts.InternalServerErrorMessage, struct{}{})
	}
	return helper.NewResponse(http.StatusOK, consts.SuccessMessage, struct{}{})
}

func (o *orderRepository) Get(ctx context.Context, id uint) entity.Response {
	var order entity.Order
	err := o.db.WithContext(ctx).Model(entity.Order{}).Joins("left join order_items on order_items.order_id = orders.id").Preload("OrderItem").Where("orders.id=?", id).First(&order).Error
	log.Println(order)
	if err != nil {
		log.Println(err)
		return helper.NewResponse(http.StatusInternalServerError, consts.InternalServerErrorMessage, struct{}{})
	}
	return helper.NewResponse(http.StatusOK, consts.SuccessMessage, order)
}
