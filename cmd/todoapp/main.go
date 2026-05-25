package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth"
	core_config "github.com/daf32/golang-todoapp/internal/core/config"
	core_logger "github.com/daf32/golang-todoapp/internal/core/logger"
	core_mailer "github.com/daf32/golang-todoapp/internal/core/mailer"
	core_oauth "github.com/daf32/golang-todoapp/internal/core/oauth"
	core_pgx_pool "github.com/daf32/golang-todoapp/internal/core/repository/postgres/pool/pgx"
	core_http_cookie "github.com/daf32/golang-todoapp/internal/core/transport/http/cookie"
	core_http_middleware "github.com/daf32/golang-todoapp/internal/core/transport/http/middleware"
	core_ratelimit "github.com/daf32/golang-todoapp/internal/core/transport/http/middleware/ratelimit"
	core_http_server "github.com/daf32/golang-todoapp/internal/core/transport/http/server"
	auth_postgres_repository "github.com/daf32/golang-todoapp/internal/features/auth/repository/postgres"
	auth_service "github.com/daf32/golang-todoapp/internal/features/auth/service"
	auth_transport_http "github.com/daf32/golang-todoapp/internal/features/auth/transport/http"
	statistics_postgres_repository "github.com/daf32/golang-todoapp/internal/features/statistics/repository/postgres"
	statistics_service "github.com/daf32/golang-todoapp/internal/features/statistics/service"
	statistics_transport_http "github.com/daf32/golang-todoapp/internal/features/statistics/transport/http"
	tasks_postgres_repository "github.com/daf32/golang-todoapp/internal/features/tasks/repository/postgres"
	tasks_service "github.com/daf32/golang-todoapp/internal/features/tasks/service"
	tasks_transport_http "github.com/daf32/golang-todoapp/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/daf32/golang-todoapp/internal/features/users/repository/postgres"
	users_service "github.com/daf32/golang-todoapp/internal/features/users/service"
	users_transport_http "github.com/daf32/golang-todoapp/internal/features/users/transport/http"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	_ "github.com/daf32/golang-todoapp/docs"
)

func WithMiddleware(
	routes []core_http_server.Route,
	m ...core_http_middleware.Middleware,
) []core_http_server.Route {
	for i := range routes {
		routes[i].Middleware = append(m, routes[i].Middleware...)
	}

	return routes
}

// @title 		 Golang Todo API
// @version 	 1.0
// @description  Todo Application REST-API scheme
// @host 		 127.0.0.1:5050
// @BasePath 	 /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer {access_token}"
func main() {
	apiVersion := core_http_server.ApiVersion1

	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger: ", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", time.Local))

	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHanlder(usersService)

	logger.Debug("initialing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTaskHTTPHandler(tasksService)

	logger.Debug("initializing feature", zap.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService)

	authConfig := core_auth.NewConfigMust()
	smtpConfig := core_mailer.NewConfigMust()
	oauthConfig := core_oauth.NewConfigMust()

	mailer := core_mailer.NewSMTPMailer(smtpConfig)

	googleProvider, err := core_oauth.NewGoogleProvider(
		ctx,
		oauthConfig.Google,
		cfg.AppBaseURL,
	)
	if err != nil {
		logger.Fatal("failed to set up google OAuth provider", zap.Error(err))
	}

	loginRL := core_ratelimit.NewMemoryLimiter(
		rate.Every(10*time.Second), 5, 10*time.Minute,
	)
	registerRL := core_ratelimit.NewMemoryLimiter(
		rate.Every(20*time.Second), 3, 10*time.Minute,
	)
	resendRL := core_ratelimit.NewMemoryLimiter(
		rate.Every(60*time.Second), 3, 10*time.Minute,
	)
	refreshRL := core_ratelimit.NewMemoryLimiter(
		rate.Every(time.Second), 30, 10*time.Minute,
	)

	logger.Debug("initializing feature", zap.String("feature", "auth"))
	authRepository := auth_postgres_repository.NewAuthRepository(pool)
	authService := auth_service.NewAuthService(
		authRepository,
		usersRepository,
		mailer,
		logger,
		authConfig.JWTSecret,
		authConfig.AccessTokenExpiry,
		authConfig.RefreshTokenExpiry,
		authConfig.EmailConfirmationTokenExpiry,
		[]core_oauth.Provider{googleProvider},
	)
	cookieManager := core_http_cookie.NewManager(cfg.CookieSecure)
	authTransportHTTP := auth_transport_http.NewAuthHTTPHandler(
		authService,
		apiVersion,
		cfg.AppBaseURL,
		cookieManager,
		auth_transport_http.RateLimiters{
			Login:    loginRL,
			Register: registerRL,
			Resend:   resendRL,
			Refresh:  refreshRL,
		},
	)

	logger.Debug("initializing HTTP server")
	httpConfig := core_http_server.NewConfigMust()
	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		logger,
		core_http_middleware.CORS(httpConfig.AllowedOrigins),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	authMW := core_http_middleware.Auth(authService)

	apiVersionRouterV1 := core_http_server.NewApiVersionRouter(apiVersion)
	apiVersionRouterV1.RegisterRoutes(WithMiddleware(usersTransportHTTP.Routes(), authMW)...)
	apiVersionRouterV1.RegisterRoutes(WithMiddleware(tasksTransportHTTP.Routes(), authMW)...)
	apiVersionRouterV1.RegisterRoutes(WithMiddleware(statisticsTransportHTTP.Routes(), authMW)...)
	apiVersionRouterV1.RegisterRoutes(authTransportHTTP.Routes()...)

	/*
		Example of usage apiVersionRouter V2 with separate Middlewares

		apiVersionRouterV2 := core_http_server.NewApiVersionRouter(
			core_http_server.ApiVersion2,
			core_http_middleware.Dummy("api v2 middleware"),
		)
		apiVersionRouterV2.RegisterRoutes(usersTransportHTTP.Routes()...)
	*/

	httpServer.RegisterAPIRouters(
		apiVersionRouterV1,
		// apiVersionRouterV2,
	)
	httpServer.RegisterSPA("frontend/dist")
	httpServer.RegisterSwagger()
	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
