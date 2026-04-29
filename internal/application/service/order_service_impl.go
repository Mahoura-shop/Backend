package service

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	orderdto "github.com/Mahoura-shop/Backend/internal/application/dto/order"
	"github.com/Mahoura-shop/Backend/internal/application/service/pricing"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type OrderService struct {
	constants         *bootstrap.Constants
	orderRepository   postgres.OrderRepository
	productRepository postgres.ProductRepository
	currencyService   usecase.CurrencyService
	productService    usecase.ProductService
	userService       usecase.UserService
	cartService       usecase.CartService
	db                database.Database
}

type OrderServiceDeps struct {
	Constants         *bootstrap.Constants
	OrderRepository   postgres.OrderRepository
	ProductRepository postgres.ProductRepository
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
		productRepository: deps.ProductRepository,
		currencyService:   deps.CurrencyService,
		productService:    deps.ProductService,
		userService:       deps.UserService,
		cartService:       deps.CartService,
		db:                deps.DB,
	}
}

func (orderService *OrderService) ParseOrderItem(orderItem entity.OrderItem) orderdto.OrderItemCredential {
	response := orderdto.OrderItemCredential{
		ID:            orderItem.ID,
		Count:         orderItem.Count,
		PriceSnapshot: orderItem.PriceSnapshot,
		Tier:          orderItem.Tier,
	}
	response.Product = orderService.productService.ParseProduct(orderItem.Product)
	return response
}

func (orderService *OrderService) ParseOrder(order entity.Order) orderdto.OrderCredential {
	response := orderdto.OrderCredential{
		ID:            order.ID,
		UserID:        order.UserID,
		PaymentMethod: order.PaymentMethod,
	}

	var items []orderdto.OrderItemCredential
	for _, item := range order.Items {
		items = append(items, orderService.ParseOrderItem(item))
	}
	response.Items = items

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

func validatePaymentMethod(userType enum.UserType, paymentMethod enum.PaymentMethod) error {
	switch paymentMethod {
	case enum.PaymentMethodCash:
		if userType != enum.UserTypeShopkeeperCash {
			return exception.ForbiddenError{Message: "cash payment is only available for shopkeeperCash users"}
		}
	case enum.PaymentMethodInstallment:
		if userType != enum.UserTypeShopkeeperCheque {
			return exception.ForbiddenError{Message: "instalment payment is only available for shopkeeperCheque users"}
		}
	}
	return nil
}

func (orderService *OrderService) RegisterOrder(userID uint, paymentMethod enum.PaymentMethod) error {
	user, err := orderService.userService.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return exception.NotFoundError{Item: orderService.constants.Field.User}
	}

	if err := validatePaymentMethod(user.Type, paymentMethod); err != nil {
		return err
	}

	cart, err := orderService.cartService.GetUserCart(userID)
	if err != nil {
		return err
	}

	if len(cart.Items) == 0 {
		return exception.NotFoundError{Item: orderService.constants.Field.CartItem}
	}

	return orderService.db.WithTransaction(func(tx database.Database) error {
		order := entity.Order{
			UserID:        userID,
			PaymentMethod: paymentMethod,
		}
		createdOrder, err := orderService.orderRepository.CreateOrder(tx, order)
		if err != nil {
			return err
		}

		for _, cartItem := range cart.Items {
			product, err := orderService.productRepository.FindProductByID(tx, cartItem.Product.ID)
			if err != nil {
				return err
			}
			if product == nil {
				return exception.NotFoundError{Item: orderService.constants.Field.Product}
			}

			if cartItem.Count < product.MinOrder {
				return exception.ForbiddenError{Message: "order quantity is below minimum order requirement"}
			}
			if product.Quantity < cartItem.Count {
				return exception.ForbiddenError{Message: "insufficient stock for product"}
			}

			resolvedPrice := pricing.ResolvePrice(user.Type, *product)

			orderItem := entity.OrderItem{
				OrderID:       createdOrder.ID,
				ProductID:     product.ID,
				Count:         cartItem.Count,
				PriceSnapshot: resolvedPrice,
				Tier:          user.Type,
			}
			if err := orderService.orderRepository.CreateOrderItem(tx, orderItem); err != nil {
				return err
			}

			product.Quantity -= cartItem.Count
			if err := orderService.productRepository.UpdateProduct(tx, *product); err != nil {
				return err
			}
		}

		return orderService.cartService.DeleteCartItems(cart.ID)
	})
}
