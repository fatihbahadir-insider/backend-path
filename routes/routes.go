package routes

import (
	"backend-path/app/controllers"
	"backend-path/app/metrics"
	"backend-path/app/middlewares"
	"backend-path/app/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Setup(app *fiber.App) {
	apiRoute := app.Group("/api/v1", middlewares.JwtMiddleware)

    app.Get("/metrics", adaptor.HTTPHandler(
        promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{}),
    ))
	auth := app.Group("/api/v1/auth", middlewares.SetupAuthRateLimiter())
	authController := controllers.NewAuthController()
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/refresh", authController.RefreshToken)

	users := apiRoute.Group("/users")
	userController := controllers.NewUserController()
	users.Get("/", middlewares.Role(models.RoleAdmin, models.RoleMod), userController.GetAll)
	users.Get("/:id", middlewares.Role(models.RoleAdmin, models.RoleMod), userController.GetByID)
	users.Put("/:id", middlewares.Role(models.RoleAdmin), userController.Update)
	users.Delete("/:id", middlewares.Role(models.RoleAdmin), userController.Delete)

	balances := apiRoute.Group("/balances")
	balanceController := controllers.NewBalanceController()

	balances.Get("/current", balanceController.GetCurrent)
	balances.Get("/historical", balanceController.GetHistorical)
	balances.Get("/at-time", balanceController.GetAtTime)

	transactions := apiRoute.Group("/transactions")
	transactionController := controllers.NewTransactionController()
	transactions.Post("/credit", transactionController.Credit)
	transactions.Post("/debit", transactionController.Debit)
	transactions.Post("/transfer", transactionController.Transfer)
	transactions.Get("/history", transactionController.GetHistory)
	transactions.Get("/stats", middlewares.Role(models.RoleAdmin), transactionController.GetStats)
	transactions.Get("/:id", transactionController.GetByID)
}