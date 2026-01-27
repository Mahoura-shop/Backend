//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/Mahoura-shop/Backend/internal/application/service"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/domain/communication"
	domainLogger "github.com/Mahoura-shop/Backend/internal/domain/logger"
	"github.com/Mahoura-shop/Backend/internal/domain/message"
	domainPostgres "github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	domainRedis "github.com/Mahoura-shop/Backend/internal/domain/repository/redis"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/communication/email"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/communication/sms"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	infraJWT "github.com/Mahoura-shop/Backend/internal/infrastructure/jwt"
	infraLocalization "github.com/Mahoura-shop/Backend/internal/infrastructure/localization"
	infraLogger "github.com/Mahoura-shop/Backend/internal/infrastructure/logger"
	infraRabbitMQ "github.com/Mahoura-shop/Backend/internal/infrastructure/rabbitmq"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/rabbitmq/consumer"
	infraPostgres "github.com/Mahoura-shop/Backend/internal/infrastructure/repository/postgres"
	infraRedis "github.com/Mahoura-shop/Backend/internal/infrastructure/repository/redis"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/seed"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/websocket"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller/v1/address"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller/v1/user"
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
	infraPostgres.NewUserRepository,
	infraPostgres.NewAddressRepository,
	infraRedis.NewUserCacheRepository,
	infraPostgres.NewNotificationRepository,
	wire.Bind(new(domainPostgres.UserRepository), new(*infraPostgres.UserRepository)),
	wire.Bind(new(domainPostgres.AddressRepository), new(*infraPostgres.AddressRepository)),
	wire.Bind(new(domainRedis.UserCacheRepository), new(*infraRedis.UserCacheRepository)),
	wire.Bind(new(domainPostgres.NotificationRepository), new(*infraPostgres.NotificationRepository)),
)

var ServiceProviderSet = wire.NewSet(
	wire.Struct(new(service.UserServiceDeps), "*"),
	service.NewUserService,
	service.NewOTPService,
	sms.NewSMSService,
	email.NewEmailService,
	service.NewJWTService,
	service.NewAddressService,
	wire.Bind(new(usecase.UserService), new(*service.UserService)),
	wire.Bind(new(usecase.OTPService), new(*service.OTPService)),
	wire.Bind(new(communication.SMSService), new(*sms.SMSService)),
	wire.Bind(new(communication.EmailService), new(*email.EmailService)),
	wire.Bind(new(usecase.JWTService), new(*service.JWTService)),
	wire.Bind(new(usecase.AddressService), new(*service.AddressService)),
)

var AdapterProviderSet = wire.NewSet(
	infraLocalization.NewTranslationService,
	infraLogger.NewLogger,
	infraJWT.NewJWTKeyManager,
	infraRabbitMQ.NewRabbitMQ,
	wire.Bind(new(domainLogger.Logger), new(*infraLogger.Logger)),
	wire.Bind(new(message.Broker), new(*infraRabbitMQ.RabbitMQ)),
)

var GeneralControllerProviderSet = wire.NewSet(
	user.NewGeneralUserController,
	address.NewGeneralAddressController,
	wire.Struct(new(GeneralControllers), "*"),
)

var CustomerControllerProviderSet = wire.NewSet(
	user.NewCustomerUserController,
	address.NewCustomerAddressController,
	wire.Struct(new(CustomerControllers), "*"),
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
	wire.Struct(new(Middlewares), "*"),
)

var SeederProviderSet = wire.NewSet(
	seed.NewAddressSeeder,
	seed.NewNotificationTypeSeeder,
	seed.NewRoleSeeder,
	wire.Struct(new(Seeds), "*"),
)

var ConsumerProviderSet = wire.NewSet(
	consumer.NewRegisterConsumer,
	consumer.NewPushConsumer,
	consumer.NewEmailConsumer,
	consumer.NewSendNotificationConsumer,
	wire.Struct(new(Consumers), "*"),
)

func ProvideConstants(container *bootstrap.Config) *bootstrap.Constants {
	return container.Constants
}

func ProvideLoggerConfig(container *bootstrap.Config) *bootstrap.Logger {
	return &container.Env.Logger
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

func ProvideStorageConfig(container *bootstrap.Config) *bootstrap.S3 {
	return &container.Env.Storage
}

func ProvideWebsocketSetting(container *bootstrap.Config) *bootstrap.WebsocketSetting {
	return &container.Env.WebsocketSetting
}

func ProvideEmailSenderAccount(container *bootstrap.Config) *bootstrap.EmailAccount {
	return &container.Env.EmailSenderAccount
}

func ProvideSuperAdminCredential(container *bootstrap.Config) *bootstrap.AdminCredentials {
	return &container.Env.SuperAdmin
}

func ProvideRabbitMQConfig(container *bootstrap.Config) *bootstrap.RabbitMQ {
	return &container.Env.RabbitMQ
}

func ProvideRabbitMQConstants(container *bootstrap.Config) *bootstrap.RabbitMQConstants {
	return &container.Constants.RabbitMQ
}

func ProvideWebSocketHub(container *bootstrap.Config) *websocket.Hub {
    return websocket.NewHub()
}

var ProviderSet = wire.NewSet(
	DatabaseProviderSet,
	RepositoryProviderSet,
	ServiceProviderSet,
	AdapterProviderSet,
	GeneralControllerProviderSet,
	CustomerControllerProviderSet,
	ControllersProviderSet,
	MiddlewareProviderSet,
	SeederProviderSet,
	ConsumerProviderSet,
	ProvideWebSocketHub,
	ProvideConstants,
	ProvideLoggerConfig,
	ProvideRateLimitConfig,
	ProvideDBConfig,
	ProvideRDBConfig,
	ProvideOTPConfig,
	ProvideSMSGatewayConfig,
	ProvideSMSTemplates,
	ProvideEmailTemplates,
	ProvideJWTKeysPath,
	ProvidePaginationConfig,
	ProvideStorageConfig,
	ProvideWebsocketSetting,
	ProvideEmailSenderAccount,
	ProvideSuperAdminCredential,
	ProvideRabbitMQConfig,
	ProvideRabbitMQConstants,
)

type Database struct {
	DB  database.Database
	RDB database.Cache
}

type GeneralControllers struct {
	UserController    *user.GeneralUserController
	AddressController *address.GeneralAddressController
}

type CustomerControllers struct {
	UserController    *user.CustomerUserController
	AddressController *address.CustomerAddressController
}

type Controllers struct {
	General  *GeneralControllers
	Customer *CustomerControllers
}

type Middlewares struct {
	Authentication *middleware.AuthMiddleware
	CORS           *middleware.CORSMiddleware
	Recovery       *middleware.RecoveryMiddleware
	Localization   *middleware.LocalizationMiddleware
	RateLimit      *middleware.RateLimitMiddleware
	Logger         *middleware.LoggerMiddleware
}

type Seeds struct {
	AddressSeeder          *seed.AddressSeeder
	NotificationTypeSeeder *seed.NotificationTypeSeeder
	RoleSeeder             *seed.RoleSeeder
}

type Consumers struct {
	Register     *consumer.RegisterConsumer
	Push         *consumer.PushConsumer
	Email        *consumer.EmailConsumer
	Notification *consumer.SendNotificationConsumer
}

type Application struct {
	Database    *Database
	Controllers *Controllers
	Middlewares *Middlewares
	Seeds       *Seeds
	Consumers   *Consumers
	Hub         *websocket.Hub
}

func NewApplication(
	database *Database,
	controllers *Controllers,
	middlewares *Middlewares,
	seeds *Seeds,
	consumers *Consumers,
	hub *websocket.Hub,
) *Application {
	return &Application{
		Database:    database,
		Controllers: controllers,
		Middlewares: middlewares,
		Seeds:       seeds,
		Consumers:   consumers,
		Hub:         hub,
	}
}

func InitializeApplication(container *bootstrap.Config) (*Application, error) {
	wire.Build(
		ProviderSet,
		NewApplication,
	)
	return &Application{}, nil
}
