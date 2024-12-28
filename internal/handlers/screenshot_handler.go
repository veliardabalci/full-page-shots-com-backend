package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"backend/internal/utils"
	"fmt"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
)

func TakeScreenshot(c *fiber.Ctx) error {
	type request struct {
		URL string `json:"url"`
	}
	fmt.Println("geldi ")

	var req request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	parsedURL, err := url.ParseRequestURI(req.URL)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL",
		})
	}

	filePath, err := services.CaptureScreenshot(parsedURL.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Screenshot captured successfully",
		"url":     filePath,
	})
}

func Download(c *fiber.Ctx) error {
	filename := c.Params("filename")
	filePath := "./screenshots/" + filename

	// Dosya var mı kontrol et
	if _, err := os.Stat(filePath); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	// Dosyayı indirme olarak döndür
	return c.SendFile(filePath)
}

func Contact(c *fiber.Ctx) error {
	var contactMeModel models.ContactMe

	if err := c.BodyParser(&contactMeModel); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	mailBody := fmt.Sprintf("Email: %s\nMessage: %s\nName: %s",
		contactMeModel.Email,
		contactMeModel.Message,
		contactMeModel.Name,
	)

	isMailSend := utils.SendMail("veliarda.balci@gmail.com", mailBody, "Contact Me")
	if isMailSend != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error at send mail",
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Mesajınız iletildi",
	})
}
