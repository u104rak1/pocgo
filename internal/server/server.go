package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	accountApp "github.com/ucho456job/pocgo/internal/application/account"
	authApp "github.com/ucho456job/pocgo/internal/application/authentication"
	userApp "github.com/ucho456job/pocgo/internal/application/user"
	authDomain "github.com/ucho456job/pocgo/internal/domain/authentication"
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/internal/environment"
	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/config"
	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/repository"
	signinPre "github.com/ucho456job/pocgo/internal/presentation/signin"
	signupPre "github.com/ucho456job/pocgo/internal/presentation/signup"
	myMiddleware "github.com/ucho456job/pocgo/internal/server/middleware"
	"github.com/uptrace/bun"
)

func Start() {
	db, err := config.LoadDB()
	if err != nil {
		panic(err)
	}
	defer config.CloseDB(db)

	e := SetupEcho(db)

	startServer(e)
}

func SetupEcho(db *bun.DB) *echo.Echo {
	e := echo.New()
	e.Use(echoMiddleware.RequestID())
	myMiddleware.SetLoggerMiddleware(e)

	/** Repository */
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthenticationRepository(db)
	accountRepo := repository.NewAccountRepository(db)

	/** Domain Service */
	userServ := userDomain.NewService(userRepo)
	authServ := authDomain.NewService(authRepo, userRepo)

	/** Unit of Work */
	signupUW := repository.NewUnitOfWorkWithResult[authApp.SignupDTO](db)

	/** Usecase */
	createUserUC := userApp.NewCreateUserUsecase(userRepo, authRepo, userServ, authServ)
	createAccountUC := accountApp.NewCreateAccountUsecase(accountRepo)
	signupUC := authApp.NewSignupUsecase(createUserUC, createAccountUC, authServ, signupUW)
	signinUC := authApp.NewSigninUsecase(userRepo, authRepo, authServ)

	/** Handler */
	signupHandler := signupPre.NewSignupHandler(signupUC)
	signinHandler := signinPre.NewSigninHandler(signinUC)

	v1 := e.Group("/api/v1")

	/** Authentication Endpoint */
	v1.POST("/signup", signupHandler.Run)
	v1.POST("/signin", signinHandler.Run)

	/** Swagger */
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	return e
}

func startServer(e *echo.Echo) {
	env := environment.New()
	port := ":" + env.APP_PORT
	go func() {
		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
