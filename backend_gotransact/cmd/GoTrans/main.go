package main

import (
	"gotransact/apps/Accounts/handlers"
	accounts "gotransact/apps/Accounts/models"
	validators "gotransact/apps/Accounts/validator"
	base "gotransact/apps/Base/models"
	handler2 "gotransact/apps/transaction/handlers"
	"gotransact/logger"

	//	"gotransact/apps/transaction/models"
	transaction "gotransact/apps/transaction/models"
	"gotransact/apps/transaction/utils"
	validators2 "gotransact/apps/transaction/validators"
	"gotransact/config"
	"gotransact/middleware"

	//	"gotransact/pkg/db"
	"gotransact/pkg/db"
	database "gotransact/pkg/db"
	"gotransact/responses"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

func main() {
	config.Loadenv()
	database.InitDB("")
	validators.Init()
	validators2.InitValidation()
	logger.Init()

	database.DB.AutoMigrate(&base.Base{}, &accounts.User{}, &accounts.Company{}, &transaction.PaymentGateway{}, &transaction.TransactionRequest{}, &transaction.TransactionHistory{})
	// Define your routes here
	c := cron.New()
	c.AddFunc("@every 1h", func() {
		transactions := utils.FetchTransactionsLast24Hours()
		filePath, err := utils.GenerateExcel(transactions)
		if err != nil {
			log.Fatalf("failed to generate excel: %v", err)
		}
		utils.SendMailWithAttachment("sangadevshatwar143@gmail.com", filePath)
	})
	c.Start()
	r := gin.Default()

	r.POST("/api/signup", handlers.SignupHandler)

	r.POST("/api/signin", handlers.LoginHandler)
	r.GET("/api/confirm_payment", handler2.ConfirmationPayment)
	//r.POST("/postpayment", handler2.PostPayment)
	auth := r.Group("/protected")
	{
		auth.Use(middleware.AuthMiddleware())
		auth.POST("/postpayment", handler2.PostPayment)
		auth.POST("/Logout", handlers.LogoutHandler)

	}
	r.POST("/api/create", func(c *gin.Context) {
		gatway := transaction.PaymentGateway{
			Slug:  "card",
			Label: "Card",
		}
		if err := db.DB.Create(&gatway).Error; err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			})
			return
		}
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
