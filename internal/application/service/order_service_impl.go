package service

import (
	"fmt"
	"time"

	"github.com/Mahoura-shop/Backend/bootstrap"
	orderdto "github.com/Mahoura-shop/Backend/internal/application/dto/order"
	"github.com/Mahoura-shop/Backend/internal/application/service/pricing"
	"github.com/Mahoura-shop/Backend/internal/application/service/shipping"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/domain/communication"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type OrderService struct {
	constants             *bootstrap.Constants
	orderRepository       postgres.OrderRepository
	productRepository     postgres.ProductRepository
	walletRepository      postgres.WalletRepository
	transactionRepository postgres.TransactionRepository
	paymentRepository     postgres.PaymentRepository
	instalmentRepository  postgres.InstalmentRepository
	addressRepository     postgres.AddressRepository
	currencyService       usecase.CurrencyService
	productService        usecase.ProductService
	userService           usecase.UserService
	cartService           usecase.CartService
	paymentService        usecase.PaymentService
	notificationService   usecase.NotificationService
	smsService            communication.SMSService
	emailService          communication.EmailService
	db                    database.Database
}

type OrderServiceDeps struct {
	Constants             *bootstrap.Constants
	OrderRepository       postgres.OrderRepository
	ProductRepository     postgres.ProductRepository
	WalletRepository      postgres.WalletRepository
	TransactionRepository postgres.TransactionRepository
	PaymentRepository     postgres.PaymentRepository
	InstalmentRepository  postgres.InstalmentRepository
	AddressRepository     postgres.AddressRepository
	CurrencyService       usecase.CurrencyService
	ProductService        usecase.ProductService
	UserService           usecase.UserService
	CartService           usecase.CartService
	PaymentService        usecase.PaymentService
	NotificationService   usecase.NotificationService
	SMSService            communication.SMSService
	EmailService          communication.EmailService
	DB                    database.Database
}

func NewOrderService(deps OrderServiceDeps) *OrderService {
	return &OrderService{
		constants:             deps.Constants,
		orderRepository:       deps.OrderRepository,
		productRepository:     deps.ProductRepository,
		walletRepository:      deps.WalletRepository,
		transactionRepository: deps.TransactionRepository,
		paymentRepository:     deps.PaymentRepository,
		instalmentRepository:  deps.InstalmentRepository,
		addressRepository:     deps.AddressRepository,
		currencyService:       deps.CurrencyService,
		productService:        deps.ProductService,
		userService:           deps.UserService,
		cartService:           deps.CartService,
		paymentService:        deps.PaymentService,
		notificationService:   deps.NotificationService,
		smsService:            deps.SMSService,
		emailService:          deps.EmailService,
		db:                    deps.DB,
	}
}

func (s *OrderService) ParseOrder(order entity.Order) orderdto.OrderCredential {
	resp := orderdto.OrderCredential{
		ID:            order.ID,
		UserID:        order.UserID,
		AddressID:     order.AddressID,
		PaymentMethod: order.PaymentMethod,
		Status:        order.Status,
		TotalAmount:   order.TotalAmount,
		ShippingCost:  order.ShippingCost,
		RefundFlag:    order.RefundFlag,
		TrackingCode:  order.TrackingCode,
		CreatedAt:     order.CreatedAt,
	}

	for _, item := range order.Items {
		resp.Items = append(resp.Items, orderdto.OrderItemCredential{
			ID:            item.ID,
			Product:       s.productService.ParseProduct(item.Product),
			Count:         item.Count,
			PriceSnapshot: item.PriceSnapshot,
			Tier:          item.Tier,
		})
	}

	for _, h := range order.StatusHistory {
		resp.StatusHistory = append(resp.StatusHistory, orderdto.OrderStatusHistoryEntry{
			Status:      h.Status,
			Note:        h.Note,
			ChangedByID: h.ChangedByID,
			CreatedAt:   h.CreatedAt,
		})
	}

	return resp
}

func (s *OrderService) GetOrder(orderID uint) (*orderdto.OrderCredential, error) {
	order, err := s.orderRepository.FindOrderByID(s.db, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, exception.NotFoundError{Item: s.constants.Field.Order}
	}
	parsed := s.ParseOrder(*order)
	return &parsed, nil
}

func (s *OrderService) GetOrders() ([]orderdto.OrderCredential, error) {
	orders, err := s.orderRepository.GetOrders(s.db)
	if err != nil {
		return nil, err
	}
	var result []orderdto.OrderCredential
	for _, o := range orders {
		result = append(result, s.ParseOrder(*o))
	}
	return result, nil
}

func (s *OrderService) GetUserOrders(userID uint) ([]orderdto.OrderCredential, error) {
	orders, err := s.orderRepository.GetOrdersByUserID(s.db, userID)
	if err != nil {
		return nil, err
	}
	var result []orderdto.OrderCredential
	for _, o := range orders {
		result = append(result, s.ParseOrder(*o))
	}
	return result, nil
}

func (s *OrderService) GetOrdersByStatus(status enum.OrderStatus) ([]orderdto.OrderCredential, error) {
	orders, err := s.orderRepository.GetOrders(s.db)
	if err != nil {
		return nil, err
	}
	var result []orderdto.OrderCredential
	for _, o := range orders {
		if o.Status == status {
			result = append(result, s.ParseOrder(*o))
		}
	}
	return result, nil
}

func validatePaymentMethod(userType enum.UserType, paymentMethod enum.PaymentMethod) error {
	switch paymentMethod {
	case enum.PaymentMethodCash:
		if userType != enum.UserTypeShopkeeperCash && userType != enum.UserTypeShopkeeperCheque {
			return exception.ForbiddenError{Message: "cash payment is only available for shopkeeper users"}
		}
	case enum.PaymentMethodInstallment:
		return exception.ForbiddenError{Message: "instalment payment is no longer available online — please contact sales"}
	}
	return nil
}

func (s *OrderService) RegisterOrder(userID uint, req orderdto.CreateOrderRequest) (uint, error) {
	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		return 0, err
	}
	if user == nil {
		return 0, exception.NotFoundError{Item: s.constants.Field.User}
	}

	if err := validatePaymentMethod(user.Type, req.PaymentMethod); err != nil {
		return 0, err
	}

	cart, err := s.cartService.GetUserCart(userID)
	if err != nil {
		return 0, err
	}
	if len(cart.Items) == 0 {
		return 0, exception.NotFoundError{Item: s.constants.Field.CartItem}
	}

	var shippingCost uint
	if req.AddressID != nil {
		addr, err := s.addressRepository.GetAddressByID(s.db, *req.AddressID)
		if err != nil {
			return 0, err
		}
		if addr != nil {
			shippingCost = shipping.CalculateShipping(&addr.Province)
		}
	}

	var createdOrderID uint
	err = s.db.WithTransaction(func(tx database.Database) error {
		var totalAmount uint

		order := entity.Order{
			UserID:        userID,
			AddressID:     req.AddressID,
			PaymentMethod: req.PaymentMethod,
			Status:        enum.OrderStatusPending,
			ShippingCost:  shippingCost,
		}
		createdOrder, err := s.orderRepository.CreateOrder(tx, order)
		if err != nil {
			return err
		}
		createdOrderID = createdOrder.ID

		for _, cartItem := range cart.Items {
			product, err := s.productRepository.FindProductByID(tx, cartItem.Product.ID)
			if err != nil {
				return err
			}
			if product == nil {
				return exception.NotFoundError{Item: s.constants.Field.Product}
			}
			if cartItem.Count < product.MinOrder {
				return exception.ForbiddenError{Message: "order quantity is below minimum order requirement"}
			}
			if product.Quantity < cartItem.Count {
				return exception.ForbiddenError{Message: "insufficient stock for product"}
			}

			resolvedPrice := pricing.ResolvePrice(user.Type, *product)
			lineTotal := resolvedPrice * cartItem.Count
			totalAmount += lineTotal

			orderItem := entity.OrderItem{
				OrderID:       createdOrder.ID,
				ProductID:     product.ID,
				Count:         cartItem.Count,
				PriceSnapshot: resolvedPrice,
				Tier:          user.Type,
			}
			if err := s.orderRepository.CreateOrderItem(tx, orderItem); err != nil {
				return err
			}

			product.Quantity -= cartItem.Count
			if err := s.productRepository.UpdateProduct(tx, *product); err != nil {
				return err
			}
		}

		createdOrder.TotalAmount = totalAmount + shippingCost
		if err := s.orderRepository.UpdateOrder(tx, *createdOrder); err != nil {
			return err
		}

		if err := s.orderRepository.CreateOrderStatusHistory(tx, entity.OrderStatusHistory{
			OrderID:     createdOrder.ID,
			Status:      enum.OrderStatusPending,
			Note:        "سفارش ثبت شد",
			ChangedByID: userID,
		}); err != nil {
			return err
		}

		if req.PaymentMethod == enum.PaymentMethodInstallment && req.InstalmentCount > 0 {
			count := req.InstalmentCount
			intervalDays := req.InstalmentIntervalDays
			if intervalDays == 0 {
				intervalDays = 30
			}
			instalmentAmount := createdOrder.TotalAmount / count
			for i := uint(1); i <= count; i++ {
				instalment := entity.Instalment{
					OrderID: createdOrder.ID,
					Number:  i,
					Amount:  instalmentAmount,
					DueDate: time.Now().AddDate(0, 0, int(i*intervalDays)),
					Status:  enum.InstalmentStatusPending,
				}
				if err := s.instalmentRepository.CreateInstalment(tx, instalment); err != nil {
					return err
				}
			}
		}

		if err := s.cartService.DeleteCartItems(cart.ID); err != nil {
			return err
		}

		return nil
	})

	if err == nil {
		order, _ := s.orderRepository.FindOrderByID(s.db, createdOrderID)
		if order != nil {
			go s.sendOrderConfirmationEmail(order)
		}
	}

	return createdOrderID, err
}

func (s *OrderService) UpdateOrderStatus(orderID uint, req orderdto.UpdateOrderStatusRequest) error {
	order, err := s.orderRepository.FindOrderByID(s.db, orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return exception.NotFoundError{Item: s.constants.Field.Order}
	}

	if order.Status == enum.OrderStatusCancelled {
		return exception.ForbiddenError{Message: "cannot update a cancelled order"}
	}

	return s.db.WithTransaction(func(tx database.Database) error {
		order.Status = req.Status
		if req.Status == enum.OrderStatusShipped && req.TrackingCode != "" {
			order.TrackingCode = &req.TrackingCode
		}
		if err := s.orderRepository.UpdateOrder(tx, *order); err != nil {
			return err
		}

		if err := s.orderRepository.CreateOrderStatusHistory(tx, entity.OrderStatusHistory{
			OrderID:     order.ID,
			Status:      req.Status,
			Note:        req.Note,
			ChangedByID: req.ChangedByID,
		}); err != nil {
			return err
		}

		go s.notifyOrderStatus(order.User.Phone, order.ID, req.Status)
		go s.createOrderNotification(order.UserID, order.ID, req.Status)
		return nil
	})
}

func (s *OrderService) createOrderNotification(userID, orderID uint, status enum.OrderStatus) {
	if s.notificationService == nil {
		return
	}
	title := "به‌روزرسانی وضعیت سفارش"
	body := fmt.Sprintf("وضعیت سفارش شما به «%s» تغییر یافت", status.String())
	ref := orderID
	_ = s.notificationService.CreateNotification(userID, 1, title, body, &ref)
}

func (s *OrderService) notifyOrderStatus(phone string, orderID uint, status enum.OrderStatus) {
	if phone == "" {
		return
	}
	message := fmt.Sprintf("سفارش شماره %d: %s", orderID, status.String())
	_ = s.smsService.SendMessage(phone, message)
}

func (s *OrderService) PayOrderByWallet(userID, orderID uint) error {
	order, err := s.orderRepository.FindOrderByID(s.db, orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return exception.NotFoundError{Item: s.constants.Field.Order}
	}
	if order.UserID != userID {
		return exception.ForbiddenError{Message: "this order does not belong to you"}
	}
	if order.Status != enum.OrderStatusPending {
		return exception.ForbiddenError{Message: "order is not in pending status"}
	}

	return s.db.WithTransaction(func(tx database.Database) error {
		wallet, err := s.walletRepository.FindWalletByUserID(tx, userID)
		if err != nil {
			return err
		}
		if wallet == nil {
			return exception.NotFoundError{Item: "wallet"}
		}
		if wallet.Balance < order.TotalAmount {
			return exception.ForbiddenError{Message: "insufficient wallet balance"}
		}

		if _, err := s.walletRepository.WithdrawWallet(tx, userID, order.TotalAmount); err != nil {
			return err
		}

		if err := s.transactionRepository.CreateTransaction(tx, entity.Transaction{
			Amount:   order.TotalAmount,
			WalletID: wallet.ID,
			Type:     enum.TransactionTypeWithDraw,
		}); err != nil {
			return err
		}

		order.Status = enum.OrderStatusPaid
		if err := s.orderRepository.UpdateOrder(tx, *order); err != nil {
			return err
		}

		if err := s.orderRepository.CreateOrderStatusHistory(tx, entity.OrderStatusHistory{
			OrderID:     order.ID,
			Status:      enum.OrderStatusPaid,
			Note:        "پرداخت از کیف پول",
			ChangedByID: userID,
		}); err != nil {
			return err
		}

		go s.notifyOrderStatus(order.User.Phone, order.ID, enum.OrderStatusPaid)
		go s.createOrderNotification(order.UserID, order.ID, enum.OrderStatusPaid)
		return nil
	})
}

func (s *OrderService) InitiateGatewayPayment(userID, orderID uint) (*orderdto.PaymentGatewayResponse, error) {
	// MOCKED: skip gateway, mark order as paid immediately
	order, err := s.orderRepository.FindOrderByID(s.db, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, exception.NotFoundError{Item: s.constants.Field.Order}
	}
	order.Status = enum.OrderStatusPaid
	if err := s.orderRepository.UpdateOrder(s.db, *order); err != nil {
		return nil, err
	}
	return &orderdto.PaymentGatewayResponse{GatewayURL: ""}, nil

	// order, err := s.orderRepository.FindOrderByID(s.db, orderID)
	// if err != nil {
	// 	return nil, err
	// }
	// if order == nil {
	// 	return nil, exception.NotFoundError{Item: s.constants.Field.Order}
	// }
	// if order.UserID != userID {
	// 	return nil, exception.ForbiddenError{Message: "this order does not belong to you"}
	// }
	// if order.Status != enum.OrderStatusPending {
	// 	return nil, exception.ForbiddenError{Message: "order is not in pending status"}
	// }
	//
	// authority, gatewayURL, err := s.paymentService.InitiateGatewayPayment(orderID, order.TotalAmount)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// payment := entity.Payment{
	// 	OrderID:    orderID,
	// 	Amount:     order.TotalAmount,
	// 	Status:     enum.PaymentStatusPaying,
	// 	Authority:  authority,
	// 	GatewayURL: gatewayURL,
	// }
	// if _, err := s.paymentRepository.CreatePayment(s.db, payment); err != nil {
	// 	return nil, err
	// }
	//
	// return &orderdto.PaymentGatewayResponse{GatewayURL: gatewayURL}, nil
}

func (s *OrderService) VerifyGatewayPayment(authority string, status string) error {
	payment, err := s.paymentRepository.FindPaymentByAuthority(s.db, authority)
	if err != nil {
		return err
	}
	if payment == nil {
		return exception.NotFoundError{Item: s.constants.Field.Payment}
	}

	// Status "OK" means user completed payment on the gateway
	if status != "OK" {
		payment.Status = enum.PaymentStatusError
		_ = s.paymentRepository.UpdatePayment(s.db, *payment)
		return exception.ForbiddenError{Message: "payment was not completed"}
	}

	order, err := s.orderRepository.FindOrderByID(s.db, payment.OrderID)
	if err != nil {
		return err
	}
	if order == nil {
		return exception.NotFoundError{Item: s.constants.Field.Order}
	}

	refCode, err := s.paymentService.VerifyGatewayPayment(authority, payment.Amount)
	if err != nil {
		payment.Status = enum.PaymentStatusError
		_ = s.paymentRepository.UpdatePayment(s.db, *payment)
		return err
	}

	return s.db.WithTransaction(func(tx database.Database) error {
		payment.Status = enum.PaymentStatusPaid
		payment.RefCode = refCode
		if err := s.paymentRepository.UpdatePayment(tx, *payment); err != nil {
			return err
		}

		order.Status = enum.OrderStatusPaid
		if err := s.orderRepository.UpdateOrder(tx, *order); err != nil {
			return err
		}

		if err := s.orderRepository.CreateOrderStatusHistory(tx, entity.OrderStatusHistory{
			OrderID:     order.ID,
			Status:      enum.OrderStatusPaid,
			Note:        fmt.Sprintf("پرداخت آنلاین - کد پیگیری: %s", refCode),
			ChangedByID: order.UserID,
		}); err != nil {
			return err
		}

		go s.notifyOrderStatus(order.User.Phone, order.ID, enum.OrderStatusPaid)
		go s.createOrderNotification(order.UserID, order.ID, enum.OrderStatusPaid)
		return nil
	})
}

func (s *OrderService) CancelOrder(orderID uint) error {
	order, err := s.orderRepository.FindOrderByID(s.db, orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return exception.NotFoundError{Item: s.constants.Field.Order}
	}
	if order.Status == enum.OrderStatusShipped {
		return exception.ForbiddenError{Message: "cannot cancel a shipped order"}
	}
	if order.Status == enum.OrderStatusCancelled {
		return exception.ForbiddenError{Message: "order is already cancelled"}
	}

	return s.db.WithTransaction(func(tx database.Database) error {
		// Restore inventory
		for _, item := range order.Items {
			product, err := s.productRepository.FindProductByID(tx, item.ProductID)
			if err != nil {
				return err
			}
			if product != nil {
				product.Quantity += item.Count
				if err := s.productRepository.UpdateProduct(tx, *product); err != nil {
					return err
				}
			}
		}

		order.Status = enum.OrderStatusCancelled
		if err := s.orderRepository.UpdateOrder(tx, *order); err != nil {
			return err
		}

		return s.orderRepository.CreateOrderStatusHistory(tx, entity.OrderStatusHistory{
			OrderID: order.ID,
			Status:  enum.OrderStatusCancelled,
			Note:    "لغو شده توسط مدیر",
		})
	})
}

func (s *OrderService) FlagOrderRefund(orderID uint) error {
	order, err := s.orderRepository.FindOrderByID(s.db, orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return exception.NotFoundError{Item: s.constants.Field.Order}
	}

	order.RefundFlag = true
	return s.orderRepository.UpdateOrder(s.db, *order)
}

func (s *OrderService) GetOrderInstalments(orderID uint) ([]orderdto.InstalmentCredential, error) {
	instalments, err := s.instalmentRepository.GetInstalmentsByOrderID(s.db, orderID)
	if err != nil {
		return nil, err
	}

	var result []orderdto.InstalmentCredential
	for _, inst := range instalments {
		result = append(result, orderdto.InstalmentCredential{
			ID:      inst.ID,
			OrderID: inst.OrderID,
			Number:  inst.Number,
			Amount:  inst.Amount,
			DueDate: inst.DueDate,
			Status:  inst.Status,
			PaidAt:  inst.PaidAt,
		})
	}
	return result, nil
}

func (s *OrderService) sendOrderConfirmationEmail(order *entity.Order) {
	if order.User.Phone == "" {
		return
	}

	items := make([]map[string]interface{}, 0)
	for _, item := range order.Items {
		itemMap := map[string]interface{}{
			"ProductName": item.Product.Name,
			"Quantity":    item.Count,
			"Price":       fmt.Sprintf("%d", item.PriceSnapshot),
		}
		items = append(items, itemMap)
	}

	shippingAddr := "—"
	if order.Address != nil {
		shippingAddr = fmt.Sprintf("%s، %s - %s", order.Address.StreetAddress, order.Address.City.Name, order.Address.Province.Name)
	}

	paymentMethodStr := "—"
	switch order.PaymentMethod {
	case enum.PaymentMethodCash:
		paymentMethodStr = "نقدی"
	case enum.PaymentMethodInstallment:
		paymentMethodStr = "اقساطی"
	case enum.PaymentMethodWallet:
		paymentMethodStr = "کیف پول"
	}

	emailData := map[string]interface{}{
		"CustomerName":    order.User.Phone,
		"OrderID":         order.ID,
		"OrderDate":       order.CreatedAt.Format("2006-01-02"),
		"OrderStatus":     "در انتظار پرداخت",
		"Items":           items,
		"TotalAmount":     fmt.Sprintf("%d", order.TotalAmount),
		"ShippingAddress": shippingAddr,
		"PaymentMethod":   paymentMethodStr,
		"SupportEmail":    "support@mahoura.com",
	}

	_ = s.emailService.SendEmail(order.User.Phone, "تأیید سفارش - Mahoura", "order_confirmation/fa.html", emailData)
}
