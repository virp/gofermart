package handlers

import (
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/virp/gofermart/internal/handlers/middleware"
	"github.com/virp/gofermart/internal/handlers/order"
	"github.com/virp/gofermart/internal/handlers/user"
	"github.com/virp/gofermart/internal/handlers/withdrawal"
	"github.com/virp/gofermart/internal/repository"
	"github.com/virp/gofermart/internal/usecase"
	"github.com/virp/gofermart/pkg/accrual"
	"github.com/virp/gofermart/pkg/hash"
	"github.com/virp/gofermart/pkg/web"
	"go.uber.org/zap"
)

type Config struct {
	Shutdown                chan os.Signal
	DB                      *sqlx.DB
	Log                     *zap.SugaredLogger
	AppSecret               string
	AppUserAuthCookieName   string
	AccrualSystem           accrual.SDK
	OrderStatusWorkersCount int
}

func New(cfg Config) http.Handler {
	app := web.NewApp(
		cfg.Shutdown,
		middleware.Logger(cfg.Log),
		middleware.Errors(cfg.Log),
		middleware.Panics(),
	)

	authMiddleware := middleware.Auth(cfg.AppUserAuthCookieName, cfg.AppSecret)

	userRepository := repository.NewUserRepository(cfg.DB)
	passwordHash := hash.NewPasswordHash()
	userUseCase := usecase.NewUserUseCase(userRepository, passwordHash)

	userHandlers := user.Handlers{
		User:                  userUseCase,
		AppSecret:             cfg.AppSecret,
		AppUserAuthCookieName: cfg.AppUserAuthCookieName,
	}
	app.Handle(http.MethodPost, "/api/user/register", userHandlers.Register)
	app.Handle(http.MethodPost, "/api/user/login", userHandlers.Login)
	app.Handle(http.MethodGet, "/api/user/balance", userHandlers.Balance, authMiddleware)

	orderRepository := repository.NewOrderRepository(cfg.DB)
	ordersQueue := make(chan int)
	orderUseCase := usecase.NewOrderUseCase(orderRepository, ordersQueue)
	_ = usecase.NewOrderStatusUseCase(orderRepository, cfg.AccrualSystem, ordersQueue, cfg.OrderStatusWorkersCount)

	orderHandlers := order.Handlers{
		Order: orderUseCase,
	}
	app.Handle(http.MethodPost, "/api/user/orders", orderHandlers.Upload, authMiddleware)
	app.Handle(http.MethodGet, "/api/user/orders", orderHandlers.List, authMiddleware)

	withdrawalRepository := repository.NewWithdrawalRepository(cfg.DB)
	withdrawalUseCase := usecase.NewWithdrawalUseCase(
		withdrawalRepository,
		userRepository,
	)

	withdrawalHandlers := withdrawal.Handlers{
		Withdrawal: withdrawalUseCase,
	}
	app.Handle(http.MethodPost, "/api/user/balance/withdraw", withdrawalHandlers.Withdraw, authMiddleware)
	app.Handle(http.MethodGet, "/api/user/withdrawals", withdrawalHandlers.List, authMiddleware)

	return app
}
