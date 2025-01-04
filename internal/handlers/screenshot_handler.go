package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"backend/internal/utils"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func TakeScreenshot(c *fiber.Ctx) error {
	var req models.Request
	if err := c.BodyParser(&req); err != nil {
		logrus.WithError(err).Error("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Message: "Invalid input",
			Details: "Unable to parse request body",
		})
	}

	if req.URL == "" {
		logrus.Error("URL is missing in the request")
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Message: "URL cannot be empty",
		})
	}

	// URL'nin başında http veya https yoksa ekle
	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		req.URL = "https://" + req.URL // Varsayılan olarak http ekliyoruz
	}

	if req.Width <= 0 || req.Height <= 0 {
		logrus.Error("Invalid width or height in the request")
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Message: "Invalid dimensions",
			Details: "Width and height must be positive integers",
		})
	}

	parsedURL, err := url.ParseRequestURI(req.URL)
	if err != nil {
		logrus.WithField("url", req.URL).WithError(err).Error("Invalid URL")
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Message: "Invalid URL",
			Details: "The provided URL is not valid",
		})
	}

	filePath, err := services.CaptureScreenshot(req, parsedURL.String())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url":   parsedURL.String(),
			"error": err,
		}).Error("Failed to capture screenshot")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Message: "Failed to capture screenshot",
			Details: "An error occurred while capturing the screenshot",
		})
	}

	downloadURL := fmt.Sprintf("%s", filePath)
	return c.JSON(fiber.Map{
		"message": "Screenshot captured successfully",
		"url":     downloadURL,
	})
}

func Download(c *fiber.Ctx) error {
	filename := c.Params("filename")
	filePath := "./screenshots/" + filename

	// Dosya var mı kontrol et
	if _, err := os.Stat(filePath); err != nil {
		logrus.WithFields(logrus.Fields{
			"filename": filename,
			"error":    err,
		}).Error("File not found")
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Message: "File not found",
			Details: "The requested file does not exist",
		})
	}

	return c.SendFile(filePath)
}

func Contact(c *fiber.Ctx) error {
	var contactMeModel models.ContactMe

	if err := c.BodyParser(&contactMeModel); err != nil {
		logrus.WithError(err).Error("Failed to parse contact form input")
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Message: "Invalid input",
			Details: "Unable to parse contact form data",
		})
	}

	mailBody := fmt.Sprintf("Email: %s\nMessage: %s\nName: %s",
		contactMeModel.Email,
		contactMeModel.Message,
		contactMeModel.Name,
	)

	isMailSend := utils.SendMail("veliarda.balci@gmail.com", mailBody, "Contact Me")
	if isMailSend != nil {
		logrus.WithError(isMailSend).Error("Failed to send email")
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Message: "Failed to send email",
			Details: "An error occurred while sending the email",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Mesajınız iletildi",
	})
}

func Ping(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{"Pong": true})
}
