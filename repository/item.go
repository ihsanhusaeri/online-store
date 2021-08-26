package repository

import (
	"context"
	"log"
	"net/http"

	"github.com/online-store/consts"
	"github.com/online-store/entity"
	"gorm.io/gorm"
)

type itemRepository struct {
	db *gorm.DB
}

type ItemRepository interface {
	Get(ctx context.Context, id uint) entity.Response
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{
		db: db,
	}
}

//Get digunakan untuk mendapatkan satu data item sesuai id yang dikirim
func (o *itemRepository) Get(ctx context.Context, id uint) entity.Response {
	var item entity.Item
	err := o.db.WithContext(ctx).Model(entity.Item{}).Where("id=?", id).First(&item).Error

	if err != nil {
		log.Println(err)
		return entity.NewResponse(http.StatusInternalServerError, consts.InternalServerErrorMessage, struct{}{})
	}
	return entity.NewResponse(http.StatusOK, consts.CreatedMessage, item)
}
