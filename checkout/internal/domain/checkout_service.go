package domain

import (
	"context"
	"route256/checkout/internal/models"

	"github.com/pkg/errors"
)

type StocksChecker interface {
	GetStocks(ctx context.Context, sku uint32) ([]models.Stock, error)
}

type ProductsChecker interface {
	GetProducts(ctx context.Context, skus []uint32) ([]models.Product, error)
}

type CreateOrderChecker interface {
	CreateOrder(ctx context.Context, user int64, items []models.CreateOrderItem) error
}

type Service struct {
	stocksChecker      StocksChecker
	productsChecker    ProductsChecker
	createOrderChecker CreateOrderChecker
}

var (
	ErrInsufficientStocks = errors.New("insufficient stocks")
)

func New(
	stocksChecker StocksChecker,
	productsChecker ProductsChecker,
	createOrderChecker CreateOrderChecker,
) *Service {
	return &Service{
		stocksChecker:      stocksChecker,
		productsChecker:    productsChecker,
		createOrderChecker: createOrderChecker,
	}
}

func (s *Service) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	stocks, err := s.stocksChecker.GetStocks(ctx, sku)
	if err != nil {
		return errors.WithMessage(err, "checking stocks")
	}

	counter := int64(count)
	for _, stock := range stocks {
		counter -= int64(stock.Count)
		if counter <= 0 {
			return nil
		}
	}

	return ErrInsufficientStocks
}

func (s *Service) DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	return nil
}

func (s *Service) ListCart(ctx context.Context, user int64) (*models.CartInfo, error) {
	skus := []uint32{
		1076963,
		1148162,
		1625903,
		2618151,
		2956315,
		2958025,
		3596599,
		3618852,
		4288068,
		4465995,
	}

	products, err := s.productsChecker.GetProducts(ctx, skus)
	if err != nil {
		return nil, errors.WithMessage(err, "get products")
	}
	items := make([]models.Item, 0, len(products))
	cartInfo := models.CartInfo{Items: items}
	for _, product := range products {
		cartInfo.TotalPrice += product.Price
		cartInfo.Items = append(
			cartInfo.Items,
			models.Item{
				SKU:   1,
				Count: 2,
				Name:  product.Name,
				Price: product.Price,
			},
		)
	}
	return &cartInfo, nil
}

func (s *Service) Purchase(ctx context.Context, user int64) error {
	items := []models.CreateOrderItem{
		{Sku: 1, Count: 10},
		{Sku: 2, Count: 20},
	}
	err := s.createOrderChecker.CreateOrder(ctx, user, items)
	if err != nil {
		return errors.WithMessage(err, "purchase")
	}
	return nil
}
