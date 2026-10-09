package r2

import (
	"strconv"
	"time"

	"server/internal/middleware"
	_package "server/package"

	"github.com/gofiber/fiber/v3"
)


type Handler struct {
	service Service
}


func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}


func (h *Handler) GetAllUploadLog(c fiber.Ctx) error {
	// รับ query parameters
	limit := parseIntOrDefault(c.Query("limit"), 10)
	offset := parseIntOrDefault(c.Query("offset"), 0)
	page := parseIntOrDefault(c.Query("page"), 0)

	// ถ้ามี page parameter ให้คำนวณ offset จาก page
	if page > 0 {
		offset = (page - 1) * limit
	}

	// เรียก service เพื่อดึงข้อมูล
	logs, totalCount, err := h.service.GetAllUploadLog(c.Context(), limit, offset)
	if err != nil {
		return _package.SendJSON(
			c,
			fiber.StatusInternalServerError,
			_package.ErrorResponse(
				"Failed to fetch upload logs",
				"DATABASE_ERROR",
				err.Error(),
			),
		)
	}

	// คำนวณ pagination metadata
	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))
	currentPage := (offset / limit) + 1

	// สร้าง response data
	responseData := map[string]interface{}{
		"logs": logs,
		"pagination": map[string]interface{}{
			"total_records": totalCount,
			"total_pages":   totalPages,
			"current_page":  currentPage,
			"limit":         limit,
			"offset":        offset,
		},
	}

	return _package.SendJSON(
		c,
		fiber.StatusOK,
		_package.SuccessResponse(
			"Upload logs retrieved successfully",
			responseData,
		),
	)
}


func (h *Handler) GenerateUploadURL(c fiber.Ctx) error {
	
	// Generate unique object_key automatically
	objectKey := generateObjectKey()
	
	// Set default expires_in to 5 minutes
	expiresIn := 5 * time.Minute

	// Generate presigned URL
	presignedURL, err := h.service.GeneratePresignedUploadURL(c.Context(), objectKey, expiresIn)
	if err != nil {
		return _package.SendJSON(
			c,
			fiber.StatusInternalServerError,
			_package.ErrorResponse(
				"Failed to generate presigned URL",
				"R2_ERROR",
				err.Error(),
			),
		)
	}

	// Response
	responseData := map[string]interface{}{
		"presigned_url":       presignedURL,
		"object_key":          objectKey,
		"expires_in":          int(expiresIn.Seconds()),
		"method":              "PUT",
		"max_file_size":       5 * 1024 * 1024 * 1024, // 5GB in bytes
		"allowed_type":        "application/x-tar",
		"allowed_extensions":  []string{".tar"},
	}

	return _package.SendJSON(
		c,
		fiber.StatusOK,
		_package.SuccessResponse(
			"Presigned upload URL generated successfully",
			responseData,
		),
	)
}


func (h *Handler) GenerateDownloadURL(c fiber.Ctx) error {

	type PresignDownloadRequest struct {
		ObjectKey string `json:"object_key" validate:"required"`
	}

	var req PresignDownloadRequest
	if err := c.Bind().JSON(&req); err != nil {
		return _package.SendJSON(
			c,
			fiber.StatusBadRequest,
			_package.ErrorResponse(
				"Invalid request body",
				"VALIDATION_ERROR",
				err.Error(),
			),
		)
	}

	if req.ObjectKey == "" {
		return _package.SendJSON(
			c,
			fiber.StatusBadRequest,
			_package.ErrorResponse(
				"object_key is required",
				"VALIDATION_ERROR",
				"object_key cannot be empty",
			),
		)
	}

	// Set default expires_in to 5 minutes (same as upload)
	expiresIn := 5 * time.Minute

	presignedURL, err := h.service.GeneratePresignedDownloadURL(c.Context(), req.ObjectKey, expiresIn)
	if err != nil {
		return _package.SendJSON(
			c,
			fiber.StatusInternalServerError,
			_package.ErrorResponse(
				"Failed to generate presigned URL",
				"R2_ERROR",
				err.Error(),
			),
		)
	}

	// Response
	responseData := map[string]interface{}{
		"presigned_url": presignedURL,
		"object_key":    req.ObjectKey,
		"expires_in":    int(expiresIn.Seconds()),
		"method":        "GET",
	}

	return _package.SendJSON(
		c,
		fiber.StatusOK,
		_package.SuccessResponse(
			"Presigned download URL generated successfully",
			responseData,
		),
	)
}


func (h *Handler) RegisterRoutes(router fiber.Router) {

	router.Get("/r2", h.GetAllUploadLog)
	router.Post("/r2/presigned-upload", h.GenerateUploadURL)
	router.Post("/r2/presigned-download",
		middleware.ReqEmptyCheck, 
		h.GenerateDownloadURL,
	)
}

// --- Helper functions (optional) ---
func parseIntOrDefault(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// generateObjectKey generates a unique object key using timestamp and random string
func generateObjectKey() string {
	timestamp := time.Now().UnixNano()
	return strconv.FormatInt(timestamp, 10)
}
