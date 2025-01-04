package infrastructure

import (
	"backend/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRouter() *fiber.App {
	app := fiber.New()
	app.Use(logger.New())
	// CORS ayarları
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",                            // İzin verilen originler
		AllowMethods:     "GET,POST,PUT,DELETE",          // İzin verilen HTTP yöntemleri
		AllowHeaders:     "Origin, Content-Type, Accept", // İzin verilen headerlar
		AllowCredentials: false,                          // Credential paylaşımına izin
	}))
	// Screenshot endpoint
	app.Post("/screenshot", handlers.TakeScreenshot)
	app.Get("/download/:filename", handlers.Download)
	app.Post("/contact", handlers.Contact)

	return app
}
