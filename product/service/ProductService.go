package service

import (
	"context"
	"product/model"
	"product/proto/compiled"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type ProductService struct {
	compiled.UnimplementedProductServiceServer
	DB *gorm.DB
}

func (ps ProductService) GetProducts(_ context.Context, payload *compiled.ProductIds) (*compiled.Products, error) {
	var products []model.Product

	result := ps.DB.Find(&products, "id in ?", payload.Ids)
	if result.Error != nil {
		return nil, result.Error
	}

	return modelToDTO(products), nil
}

func modelToDTO(models []model.Product) *compiled.Products {
	var products = make([]*compiled.Product, 0, len(models))
	for _, model := range models {
		products = append(products, &compiled.Product{
			Id:        model.ID,
			Name:      model.Name,
			Stock:     model.Stock,
			Price:     model.Price,
			CreatedAt: timestamppb.New(model.CreatedAt),
			UpdatedAt: timestamppb.New(model.UpdatedAt),
		})
	}

	return &compiled.Products{Products: products}
}
