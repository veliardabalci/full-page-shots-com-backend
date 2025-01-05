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

type Request struct {
	URL                 string `json:"url"`
	Width               int64  `json:"width"`
	Height              int64  `json:"height"`
	DeviceType          string `json:"deviceType"`
	SaveAs              string `json:"saveAs"`
	FullPageScreenshots bool   `json:"fullPageScreenshots"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
