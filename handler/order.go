package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/online-store/entity"
	"github.com/online-store/helper"
	"github.com/online-store/service"
)

type orderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(app *fiber.App, orderS service.OrderService) {
	handler := &orderHandler{
		orderService: orderS,
	}
	orderGroup := app.Group("/order")
	orderGroup.Post("", handler.Create)
	orderGroup.Put("/:id", handler.Update)

}
func (o *orderHandler) Create(c *fiber.Ctx) error {
	var order entity.Order
	if err := c.BodyParser(&order); err != nil {
		log.Println(err)
		return helper.WriteResponse(c, entity.NewResponse(http.StatusBadRequest, err.Error(), struct{}{}))
	}

	response := o.orderService.Create(c.Context(), order)

	return helper.WriteResponse(c, response)
}

func (o *orderHandler) Update(c *fiber.Ctx) error {
	paramsId := c.Params("id")
	if paramsId == "" {
		return helper.WriteResponse(c, entity.NewResponse(http.StatusBadRequest, "id cannot be empty", struct{}{}))
	}

	id, err := strconv.ParseUint(paramsId, 10, 64)

	if err != nil {
		return helper.WriteResponse(c, entity.NewResponse(http.StatusBadRequest, "id is invalid", struct{}{}))
	}
	var order entity.Order
	if err := c.BodyParser(&order); err != nil {
		log.Println(err)
		return helper.WriteResponse(c, entity.NewResponse(http.StatusBadRequest, err.Error(), struct{}{}))
	}
	response := o.orderService.Update(c.Context(), uint(id), order.Status)

	return helper.WriteResponse(c, response)

}
