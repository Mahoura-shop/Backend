//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/Mahoura-shop/Backend/internal/application/service"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/domain/communication"
	domainLogger "github.com/Mahoura-shop/Backend/internal/domain/logger"
	domainS3 "github.com/Mahoura-shop/Backend/internal/domain/s3"
	domainPostgres "github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	domainRedis "github.com/Mahoura-shop/Backend/internal/domain/repository/redis"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/communication/email"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/communication/sms"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	infraJWT "github.com/Mahoura-shop/Backend/internal/infrastructure/jwt"
	infraLocalization "github.com/Mahoura-shop/Backend/internal/infrastructure/localization"
	infraLogger "github.com/Mahoura-shop/Backend/internal/infrastructure/logger"
	infraStorage "github.com/Mahoura-shop/Backend/internal/infrastructure/storage"
	infraPostgres "github.com/Mahoura-shop/Backend/internal/infrastructure/repository/postgres"
	infraRedis "github.com/Mahoura-shop/Backend/internal/infrastructure/repository/redis"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/seed"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller/v1/address"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller/v1/category"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller/v1/health"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller/v1/coupon"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller/v1/currency"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller/v1/brand"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller/v1/product"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller/v1/review"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller/v1/test"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller/v1/user"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller/v1/cart"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller/v1/order"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller/v1/wishlist"
	upgraderequest "github.com/Mahoura-shop/Backend/internal/presentation/controller/v1/upgrade_request"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller/v1/role"
	returnctrl "github.com/Mahoura-shop/Backend/internal/presentation/controller/v1/return"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller/v1/contact"
	notificationctrl "github.com/Mahoura-shop/Backend/internal/presentation/controller/v1/notification"
	adminlog "github.com/Mahoura-shop/Backend/internal/presentation/controller/v1/admin_log"
	"github.com/Mahoura-shop/Backend/internal/presentation/middleware"
	"github.com/google/wire"
)

var DatabaseProviderSet = wire.NewSet(
	database.NewPostgresDatabase,
	database.NewRedisDatabase,
	wire.Bind(new(database.Database), new(*database.PostgresDatabase)),
	wire.Bind(new(database.Cache), new(*database.RedisDatabase)),
	wire.Struct(new(Database), "*"),
)

var RepositoryProviderSet = wire.NewSet(
	infraPostgres.NewRoleRepository,
	infraPostgres.NewUserRepository,
	infraPostgres.NewAddressRepository,
	infraPostgres.NewCategoryRepository,
	infraPostgres.NewCurrencyRepository,
	infraPostgres.NewBrandRepository,
	infraPostgres.NewProductRepository,
	infraPostgres.NewProductImageRepository,
	infraPostgres.NewProductVisitRepository,
	infraPostgres.NewReviewRepository,
	infraPostgres.NewCouponRepository,
	infraPostgres.NewWishlistRepository,
	infraPostgres.NewWalletRepository,
	infraPostgres.NewCartRepository,
	infraPostgres.NewOrderRepository,
	infraPostgres.NewTransactionRepository,
	infraPostgres.NewPaymentRepository,
	infraPostgres.NewInstalmentRepository,
	infraPostgres.NewUpgradeRequestRepository,
	infraPostgres.NewUserAuditLogRepository,
	infraPostgres.NewReturnRepository,
	infraPostgres.NewContactMessageRepository,
	infraPostgres.NewNotificationRepository,
	infraPostgres.NewAdminActivityLogRepository,
	infraRedis.NewUserCacheRepository,
	wire.Bind(new(domainPostgres.UserRepository), new(*infraPostgres.UserRepository)),
	wire.Bind(new(domainPostgres.AddressRepository), new(*infraPostgres.AddressRepository)),
	wire.Bind(new(domainRedis.UserCacheRepository), new(*infraRedis.UserCacheRepository)),
	wire.Bind(new(domainPostgres.CategoryRepository), new(*infraPostgres.CategoryRepository)),
	wire.Bind(new(domainPostgres.CurrencyRepository), new(*infraPostgres.CurrencyRepository)),
	wire.Bind(new(domainPostgres.BrandRepository), new(*infraPostgres.BrandRepository)),
	wire.Bind(new(domainPostgres.ProductRepository), new(*infraPostgres.ProductRepository)),
	wire.Bind(new(domainPostgres.ProductImageRepository), new(*infraPostgres.ProductImageRepository)),
	wire.Bind(new(domainPostgres.ProductVisitRepository), new(*infraPostgres.ProductVisitRepository)),
	wire.Bind(new(domainPostgres.ReviewRepository), new(*infraPostgres.ReviewRepository)),
	wire.Bind(new(domainPostgres.CouponRepository), new(*infraPostgres.CouponRepository)),
	wire.Bind(new(domainPostgres.WishlistRepository), new(*infraPostgres.WishlistRepository)),
	wire.Bind(new(domainPostgres.WalletRepository), new(*infraPostgres.WalletRepository)),
	wire.Bind(new(domainPostgres.CartRepository), new(*infraPostgres.CartRepository)),
	wire.Bind(new(domainPostgres.OrderRepository), new(*infraPostgres.OrderRepository)),
	wire.Bind(new(domainPostgres.TransactionRepository), new(*infraPostgres.TransactionRepository)),
	wire.Bind(new(domainPostgres.PaymentRepository), new(*infraPostgres.PaymentRepository)),
	wire.Bind(new(domainPostgres.InstalmentRepository), new(*infraPostgres.InstalmentRepository)),
	wire.Bind(new(domainPostgres.UpgradeRequestRepository), new(*infraPostgres.UpgradeRequestRepository)),
	wire.Bind(new(domainPostgres.UserAuditLogRepository), new(*infraPostgres.UserAuditLogRepository)),
	wire.Bind(new(domainPostgres.RoleRepository), new(*infraPostgres.RoleRepository)),
	wire.Bind(new(domainPostgres.ReturnRepository), new(*infraPostgres.ReturnRepository)),
	wire.Bind(new(domainPostgres.ContactMessageRepository), new(*infraPostgres.ContactMessageRepository)),
	wire.Bind(new(domainPostgres.NotificationRepository), new(*infraPostgres.NotificationRepository)),
	wire.Bind(new(domainPostgres.AdminActivityLogRepository), new(*infraPostgres.AdminActivityLogRepository)),
)

var ServiceProviderSet = wire.NewSet(
	wire.Struct(new(service.UserServiceDeps), "*"),
	wire.Struct(new(service.CategoryServiceDeps), "*"),
	wire.Struct(new(service.CurrencyServiceDeps), "*"),
	wire.Struct(new(service.BrandServiceDeps), "*"),
	wire.Struct(new(service.ProductServiceDeps), "*"),
	wire.Struct(new(service.ReviewServiceDeps), "*"),
	wire.Struct(new(service.CouponServiceDeps), "*"),
	wire.Struct(new(service.WishlistServiceDeps), "*"),
	wire.Struct(new(service.CartServiceDeps), "*"),
	wire.Struct(new(service.OrderServiceDeps), "*"),
	wire.Struct(new(service.UpgradeRequestServiceDeps), "*"),
	wire.Struct(new(service.RoleServiceDeps), "*"),
	wire.Struct(new(service.ReturnServiceDeps), "*"),
	wire.Struct(new(service.ContactMessageServiceDeps), "*"),
	wire.Struct(new(service.NotificationServiceDeps), "*"),
	wire.Struct(new(service.AdminLogServiceDeps), "*"),
	service.NewUserService,
	service.NewOTPService,
	sms.NewSMSService,
	email.NewEmailService,
	service.NewJWTService,
	service.NewAddressService,
	service.NewTestService,
	service.NewCategoryService,
	service.NewCurrencyService,
	service.NewBrandService,
	service.NewProductService,
	service.NewReviewService,
	service.NewCouponService,
	service.NewWishlistService,
	service.NewCartService,
	service.NewOrderService,
	service.NewPaymentService,
	service.NewUpgradeRequestService,
	service.NewRoleService,
	service.NewReturnService,
	service.NewContactMessageService,
	service.NewNotificationService,
	service.NewAdminLogService,
	wire.Bind(new(usecase.UserService), new(*service.UserService)),
	wire.Bind(new(usecase.OTPService), new(*service.OTPService)),
	wire.Bind(new(communication.SMSService), new(*sms.SMSService)),
	wire.Bind(new(communication.EmailService), new(*email.EmailService)),
	wire.Bind(new(usecase.JWTService), new(*service.JWTService)),
	wire.Bind(new(usecase.AddressService), new(*service.AddressService)),
	wire.Bind(new(usecase.TestService), new(*service.TestService)),
	wire.Bind(new(usecase.CategoryService), new(*service.CategoryService)),
	wire.Bind(new(usecase.CurrencyService), new(*service.CurrencyService)),
	wire.Bind(new(usecase.BrandService), new(*service.BrandService)),
	wire.Bind(new(usecase.ProductService), new(*service.ProductService)),
	wire.Bind(new(usecase.ReviewService), new(*service.ReviewService)),
	wire.Bind(new(usecase.CouponService), new(*service.CouponService)),
	wire.Bind(new(usecase.WishlistService), new(*service.WishlistService)),
	wire.Bind(new(usecase.CartService), new(*service.CartService)),
	wire.Bind(new(usecase.OrderService), new(*service.OrderService)),
	wire.Bind(new(usecase.PaymentService), new(*service.PaymentService)),
	wire.Bind(new(usecase.UpgradeRequestService), new(*service.UpgradeRequestService)),
	wire.Bind(new(usecase.RoleService), new(*service.RoleService)),
	wire.Bind(new(usecase.ReturnService), new(*service.ReturnService)),
	wire.Bind(new(usecase.ContactMessageService), new(*service.ContactMessageService)),
	wire.Bind(new(usecase.NotificationService), new(*service.NotificationService)),
	wire.Bind(new(usecase.AdminLogService), new(*service.AdminLogService)),
)

var AdapterProviderSet = wire.NewSet(
	infraLocalization.NewTranslationService,
	infraLogger.NewLogger,
	infraStorage.NewS3Storage,
	infraJWT.NewJWTKeyManager,
	wire.Bind(new(domainLogger.Logger), new(*infraLogger.Logger)),
	wire.Bind(new(domainS3.S3Storage), new(*infraStorage.S3Storage)),
)

var GeneralControllerProviderSet = wire.NewSet(
	user.NewGeneralUserController,
	address.NewGeneralAddressController,
	product.NewGeneralProductController,
	review.NewGeneralReviewController,
	test.NewGeneralTestController,
	contact.NewGeneralContactController,
	health.NewHealthController,
	wire.Struct(new(GeneralControllers), "*"),
)

var CustomerControllerProviderSet = wire.NewSet(
	user.NewCustomerUserController,
	address.NewCustomerAddressController,
	cart.NewCustomerCartController,
	coupon.NewCustomerCouponController,
	order.NewCustomerOrderController,
	review.NewCustomerReviewController,
	upgraderequest.NewCustomerUpgradeRequestController,
	wishlist.NewCustomerWishlistController,
	returnctrl.NewCustomerReturnController,
	notificationctrl.NewCustomerNotificationController,
	wire.Struct(new(CustomerControllers), "*"),
)

var AdminControllerProviderSet = wire.NewSet(
	currency.NewAdminCurrencyController,
	category.NewAdminCategoryController,
	brand.NewAdminBrandController,
	product.NewAdminProductController,
	user.NewAdminUserController,
	order.NewAdminOrderController,
	upgraderequest.NewAdminUpgradeRequestController,
	role.NewAdminRoleController,
	returnctrl.NewAdminReturnController,
	contact.NewAdminContactController,
	review.NewAdminReviewController,
	adminlog.NewAdminLogController,
	wire.Struct(new(AdminControllers), "*"),
)

var ControllersProviderSet = wire.NewSet(
	wire.Struct(new(Controllers), "*"),
)

var MiddlewareProviderSet = wire.NewSet(
	middleware.NewAuthMiddleware,
	middleware.NewCorsMiddleware,
	middleware.NewRecovery,
	middleware.NewLocalization,
	middleware.NewRateLimit,
	middleware.NewLoggerMiddleware,
	middleware.NewAdminActivityLogMiddleware,
	wire.Struct(new(Middlewares), "*"),
)

var SeederProviderSet = wire.NewSet(
	seed.NewAddressSeeder,
	seed.NewAdminSeeder,
	seed.NewCurrencySeeder,
	seed.NewPermissionSeeder,
	wire.Struct(new(Seeds), "*"),
)

func ProvideConstants(container *bootstrap.Config) *bootstrap.Constants {
	return container.Constants
}

func ProvideLoggerConfig(container *bootstrap.Config) *bootstrap.Logger {
	return &container.Env.Logger
}

func ProvideStorageConfig(container *bootstrap.Config) *bootstrap.S3 {
	return &container.Env.Storage
}

func ProvideRateLimitConfig(container *bootstrap.Config) *bootstrap.RateLimit {
	return &container.Env.RateLimit
}

func ProvideDBConfig(container *bootstrap.Config) *bootstrap.Database {
	return &container.Env.PrimaryDB
}

func ProvideRDBConfig(container *bootstrap.Config) *bootstrap.Redis {
	return &container.Env.PrimaryRedis
}

func ProvideOTPConfig(container *bootstrap.Config) *bootstrap.OTP {
	return &container.Env.OTP
}

func ProvideSMSGatewayConfig(container *bootstrap.Config) *bootstrap.SMSGateway {
	return &container.Env.SMSGateway
}

func ProvideSMSTemplates(container *bootstrap.Config) *bootstrap.SMSTemplates {
	return &container.Constants.SMSTemplates
}

func ProvideJWTKeysPath(container *bootstrap.Config) *bootstrap.JWTKeysPath {
	return &container.Constants.JWTKeysPath
}

func ProvideEmailTemplates(container *bootstrap.Config) *bootstrap.EmailTemplates {
	return &container.Constants.EmailTemplates
}

func ProvidePaginationConfig(container *bootstrap.Config) *bootstrap.Pagination {
	return &container.Env.Pagination
}

func ProvideWebsocketSetting(container *bootstrap.Config) *bootstrap.WebsocketSetting {
	return &container.Env.WebsocketSetting
}

func ProvideEmailSenderAccount(container *bootstrap.Config) *bootstrap.EmailAccount {
	return &container.Env.EmailSenderAccount
}

func ProvideSuperAdminCredential(container *bootstrap.Config) *bootstrap.AdminCredentials {
	return &container.Env.Admins
}

func ProvideZarinpalConfig(container *bootstrap.Config) *bootstrap.Zarinpal {
	return &container.Env.Zarinpal
}

func ProvideRBACConfig(container *bootstrap.Config) *bootstrap.RBAC {
	return &container.Env.RBAC
}

var ProviderSet = wire.NewSet(
	DatabaseProviderSet,
	RepositoryProviderSet,
	ServiceProviderSet,
	AdapterProviderSet,

	GeneralControllerProviderSet,
	CustomerControllerProviderSet,
	AdminControllerProviderSet,
	ControllersProviderSet,

	MiddlewareProviderSet,
	SeederProviderSet,
	ProvideConstants,
	ProvideLoggerConfig,
	ProvideStorageConfig,
	ProvideRateLimitConfig,
	ProvideDBConfig,
	ProvideRDBConfig,
	ProvideOTPConfig,
	ProvideSMSGatewayConfig,
	ProvideSMSTemplates,
	ProvideEmailTemplates,
	ProvideJWTKeysPath,
	ProvidePaginationConfig,
	ProvideWebsocketSetting,
	ProvideEmailSenderAccount,
	ProvideSuperAdminCredential,
	ProvideZarinpalConfig,
	ProvideRBACConfig,
)

type Database struct {
	DB  database.Database
	RDB database.Cache
}

type GeneralControllers struct {
	UserController      *user.GeneralUserController
	AddressController   *address.GeneralAddressController
	ProductController   *product.GeneralProductController
	ReviewController    *review.GeneralReviewController
	TestController      *test.GeneralTestController
	ContactController   *contact.GeneralContactController
	HealthController    *health.HealthController
}

type CustomerControllers struct {
	UserController           *user.CustomerUserController
	AddressController        *address.CustomerAddressController
	CartController           *cart.CustomerCartController
	CouponController         *coupon.CustomerCouponController
	OrderController          *order.CustomerOrderController
	ReviewController         *review.CustomerReviewController
	UpgradeRequestController *upgraderequest.CustomerUpgradeRequestController
	WishlistController       *wishlist.CustomerWishlistController
	ReturnController         *returnctrl.CustomerReturnController
	NotificationController   *notificationctrl.CustomerNotificationController
}

type AdminControllers struct {
	CategoryController       *category.AdminCategoryController
	CurrencyController       *currency.AdminCurrencyController
	BrandController          *brand.AdminBrandController
	ProductController        *product.AdminProductController
	UserController           *user.AdminUserController
	OrderController          *order.AdminOrderController
	UpgradeRequestController *upgraderequest.AdminUpgradeRequestController
	RoleController           *role.AdminRoleController
	ReturnController         *returnctrl.AdminReturnController
	ContactController        *contact.AdminContactController
	ReviewController         *review.AdminReviewController
	AdminLogController       *adminlog.AdminLogController
}

type Controllers struct {
	General  *GeneralControllers
	Customer *CustomerControllers
	Admin    *AdminControllers
}

type Middlewares struct {
	Authentication  *middleware.AuthMiddleware
	CORS            *middleware.CORSMiddleware
	Recovery        *middleware.RecoveryMiddleware
	Localization    *middleware.LocalizationMiddleware
	RateLimit       *middleware.RateLimitMiddleware
	Logger          *middleware.LoggerMiddleware
	AdminActivityLog *middleware.AdminActivityLogMiddleware
}

type Seeds struct {
	AddressSeeder    *seed.AddressSeeder
	AdminSeeder      *seed.AdminSeeder
	CurrencySeeder   *seed.CurrencySeeder
	PermissionSeeder *seed.PermissionSeeder
}

type Application struct {
	Database    *Database
	Controllers *Controllers
	Middlewares *Middlewares
	Seeds       *Seeds
}

func NewApplication(
	database *Database,
	controllers *Controllers,
	middlewares *Middlewares,
	seeds *Seeds,
) *Application {
	return &Application{
		Database:    database,
		Controllers: controllers,
		Middlewares: middlewares,
		Seeds:       seeds,
	}
}

func InitializeApplication(container *bootstrap.Config) (*Application, error) {
	wire.Build(
		ProviderSet,
		NewApplication,
	)
	return &Application{}, nil
}
