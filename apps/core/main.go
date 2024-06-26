package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arthurdotwork/dreamtrader/core/internal/handler"
	"github.com/arthurdotwork/dreamtrader/core/internal/middleware"
	"github.com/arthurdotwork/dreamtrader/core/internal/service"
	"github.com/arthurdotwork/dreamtrader/core/internal/store"
	"github.com/arthurdotwork/dreamtrader/core/pkg/prom"
	"github.com/arthurdotwork/dreamtrader/core/pkg/psql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-done
		cancel()
	}()

	zerolog.DefaultContextLogger = func() *zerolog.Logger {
		logger := log.With().Caller().Logger()
		return &logger
	}()

	if err := run(ctx); err != nil {
		log.Error().Err(err).Msg("failed to run application")
		os.Exit(1)
	}
}

func run(parent context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// packages
	db, err := psql.Connect(
		ctx,
		env("DATABASE_USERNAME", "postgres"),
		env("DATABASE_PASSWORD", "postgres"),
		env("DATABASE_HOST", "localhost"),
		env("DATABASE_PORT", "5432"),
		env("DATABASE_NAME", "postgres"),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to database")
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	transactor := psql.NewTransactor(db(ctx))
	_ = transactor

	// stores
	userStore := store.NewUserStore(db)
	authStore := store.NewAuthStore(db)

	// services
	registerService := service.NewRegisterService(userStore)
	authService := service.NewAuthService(authStore, userStore)

	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"*"},
		AllowHeaders:    []string{"*"},
	}))
	router.Use(middleware.InstrumentedMiddleware())
	router.GET("/livez", handler.LivenessProbeHandler())
	router.GET("/readyz", handler.ReadinessProbeHandler(parent))

	apiRouter := router.Group("/api")
	v1Router := apiRouter.Group("/v1")

	v1Router.POST("/register", handler.RegisterHandler(registerService))
	v1Router.POST("/auth", handler.AuthenticateHandler(authService))
	v1Router.POST("/auth/verify", handler.VerifyAuthenticationHandler(authService))
	v1Router.POST("/auth/refresh", handler.RefreshAuthenticationHandler(authService))

	httpServer := &http.Server{
		Addr:              env("HTTP_ADDR", "0.0.0.0:8080"),
		Handler:           router,
		ReadHeaderTimeout: time.Second * 2,
	}

	errGroup, ctx := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to listen and serve: %w", err)
		}

		return nil
	})

	errGroup.Go(func() error {
		<-parent.Done()
		log.Debug().Msg("shutting down application")

		time.Sleep(time.Second * 5)
		cancel()

		if err := httpServer.Shutdown(context.WithoutCancel(ctx)); err != nil {
			return fmt.Errorf("failed to shutdown http server: %w", err)
		}

		return nil
	})

	errGroup.Go(func() error {
		if err := prom.Handler(ctx); err != nil {
			return fmt.Errorf("failed to run prometheus handler: %w", err)
		}

		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("failed to run application: %w", err)
	}

	log.Debug().Msg("application is shutting down")
	return nil
}

func env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
