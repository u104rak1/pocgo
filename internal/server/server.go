package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	accountApp "github.com/ucho456job/pocgo/internal/application/account"
	signupApp "github.com/ucho456job/pocgo/internal/application/authentication/signup"
	userApp "github.com/ucho456job/pocgo/internal/application/user"
	"github.com/ucho456job/pocgo/internal/config"
	authDomain "github.com/ucho456job/pocgo/internal/domain/authentication"
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/repository"
	"github.com/ucho456job/pocgo/internal/presentation/shared/validation"
	signupPre "github.com/ucho456job/pocgo/internal/presentation/signup"
	"github.com/uptrace/bun"
)

func Start() {
	db, err := config.LoadDB()
	if err != nil {
		panic(err)
	}
	defer config.CloseDB(db)

	e := setupEcho(db)

	startServer(e)
}

func setupEcho(db *bun.DB) *echo.Echo {
	e := echo.New()
	validation.SetupCustomValidation(e)

	/** Repository */
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthenticationRepository(db)
	accountRepo := repository.NewAccountRepository(db)

	/** Domain Service */
	verifyEmailUniqueServ := userDomain.NewVerifyEmailUniquenessService(userRepo)
	verifyAuthUniqueServ := authDomain.NewVerifyAuthenticationUniquenessService(authRepo)
	accessTokenServ := authDomain.NewAccessTokenService()

	/** Unit of Work */
	signupUW := repository.NewUnitOfWorkWithResult[signupApp.SignupDTO](db)

	/** Usecase */
	createUserUC := userApp.NewCreateUserUsecase(userRepo, authRepo, verifyEmailUniqueServ, verifyAuthUniqueServ)
	createAccountUC := accountApp.NewCreateAccountUsecase(accountRepo)
	signupUC := signupApp.NewSignupUsecase(createUserUC, createAccountUC, *accessTokenServ, signupUW)

	/** Handler */
	signupHandler := signupPre.NewSignupHandler(signupUC)

	v1 := e.Group("/v1")

	/** Authentication */
	v1.POST("/signup", signupHandler.Run)
	return e
}

func startServer(e *echo.Echo) {
	env := config.NewEnv()
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