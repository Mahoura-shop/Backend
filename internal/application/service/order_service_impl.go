package service

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	orderdto "github.com/Mahoura-shop/Backend/internal/application/dto/order"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type OrderService struct {
	constants         *bootstrap.Constants
	orderRepository   postgres.OrderRepository
	currencyService   usecase.CurrencyService
	productService    usecase.ProductService
	userService       usecase.UserService
	cartService       usecase.CartService
	db                database.Database
}

type OrderServiceDeps struct {
	Constants         *bootstrap.Constants
	OrderRepository   postgres.OrderRepository
	CurrencyService   usecase.CurrencyService
	ProductService    usecase.ProductService
	UserService       usecase.UserService
	CartService       usecase.CartService
	DB                database.Database
}

func NewOrderService(deps OrderServiceDeps) *OrderService {
	return &OrderService{
		constants:         deps.Constants,
		orderRepository:   deps.OrderRepository,
		currencyService:   deps.CurrencyService,
		productService:    deps.ProductService,
		userService:       deps.UserService,
		cartService:       deps.CartService,
		db:                deps.DB,
	}
}

func (orderService *OrderService) ParseOrderItem(orderItem entity.OrderItem) (orderdto.OrderItemCredential) {
	response := orderdto.OrderItemCredential{
		ID:    orderItem.ID,
		Count: orderItem.Count,
	}
	
	product := orderService.productService.ParseProduct(orderItem.Product)
	response.Product = product
	// response.PriceSnapshot = product

	return response
}

func (orderService *OrderService) ParseOrder(order entity.Order) (orderdto.OrderCredential) {
	response := orderdto.OrderCredential{
		ID: order.ID,
	}
	
	var items []orderdto.OrderItemCredential
	for _, item := range order.Items {
		response := orderService.ParseOrderItem(item)
		items = append(items, response)
	}
	
	return response
}

func (orderService *OrderService) GetOrder(orderID uint) (*orderdto.OrderCredential, error) {
	order, err := orderService.orderRepository.FindOrderByID(orderService.db, orderID)
	if err != nil {
		return nil, err
	}
	
	if order == nil {
		return nil, exception.NotFoundError{Item: orderService.constants.Field.Order}
	}
	
	parsedOrder := orderService.ParseOrder(*order)
	return &parsedOrder, nil
}

func (orderService *OrderService) GetOrders() ([]orderdto.OrderCredential, error) {
	orders, err := orderService.orderRepository.GetOrders(orderService.db)
	if err != nil {
		return nil, err
	}

	var responses []orderdto.OrderCredential
	for _, order := range orders {
		response := orderService.ParseOrder(*order)
		responses = append(responses, response)
	}

	return responses, nil
}

func (orderService *OrderService) RegisterOrder(userID uint) (error) {
	cart, err := orderService.cartService.GetUserCart(userID)
	if err != nil {
		return err
	}

	return orderService.db.WithTransaction(func(tx database.Database) error {
		order, err := orderService.orderRepository.CreateOrder(tx); 
		if err != nil {
			return err
		}

		for _, cartItem := range cart.Items {
			orderItem := entity.OrderItem{
				OrderID:   order.ID,
				ProductID: cartItem.Product.ID,
				Count:     cartItem.Count,
				// PriceSnapshot: getProductPriceForRole(),
			}
			err := orderService.orderRepository.CreateOrderItem(tx, orderItem)
			if err != nil {
				return err
			}
		}

		if err = orderService.cartService.DeleteCartItems(cart.ID); err != nil {
			return err
		}

		return nil
	})
}