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
	accountApp "github.com/u104rak1/pocgo/internal/application/account"
	authApp "github.com/u104rak1/pocgo/internal/application/authentication"
	transactionApp "github.com/u104rak1/pocgo/internal/application/transaction"
	userApp "github.com/u104rak1/pocgo/internal/application/user"
	"github.com/u104rak1/pocgo/internal/config"
	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	authDomain "github.com/u104rak1/pocgo/internal/domain/authentication"
	transactionDomain "github.com/u104rak1/pocgo/internal/domain/transaction"
	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	"github.com/u104rak1/pocgo/internal/infrastructure/postgres/repository"
	healthPre "github.com/u104rak1/pocgo/internal/presentation/health"
	mePre "github.com/u104rak1/pocgo/internal/presentation/me"
	accountsPre "github.com/u104rak1/pocgo/internal/presentation/me/accounts"
	transactionsPre "github.com/u104rak1/pocgo/internal/presentation/me/accounts/transactions"
	signinPre "github.com/u104rak1/pocgo/internal/presentation/signin"
	signupPre "github.com/u104rak1/pocgo/internal/presentation/signup"
	myMiddleware "github.com/u104rak1/pocgo/internal/server/middleware"
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

	/** Repository */
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthenticationRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	/** Domain Service */
	userServ := userDomain.NewService(userRepo)
	authServ := authDomain.NewService(authRepo, userRepo)
	accountServ := accountDomain.NewService(accountRepo)
	transactionServ := transactionDomain.NewService(accountRepo, transactionRepo)

	/** Middleware */
	e.Use(echoMiddleware.RequestID())
	myMiddleware.SetLoggerMiddleware(e)
	env := config.NewEnv()
	authMiddleware := myMiddleware.AuthorizationMiddleware(authServ, []byte(env.JWT_SECRET_KEY))

	/** Unit of Work */
	unitOfWork := repository.NewUnitOfWork(db)
	transactionUOW := repository.NewUnitOfWorkWithResult[transactionDomain.Transaction](db)

	/** Usecase */
	readUserUC := userApp.NewReadUserUsecase(userRepo)
	createAccountUC := accountApp.NewCreateAccountUsecase(accountRepo, accountServ, userServ, unitOfWork)
	signupUC := authApp.NewSignupUsecase(userRepo, authRepo, userServ, authServ)
	signinUC := authApp.NewSigninUsecase(authServ)
	execTransactionUC := transactionApp.NewExecuteTransactionUsecase(accountServ, transactionServ, transactionUOW)
	listTransactionsUC := transactionApp.NewListTransactionsUsecase(accountServ, transactionServ)

	/** Handler */
	healthHandler := healthPre.NewHealthHandler(db)
	signupHandler := signupPre.NewSignupHandler(signupUC)
	signinHandler := signinPre.NewSigninHandler(signinUC)
	readMyProfHandler := mePre.NewReadMyProfileHandler(readUserUC)
	createAccountHandler := accountsPre.NewCreateAccountHandler(createAccountUC)
	execTransactionHandler := transactionsPre.NewExecuteTransactionHandler(execTransactionUC)
	listTransactionsHandler := transactionsPre.NewListTransactionsHandler(listTransactionsUC)

	/** Health Endpoint */
	e.GET("/", healthHandler.Run)

	v1 := e.Group("/api/v1")

	/** Authentication Endpoint */
	v1.POST("/signup", signupHandler.Run)
	v1.POST("/signin", signinHandler.Run)

	/** User Endpoint */
	v1.GET("/me", readMyProfHandler.Run, authMiddleware)

	/** Account Endpoint */
	v1.POST("/me/accounts", createAccountHandler.Run, authMiddleware)

	/** Transaction Endpoint */
	v1.POST("/me/accounts/:account_id/transactions", execTransactionHandler.Run, authMiddleware)
	v1.GET("/me/accounts/:account_id/transactions", listTransactionsHandler.Run, authMiddleware)

	/** Swagger */
	e.GET("/swagger/*", echoSwagger.WrapHandler)
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
