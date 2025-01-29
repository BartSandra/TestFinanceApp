package main

import (
	"TestFinanceApp/configs"
	"TestFinanceApp/db"
	"TestFinanceApp/handler"
	"TestFinanceApp/repository"
	"TestFinanceApp/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {

	log := logrus.New()

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	log.SetLevel(logrus.InfoLevel)

	if err := godotenv.Load(); err != nil {
		log.Warn("No .env file found, using environment variables")
	}

	database, err := db.InitDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	repo := repository.NewTransactionRepository(database)
	svc := service.NewTransactionService(repo)
	h := handler.NewTransactionHandler(*svc)

	r := gin.Default()
	r.GET("/balance/:userID", h.GetBalance)
	r.POST("/deposit", h.Deposit)
	r.POST("/transfer", h.Transfer)
	r.GET("/transactions/:userID", h.GetTransactions)

	port := configs.GetEnv("PORT", "8080")

	log.Infof("Server is running on port %s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
