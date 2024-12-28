package models

type ScreenshotResponse struct {
	Message string `json:"message"`
	URL     string `json:"url"`
}

type ContactMe struct {
	Email   string `json:"email"`
	Message string `json:"message"`
	Name    string `json:"name"`
}
