package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/online-store/consts"
	"github.com/online-store/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

func (mock *MockOrderRepository) Create(ctx context.Context, order entity.Order) entity.Response {
	args := mock.Called(ctx, order)
	result := args.Get(0)
	return result.(entity.Response)
}

func (mock *MockOrderRepository) Update(ctx context.Context, order entity.Order) entity.Response {
	args := mock.Called(ctx, order)
	result := args.Get(0)
	return result.(entity.Response)
}

func (mock *MockOrderRepository) Get(ctx context.Context, id uint) entity.Response {
	args := mock.Called(ctx, id)
	result := args.Get(0)
	return result.(entity.Response)
}

func (mock *MockOrderRepository) GetExpiredCheckout(ctx context.Context) ([]entity.Order, error) {
	args := mock.Called(ctx)
	result := args.Get(0)
	return result.([]entity.Order), args.Error(1)
}

type MockItemRepository struct {
	mock.Mock
}

func (mock *MockItemRepository) Get(ctx context.Context, id uint) entity.Response {
	args := mock.Called(ctx, id)
	result := args.Get(0)
	return result.(entity.Response)
}

func TestCreate(t *testing.T) {
	mockOrderRepo := new(MockOrderRepository)
	mockItemRepo := new(MockItemRepository)

	order := entity.Order{
		UserId:   1,
		UserName: "Ahmad",
		TotalPay: 30000,
		OrderItems: []entity.OrderItem{
			{
				ItemId:    1,
				ItemName:  "buku",
				ItemPrice: 15000,
				ItemQty:   2,
			},
		},
	}
	response := entity.Response{Code: 201, Message: consts.CreatedMessage, Data: order}
	mockOrderRepo.On("Create").Return(response)

	testService := NewOrderService(mockOrderRepo, mockItemRepo)

	result := testService.Create(context.Background(), order)

	mockOrderRepo.AssertExpectations(t)

	assert.Equal(t, http.StatusCreated, result.Code)
	assert.Equal(t, consts.CreatedMessage, result.Message)
	assert.Equal(t, order, result.Data.(entity.Order))
}
