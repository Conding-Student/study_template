package main

import (
	"fmt"
	"log"
	"os"

	routers "github.com/Conding-Student/study_template/pkg/routers"
	middleware "github.com/Conding-Student/study_template/pkg/utils"
	"github.com/Conding-Student/study_template/pkg/utils/go-utils/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error Loading Env File: ", err)
	}
	envi := os.Getenv("ENVIRONMENT")

	err = godotenv.Load(fmt.Sprintf(".env-%v", envi)) //
	if err != nil {
		log.Fatal("Error Loading Env File 2nd: ", err)
	}

	// Initialize DB here
	// Initialize DB
	dbDriver := os.Getenv("DB_DRIVER")
	switch dbDriver {
	case "postgres":
		database.PostgreSQLConnect(
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DATABASE_NAME"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_SSLMODE"),
			os.Getenv("DB_TIMEZONE"),
		)
	case "mysql":
		database.MySQLConnect(
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_NAME"),
		)
	case "sqlite":
		database.SQLiteConnect(os.Getenv("DB_NAME"))
	default:
		log.Fatal("❌ Unsupported DB driver:", dbDriver)
	}

	if database.Err != nil {
		log.Fatal("❌ Failed to connect to DB:", database.Err)
	} else {
		log.Println("✅ Database connected successfully!")
	}
	// Declare & initialize fiber
	app := fiber.New(fiber.Config{
		UnescapePath: true,
	})

	// For GoRoutine implementation
	// appb := fiber.New(fiber.Config{
	// 	UnescapePath: true,
	// })

	// Configure application CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization", // ✅ Add Authorization
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	// For GoRoutine implementation
	// appb.Use(cors.New(cors.Config{
	// 	AllowOrigins: "*",
	// 	AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	// }))

	// Declare & initialize logger
	app.Use(logger.New())

	// For GoRoutine implementation
	// appb.Use(logger.New())

	// Declare & initialize routes
	routers.SetupPublicRoutes(app)
	routers.SetupPrivateRoutes(app)

	// For GoRoutine implementation
	// routers.SetupPublicRoutesB(appb)
	// go func() {
	// 	err := appb.Listen(fmt.Sprintf(":8002"))
	// 	if err != nil {
	// 		log.Fatal(err.Error())
	// 	}
	// }()

	fmt.Println("Port: ", middleware.GetEnv("PORT"))
	// Serve the application
	if middleware.GetEnv("SSL") == "enabled" {
		log.Fatal(app.ListenTLS(
			fmt.Sprintf(":%s", middleware.GetEnv("PORT")),
			middleware.GetEnv("SSL_CERTIFICATE"),
			middleware.GetEnv("SSL_KEY"),
		))
	} else {
		err := app.Listen(fmt.Sprintf(":%s", middleware.GetEnv("PORT")))
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
