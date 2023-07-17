package app

import (
	"context"
	"github.com/MrTomSawyer/loyalty-system/internal/app/config"
	"github.com/MrTomSawyer/loyalty-system/internal/app/handler"
	"github.com/MrTomSawyer/loyalty-system/internal/app/interfaces"
	"github.com/MrTomSawyer/loyalty-system/internal/app/logger"
	"github.com/MrTomSawyer/loyalty-system/internal/app/middleware"
	"github.com/MrTomSawyer/loyalty-system/internal/app/repository"
	"github.com/MrTomSawyer/loyalty-system/internal/app/repository/postgres"
	"github.com/MrTomSawyer/loyalty-system/internal/app/service"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	app                  *fiber.App
	Config               *config.Config
	UserRepository       interfaces.UserRepository
	OrderRepository      interfaces.OrderRepository
	WithdrawalRepository interfaces.WithdrawalRepository
	UserService          interfaces.UserService
	OrderService         interfaces.OrderService
	WithdrawalService    interfaces.WithdrawalService
	authService          interfaces.Auth
	Handler              interfaces.Handler
	dataBase             *repository.PostgresDatabase
	middlewares          interfaces.Middlewares
}

func NewServer() *Server {
	return &Server{
		Config: config.NewConfig(),
	}
}

func (s *Server) Run() error {
	s.app = fiber.New(fiber.Config{
		AppName: "GopherMart Loyalty System",
	})

	s.Config.InitConfig()

	err := logger.InitLogger(s.Config.Environment, s.Config.LogLevel)
	if err != nil {
		log.Fatal("Failed to init logger:", err)
		return err
	}

	ctx := context.Background()
	pgPool, err := pgxpool.New(ctx, s.Config.PostgresDSN)
	if err != nil {
		logger.Log.Fatalf("failed to init repositories")
		return err
	}

	s.CreateDataBase(ctx, pgPool, s.Config)
	if err != nil {
		logger.Log.Fatalf("failed to configure database")
	}

	s.InitRepositories(ctx, pgPool)
	s.InitServices()
	s.InitHandler()
	s.InitMiddlewares()
	s.InitRouter()
	err = s.app.Listen(s.Config.Port)
	if err != nil {
		logger.Log.Errorf("failed to listen port")
		return err
	}
	return nil
}

func (s *Server) InitRepositories(ctx context.Context, pool *pgxpool.Pool) {
	s.UserRepository = postgres.NewUserRepository(ctx, pool, s.Config.UserTableName)
	s.OrderRepository = postgres.NewOrderRepository(ctx, pool, s.Config.OrderTableName)
	s.WithdrawalRepository = postgres.NewWithdrawalRepository(ctx, pool)
}

func (s *Server) InitServices() {
	s.authService = service.NewAuthService(s.Config.SecretKey, s.Config.TokenExp)
	s.UserService = service.NewUserService(s.Config, s.UserRepository, s.authService)
	s.OrderService = service.NewOrderService(s.Config, s.OrderRepository)
	s.WithdrawalService = service.NewWithdrawalService(s.WithdrawalRepository)
}

func (s *Server) InitHandler() {
	s.Handler = handler.NewHandler(s.UserService, s.OrderService, s.WithdrawalService)
}

func (s *Server) InitMiddlewares() {
	s.middlewares = middleware.NewMiddlewares(s.authService)
}

func (s *Server) CreateDataBase(ctx context.Context, pool *pgxpool.Pool, config *config.Config) {
	s.dataBase = repository.NewDatabase(ctx, pool, config)
	err := s.dataBase.ConfigDataBase()
	if err != nil {
		logger.Log.Fatalf("failed to init a database")
	}
}

func (s *Server) InitRouter() {
	api := s.app.Group("/api", s.middlewares.LogReqResInfo)

	user := api.Group("/user")
	{
		user.Post("/register", s.Handler.CreateUser)
		user.Post("/login", s.Handler.Login)
		user.Get("/withdrawals", s.middlewares.AuthMiddleware, s.Handler.GetWithdrawals)

		balance := user.Group("/balance", s.middlewares.AuthMiddleware)
		{
			balance.Get("/", s.Handler.GetUserBalance)
			balance.Post("/withdraw", s.Handler.Withdraw)
		}

		orders := user.Group("/orders", s.middlewares.AuthMiddleware)
		{
			orders.Post("/", s.Handler.CreateOrder)
			orders.Get("/", s.Handler.GetAllOrders)
		}
	}
}

func (s *Server) Shutdown() {
	err := s.app.Shutdown()
	if err != nil {
		logger.Log.Errorf("failed to shutdown server: %v", err)
	}
}
