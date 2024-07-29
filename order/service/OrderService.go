package service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"order/dto"
	"order/model"
	"order/proto/compiled"
	rpc "order/service/grpc"
	"order/util"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetOrder(c *fiber.Ctx, db *gorm.DB, client *compiled.ProductServiceClient) error {
	code := c.Params("code")
	var orders []model.Order
	var productIds []int64
	orderMap := make(map[int64][2]int64)

	err := db.Find(&orders, "order_code = ?", code).Error
	if err != nil {
		return err
	}

	for _, order := range orders {
		productIds = append(productIds, order.ProductId)
		orderMap[order.ProductId] = [2]int64{int64(order.Quantity), order.FinalPrice}
	}

	products, err := rpc.GetProducts(*client, productIds)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var productDtos []dto.ProductDetail
	var totalPrice int64

	for _, product := range products.Products {
		productDtos = append(productDtos, dto.ProductDetail{
			ID:       product.Id,
			Name:     product.Name,
			Quantity: int32(orderMap[product.Id][0]),
			Price:    orderMap[product.Id][1],
		})
		totalPrice += orderMap[product.Id][1]
	}

	return c.JSON(dto.OrderResponse{
		OrderCode:  code,
		Products:   productDtos,
		TotalPrice: totalPrice,
	})
}

func AddOrder(c *fiber.Ctx, db *gorm.DB, client *compiled.ProductServiceClient) error {
	body := c.Body()
	request := new(dto.OrderRequest)
	err := json.Unmarshal(body, request)
	if err != nil {
		return err
	}

	var orders []model.Order
	var productIds []int64
	millis := time.Now().UnixNano() / int64(time.Millisecond)
	randomString := util.RandStringBytes(7, rand.New(rand.NewSource(time.Now().UnixNano())))

	for _, product := range request.Products {
		orders = append(orders, model.Order{
			OrderCode:  fmt.Sprintf("ORD-%d-%s", millis, randomString),
			CustomerId: request.CustomerId,
			ProductId:  product.ProductId,
			Quantity:   product.Quantity,
		})
		productIds = append(productIds, product.ProductId)
	}

	products, err := rpc.GetProducts(*client, productIds)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	for i := 0; i < len(products.Products); i++ {
		orders[i].FinalPrice = products.Products[i].Price * int64(orders[i].Quantity)
	}

	res := db.Create(&orders)
	if res.Error != nil {
		return res.Error
	}

	return c.SendString("Success")
}

func UpdateOrder(c *fiber.Ctx, db *gorm.DB, client *compiled.ProductServiceClient) error {
	code := c.Params("code")
	var orders []*model.Order
	orderMap := make(map[int64]int32)
	body := c.Body()
	request := new(dto.OrderRequest)
	err := json.Unmarshal(body, request)
	if err != nil {
		return err
	}
	var productIds []int64

	for _, product := range request.Products {
		orderMap[product.ProductId] = product.Quantity
		productIds = append(productIds, product.ProductId)
	}

	err = db.Find(&orders, "order_code = ?", code).Error
	if err != nil {
		return err
	}

	products, err := rpc.GetProducts(*client, productIds)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	for i, order := range orders {
		order.Quantity = orderMap[order.ProductId]
		order.FinalPrice = products.Products[i].Price * int64(orderMap[order.ProductId])
	}

	res := db.Save(&orders)
	if res.Error != nil {
		return res.Error
	}

	return c.SendString("Success")
}

func DeleteOrder(c *fiber.Ctx, db *gorm.DB) error {
	code := c.Params("code")

	err := db.Where("order_code = ?", code).Delete(&model.Order{}).Error
	if err != nil {
		return err
	}

	return c.SendString("Success")
}
