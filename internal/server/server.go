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
	"github.com/u104rak1/pocgo/internal/infrastructure/jwt"
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

func SetupEcho(db *bun.DB) *echo.Echo {
	e := echo.New()

	repositories := setupRepository(db)
	domainServices := setupDomainServices(repositories)
	usecases := setupUsecases(db, repositories, domainServices)
	handlers := setupHandlers(usecases)

	/** Middleware */
	e.Use(echoMiddleware.RequestID())
	myMiddleware.SetLoggerMiddleware(e)
	authMiddleware := myMiddleware.AuthorizationMiddleware(repositories.jwt)

	/** Health Endpoint */
	healthHandler := healthPre.NewHealthHandler(db)
	e.GET("/", healthHandler.Run)

	v1 := e.Group("/api/v1")
	setupRoutes(v1, handlers, authMiddleware)

	/** Swagger */
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return e
}

type Repositories struct {
	user        userDomain.IUserRepository
	auth        authDomain.IAuthenticationRepository
	account     accountDomain.IAccountRepository
	transaction transactionDomain.ITransactionRepository
	jwt         authApp.IJWTService
}

func setupRepository(db *bun.DB) (repositories Repositories) {
	env := config.NewEnv()

	return Repositories{
		user:        repository.NewUserRepository(db),
		auth:        repository.NewAuthenticationRepository(db),
		account:     repository.NewAccountRepository(db),
		transaction: repository.NewTransactionRepository(db),
		jwt:         jwt.NewService([]byte(env.JWT_SECRET_KEY)),
	}
}

type DomainServices struct {
	user        userDomain.IUserService
	auth        authDomain.IAuthenticationService
	account     accountDomain.IAccountService
	transaction transactionDomain.ITransactionService
}

func setupDomainServices(r Repositories) DomainServices {
	return DomainServices{
		user:        userDomain.NewService(r.user),
		auth:        authDomain.NewService(r.auth, r.user),
		account:     accountDomain.NewService(r.account),
		transaction: transactionDomain.NewService(r.account, r.transaction),
	}
}

type Usecases struct {
	signupUC           authApp.ISignupUsecase
	signinUC           authApp.ISigninUsecase
	readUserUC         userApp.IReadUserUsecase
	createAccountUC    accountApp.ICreateAccountUsecase
	execTransactionUC  transactionApp.IExecuteTransactionUsecase
	listTransactionsUC transactionApp.IListTransactionsUsecase
}

func setupUsecases(db *bun.DB, r Repositories, ds DomainServices) Usecases {
	uow := repository.NewUnitOfWork(db)
	transactionUOW := repository.NewUnitOfWorkWithResult[transactionDomain.Transaction](db)

	return Usecases{
		signupUC:           authApp.NewSignupUsecase(r.user, r.auth, ds.user, ds.auth, r.jwt),
		signinUC:           authApp.NewSigninUsecase(ds.auth, r.jwt),
		readUserUC:         userApp.NewReadUserUsecase(ds.user),
		createAccountUC:    accountApp.NewCreateAccountUsecase(r.account, ds.account, ds.user, uow),
		execTransactionUC:  transactionApp.NewExecuteTransactionUsecase(ds.account, ds.transaction, transactionUOW),
		listTransactionsUC: transactionApp.NewListTransactionsUsecase(ds.account, ds.transaction),
	}
}

type Handlers struct {
	signupHandler           *signupPre.SignupHandler
	signinHandler           *signinPre.SigninHandler
	readMyProfHandler       *mePre.ReadMyProfileHandler
	createAccountHandler    *accountsPre.CreateAccountHandler
	execTransactionHandler  *transactionsPre.ExecuteTransactionHandler
	listTransactionsHandler *transactionsPre.ListTransactionsHandler
}

func setupHandlers(u Usecases) Handlers {
	return Handlers{
		signupHandler:           signupPre.NewSignupHandler(u.signupUC),
		signinHandler:           signinPre.NewSigninHandler(u.signinUC),
		readMyProfHandler:       mePre.NewReadMyProfileHandler(u.readUserUC),
		createAccountHandler:    accountsPre.NewCreateAccountHandler(u.createAccountUC),
		execTransactionHandler:  transactionsPre.NewExecuteTransactionHandler(u.execTransactionUC),
		listTransactionsHandler: transactionsPre.NewListTransactionsHandler(u.listTransactionsUC),
	}
}

func setupRoutes(e *echo.Group, h Handlers, authMiddleware echo.MiddlewareFunc) {
	/** Authentication Endpoint */
	e.POST("/signup", h.signupHandler.Run)
	e.POST("/signin", h.signinHandler.Run)

	/** User Endpoint */
	e.GET("/me", h.readMyProfHandler.Run, authMiddleware)

	/** Account Endpoint */
	e.POST("/me/accounts", h.createAccountHandler.Run, authMiddleware)

	/** Transaction Endpoint */
	e.POST("/me/accounts/:account_id/transactions", h.execTransactionHandler.Run, authMiddleware)
	e.GET("/me/accounts/:account_id/transactions", h.listTransactionsHandler.Run, authMiddleware)
}
